package iot

import (
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
	"time"
)

type AsyncDevice interface {
	BaseDevice
	SendMessage(message Message) BooleanAsyncResult
	ReportProperties(properties DeviceProperties) BooleanAsyncResult
	BatchReportSubDevicesProperties(service DevicesService) BooleanAsyncResult
	QueryDeviceShadow(query DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler) BooleanAsyncResult
	UploadFile(filename string) BooleanAsyncResult
	DownloadFile(filename string) BooleanAsyncResult

	ReportDeviceInfo(swVersion, fwVersion string) BooleanAsyncResult
}

func CreateAsyncIotDevice(id, password, servers string) *asyncDevice {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      0,
	}

	return CreateAsyncIotDeviceWitConfig(config)
}

func CreateAsyncIotDeviceWithQos(id, password, servers string, qos byte) *asyncDevice {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      qos,
	}

	return CreateAsyncIotDeviceWitConfig(config)
}

func CreateAsyncIotDeviceWitConfig(config DeviceConfig) *asyncDevice {
	device := baseIotDevice{}
	device.Id = config.Id
	device.Password = config.Password
	device.Servers = config.Servers
	device.messageHandlers = []MessageHandler{}
	device.commandHandlers = []CommandHandler{}

	device.fileUrls = map[string]string{}

	device.qos = config.Qos
	device.batchSubDeviceSize = config.BatchSubDeviceSize

	result := &asyncDevice{
		base: device,
	}
	return result
}

type asyncDevice struct {
	base baseIotDevice
}

func (device *asyncDevice) Init() bool {
	return device.base.Init()
}

func (device *asyncDevice) DisConnect() () {
	device.base.DisConnect()
}

func (device *asyncDevice) IsConnected() bool {
	return device.base.IsConnected()
}

func (device *asyncDevice) AddMessageHandler(handler MessageHandler) {
	device.base.AddMessageHandler(handler)
}
func (device *asyncDevice) AddCommandHandler(handler CommandHandler) {
	device.base.AddCommandHandler(handler)
}
func (device *asyncDevice) AddPropertiesSetHandler(handler DevicePropertiesSetHandler) {
	device.base.AddPropertiesSetHandler(handler)
}
func (device *asyncDevice) SetPropertyQueryHandler(handler DevicePropertyQueryHandler) {
	device.base.SetPropertyQueryHandler(handler)
}

func (device *asyncDevice) SetSwFwVersionReporter(handler SwFwVersionReporter) {
	device.base.SetSwFwVersionReporter(handler)
}

func (device *asyncDevice) SetDeviceUpgradeHandler(handler DeviceUpgradeHandler) {
	device.base.SetDeviceUpgradeHandler(handler)
}

func (device *asyncDevice) SendMessage(message Message) BooleanAsyncResult {
	messageData := Interface2JsonString(message)

	result := BooleanAsyncResult{
		baseAsyncResult: baseAsyncResult{
			complete: make(chan struct{}),
		},
	}
	go func() {
		token := device.base.Client.Publish(FormatTopic(MessageUpTopic, device.base.Id), device.base.qos, false, messageData)
		token.Wait()
		if token.Error() != nil {
			result.SetResult(false)
			result.setError(token.Error())
			result.complete <- struct{}{}
		} else {
			result.SetResult(false)
			result.complete <- struct{}{}
		}
	}()

	return result
}

func (device *asyncDevice) ReportProperties(properties DeviceProperties) BooleanAsyncResult {
	propertiesData := Interface2JsonString(properties)
	result := BooleanAsyncResult{
		baseAsyncResult: baseAsyncResult{
			complete: make(chan struct{}),
		},
	}

	go func() {
		if token := device.base.Client.Publish(FormatTopic(PropertiesUpTopic, device.base.Id), device.base.qos, false, propertiesData);
			token.Wait() && token.Error() != nil {
			glog.Warningf("device %s report properties failed", device.base.Id)
			result.SetResult(false)
			result.setError(token.Error())
			result.complete <- struct{}{}
		} else {
			result.SetResult(false)
			result.complete <- struct{}{}
		}
	}()

	return result
}

func (device *asyncDevice) BatchReportSubDevicesProperties(service DevicesService) BooleanAsyncResult {
	result := BooleanAsyncResult{
		baseAsyncResult: baseAsyncResult{
			complete: make(chan struct{}),
		},
	}

	go func() {
		subDeviceCounts := len(service.Devices)

		batchReportSubDeviceProperties := 0
		if subDeviceCounts%device.base.batchSubDeviceSize == 0 {
			batchReportSubDeviceProperties = subDeviceCounts / device.base.batchSubDeviceSize
		} else {
			batchReportSubDeviceProperties = subDeviceCounts/device.base.batchSubDeviceSize + 1
		}

		loopResult := true
		for i := 0; i < batchReportSubDeviceProperties; i++ {
			begin := i * device.base.batchSubDeviceSize
			end := (i + 1) * device.base.batchSubDeviceSize
			if end > subDeviceCounts {
				end = subDeviceCounts
			}

			sds := DevicesService{
				Devices: service.Devices[begin:end],
			}

			if token := device.base.Client.Publish(FormatTopic(GatewayBatchReportSubDeviceTopic, device.base.Id), device.base.qos, false, Interface2JsonString(sds));
				token.Wait() && token.Error() != nil {
				glog.Warningf("device %s batch report sub device properties failed", device.base.Id)
				loopResult = false
				result.SetResult(false)
				result.setError(token.Error())
				result.complete <- struct{}{}
				break
			}
		}

		if loopResult {
			result.SetResult(true)
			result.complete <- struct{}{}
		}
	}()

	return result
}

func (device *asyncDevice) QueryDeviceShadow(query DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler) BooleanAsyncResult {
	device.base.propertiesQueryResponseHandler = handler
	result := BooleanAsyncResult{
		baseAsyncResult: baseAsyncResult{
			complete: make(chan struct{}),
		},
	}

	go func() {
		requestId := uuid.NewV4()
		if token := device.base.Client.Publish(FormatTopic(DeviceShadowQueryRequestTopic, device.base.Id)+requestId.String(), device.base.qos, false, Interface2JsonString(query));
			token.Wait() && token.Error() != nil {
			glog.Warningf("device %s query device shadow data failed,request id = %s", device.base.Id, requestId)
			result.setError(token.Error())
			result.SetResult(false)
		} else {
			result.SetResult(true)
		}
		result.complete <- struct{}{}
	}()

	return result
}

func (device *asyncDevice) UploadFile(filename string) BooleanAsyncResult {
	asyncResult:=BooleanAsyncResult{
		baseAsyncResult:baseAsyncResult{
			complete: make(chan struct{}),
		},
	}
	go func() {
		// 构造获取文件上传URL的请求
		requestParas := FileRequestServiceEventParas{
			FileName: filename,
		}

		serviceEvent := FileRequestServiceEvent{
			Paras: requestParas,
		}
		serviceEvent.ServiceId = "$file_manager"
		serviceEvent.EventTime = GetEventTimeStamp()
		serviceEvent.EventType = "get_upload_url"

		var services []FileRequestServiceEvent
		services = append(services, serviceEvent)
		request := FileRequest{
			Services: services,
		}

		if token := device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(request));
			token.Wait() && token.Error() != nil {
			glog.Warningf("publish file upload request url failed")
			asyncResult.setError(&DeviceError{
				errorMsg: "publish file upload request url failed",
			})
			asyncResult.SetResult(false)
			return
		}
		glog.Info("publish file upload request url success")

		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				_, ok := device.base.fileUrls[filename+FileActionUpload]
				if ok {
					glog.Infof("platform send file upload url success")
					goto ENDFOR
				}

			}
		}
	ENDFOR:

		if len(device.base.fileUrls[filename+FileActionUpload]) == 0 {
			glog.Errorf("get file upload url failed")
			asyncResult.setError(&DeviceError{
				errorMsg: "get file upload url failed",
			})
			asyncResult.SetResult(false)
			return
		}
		glog.Infof("file upload url is %s", device.base.fileUrls[filename+FileActionUpload])

		//filename = SmartFileName(filename)
		uploadFlag := CreateHttpClient().UploadFile(filename, device.base.fileUrls[filename+FileActionUpload])
		if !uploadFlag {
			glog.Errorf("upload file failed")
			asyncResult.setError(&DeviceError{
				errorMsg: "upload file failed",
			})
			asyncResult.SetResult(false)
			return
		}

		response := CreateFileUploadDownLoadResultResponse(filename, FileActionUpload, uploadFlag)

		token := device.base.Client.Publish(FormatTopic(PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, Interface2JsonString(response))
		if token.Wait() && token.Error() != nil {
			glog.Error("report file upload file result failed")
			asyncResult.setError(token.Error())
			asyncResult.SetResult(false)
		} else {
			asyncResult.SetResult(true)
		}
	}()

	return asyncResult
}

func (device *asyncDevice) DownloadFile(filename string) BooleanAsyncResult {
	asyncResult:=BooleanAsyncResult{
		baseAsyncResult:baseAsyncResult{
			complete: make(chan struct{}),
		},
	}
	go func() {
		// 构造获取文件上传URL的请求
		requestParas := FileRequestServiceEventParas{
			FileName: filename,
		}

		serviceEvent := FileRequestServiceEvent{
			Paras: requestParas,
		}
		serviceEvent.ServiceId = "$file_manager"
		serviceEvent.EventTime = GetEventTimeStamp()
		serviceEvent.EventType = "get_download_url"

		var services []FileRequestServiceEvent
		services = append(services, serviceEvent)
		request := FileRequest{
			Services: services,
		}

		if token := device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(request));
			token.Wait() && token.Error() != nil {
			glog.Warningf("publish file download request url failed")
			asyncResult.setError(&DeviceError{
				errorMsg: "publish file download request url failed",
			})
			asyncResult.SetResult(false)
			return
		}

		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				_, ok := device.base.fileUrls[filename+FileActionDownload]
				if ok {
					glog.Infof("platform send file upload url success")
					goto ENDFOR
				}

			}
		}
	ENDFOR:

		if len(device.base.fileUrls[filename+FileActionDownload]) == 0 {
			glog.Errorf("get file download url failed")
			asyncResult.setError(&DeviceError{
				errorMsg: "get file download url failed",
			})
			asyncResult.SetResult(false)
			return
		}

		downloadFlag := CreateHttpClient().DownloadFile(filename, device.base.fileUrls[filename+FileActionDownload])
		if !downloadFlag {
			glog.Errorf("down load file { %s } failed", filename)
			asyncResult.setError(&DeviceError{
				errorMsg: "down load file failed",
			})
			asyncResult.SetResult(false)
			return
		}

		response := CreateFileUploadDownLoadResultResponse(filename, FileActionDownload, downloadFlag)

		token := device.base.Client.Publish(FormatTopic(PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, Interface2JsonString(response))
		if token.Wait() && token.Error() != nil {
			glog.Error("report file upload file result failed")
			asyncResult.setError(token.Error())
			asyncResult.SetResult(false)
		} else {
			asyncResult.SetResult(true)
		}

	}()

	return asyncResult
}

func (device *asyncDevice) ReportDeviceInfo(swVersion, fwVersion string) BooleanAsyncResult {
	result := BooleanAsyncResult{
		baseAsyncResult: baseAsyncResult{
			complete: make(chan struct{}),
		},
	}
	go func() {
		event := ReportDeviceInfoServiceEvent{
			BaseServiceEvent{
				ServiceId: "$device_info",
				EventType: "device_info_report",
				EventTime: GetEventTimeStamp(),
			},
			ReportDeviceInfoEventParas{
				DeviceSdkVersion: SdkInfo()["sdk-version"],
				SwVersion:        swVersion,
				FwVersion:        fwVersion,
			},
		}

		request := ReportDeviceInfoRequest{
			ObjectDeviceId: device.base.Id,
			Services:       []ReportDeviceInfoServiceEvent{event},
		}

		token:=device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(request))
		if token.Wait() && token.Error()!= nil {
			result.setError(token.Error())
			result.SetResult(true)
		} else {
			result.SetResult(false)
		}

	}()

	return result
}
