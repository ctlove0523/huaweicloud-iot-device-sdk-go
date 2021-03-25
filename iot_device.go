package iot

import (
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Gateway interface {
	// 网关更新子设备状态
	UpdateSubDeviceState(subDevicesStatus SubDevicesStatus) bool

	// 网关删除子设备
	DeleteSubDevices(deviceIds []string) bool

	// 网关添加子设备
	AddSubDevices(deviceInfos []DeviceInfo) bool

	// 设置平台添加子设备回调函数
	SetSubDevicesAddHandler(handler SubDevicesAddHandler)

	// 设置平台删除子设备回调函数
	SetSubDevicesDeleteHandler(handler SubDevicesDeleteHandler)

	// 网关同步子设备列表,默认实现不指定版本
	SyncAllVersionSubDevices()

	// 网关同步特定版本子设备列表
	SyncSubDevices(version int)
}

type Device interface {
	Gateway
	BaseDevice
	SendMessage(message Message) bool
	ReportProperties(properties DeviceProperties) bool
	BatchReportSubDevicesProperties(service DevicesService)
	QueryDeviceShadow(query DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler)
	UploadFile(filename string) bool
	DownloadFile(filename string) bool

	ReportDeviceInfo(swVersion, fwVersion string)
}

type iotDevice struct {
	base baseIotDevice
}

func (device *iotDevice) Init() bool {
	return device.base.Init()
}

func (device *iotDevice) DisConnect() {
	device.base.DisConnect()
}

func (device *iotDevice) IsConnected() bool {
	return device.base.IsConnected()
}

func (device *iotDevice) AddMessageHandler(handler MessageHandler) {
	device.base.AddMessageHandler(handler)
}
func (device *iotDevice) AddCommandHandler(handler CommandHandler) {
	device.base.AddCommandHandler(handler)
}
func (device *iotDevice) AddPropertiesSetHandler(handler DevicePropertiesSetHandler) {
	device.base.AddPropertiesSetHandler(handler)
}
func (device *iotDevice) SetSwFwVersionReporter(handler SwFwVersionReporter) {
	device.base.SetSwFwVersionReporter(handler)
}

func (device *iotDevice) SetDeviceUpgradeHandler(handler DeviceUpgradeHandler) {
	device.base.SetDeviceUpgradeHandler(handler)
}

func (device *iotDevice) SetPropertyQueryHandler(handler DevicePropertyQueryHandler) {
	device.base.SetPropertyQueryHandler(handler)
}

func (device *iotDevice) SendMessage(message Message) bool {
	messageData := Interface2JsonString(message)
	if token := device.base.Client.Publish(FormatTopic(MessageUpTopic, device.base.Id), device.base.qos, false, messageData);
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s send messagefailed", device.base.Id)
		return false
	}
	return true
}

func (device *iotDevice) ReportProperties(properties DeviceProperties) bool {
	propertiesData := Interface2JsonString(properties)
	if token := device.base.Client.Publish(FormatTopic(PropertiesUpTopic, device.base.Id), device.base.qos, false, propertiesData);
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s report properties failed", device.base.Id)
		return false
	}
	return true
}
func (device *iotDevice) BatchReportSubDevicesProperties(service DevicesService) {

	subDeviceCounts := len(service.Devices)

	batchReportSubDeviceProperties := 0
	if subDeviceCounts%device.base.batchSubDeviceSize == 0 {
		batchReportSubDeviceProperties = subDeviceCounts / device.base.batchSubDeviceSize
	} else {
		batchReportSubDeviceProperties = subDeviceCounts/device.base.batchSubDeviceSize + 1
	}

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
		}
	}
}

func (device *iotDevice) QueryDeviceShadow(query DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler) {
	device.base.propertiesQueryResponseHandler = handler
	requestId := uuid.NewV4()
	if token := device.base.Client.Publish(FormatTopic(DeviceShadowQueryRequestTopic, device.base.Id)+requestId.String(), device.base.qos, false, Interface2JsonString(query));
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s query device shadow data failed,request id = %s", device.base.Id, requestId)
	}
}

func (device *iotDevice) UploadFile(filename string) bool {
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
		return false
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
		return false
	}
	glog.Infof("file upload url is %s", device.base.fileUrls[filename+FileActionUpload])

	//filename = SmartFileName(filename)
	uploadFlag := CreateHttpClient().UploadFile(filename, device.base.fileUrls[filename+FileActionUpload])
	if !uploadFlag {
		glog.Errorf("upload file failed")
		return false
	}

	response := CreateFileUploadDownLoadResultResponse(filename, FileActionUpload, uploadFlag)

	token := device.base.Client.Publish(FormatTopic(PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, Interface2JsonString(response))
	if token.Wait() && token.Error() != nil {
		glog.Error("report file upload file result failed")
		return false
	}

	return true
}

func (device *iotDevice) DownloadFile(filename string) bool {
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
		return false
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
		return false
	}

	downloadFlag := CreateHttpClient().DownloadFile(filename, device.base.fileUrls[filename+FileActionDownload])
	if !downloadFlag {
		glog.Errorf("down load file { %s } failed", filename)
		return false
	}

	response := CreateFileUploadDownLoadResultResponse(filename, FileActionDownload, downloadFlag)

	token := device.base.Client.Publish(FormatTopic(PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, Interface2JsonString(response))
	if token.Wait() && token.Error() != nil {
		glog.Error("report file upload file result failed")
		return false
	}

	return true
}

func (device *iotDevice) ReportDeviceInfo(swVersion, fwVersion string) {
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

	device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(request))
}

func (device *iotDevice) UpdateSubDeviceState(subDevicesStatus SubDevicesStatus) bool {
	glog.Infof("begin to update sub-devices status")

	subDeviceCounts := len(subDevicesStatus.DeviceStatuses)

	batchUpdateSubDeviceState := 0
	if subDeviceCounts%device.base.batchSubDeviceSize == 0 {
		batchUpdateSubDeviceState = subDeviceCounts / device.base.batchSubDeviceSize
	} else {
		batchUpdateSubDeviceState = subDeviceCounts/device.base.batchSubDeviceSize + 1
	}

	for i := 0; i < batchUpdateSubDeviceState; i++ {
		begin := i * device.base.batchSubDeviceSize
		end := (i + 1) * device.base.batchSubDeviceSize
		if end > subDeviceCounts {
			end = subDeviceCounts
		}

		sds := SubDevicesStatus{
			DeviceStatuses: subDevicesStatus.DeviceStatuses[begin:end],
		}

		requestEventService := DataEntry{
			ServiceId: "$sub_device_manager",
			EventType: "sub_device_update_status",
			EventTime: GetEventTimeStamp(),
			Paras:     sds,
		}

		request := Data{
			ObjectDeviceId: device.base.Id,
			Services:       []DataEntry{requestEventService},
		}

		if token := device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(request));
			token.Wait() && token.Error() != nil {
			glog.Warningf("gateway %s update sub devices status failed", device.base.Id)
			return false
		}
	}

	glog.Info("gateway  update sub devices status failed", device.base.Id)
	return true
}

func (device *iotDevice) DeleteSubDevices(deviceIds []string) bool {
	glog.Infof("begin to delete sub-devices %s", deviceIds)

	subDevices := struct {
		Devices []string `json:"devices"`
	}{
		Devices: deviceIds,
	}

	requestEventService := DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "delete_sub_device_request",
		EventTime: GetEventTimeStamp(),
		Paras:     subDevices,
	}

	request := Data{
		ObjectDeviceId: device.base.Id,
		Services:       []DataEntry{requestEventService},
	}

	if token := device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(request));
		token.Wait() && token.Error() != nil {
		glog.Warningf("gateway %s delete sub devices request send failed", device.base.Id)
		return false
	}

	glog.Warningf("gateway %s delete sub devices request send success", device.base.Id)
	return true
}

func (device *iotDevice) AddSubDevices(deviceInfos []DeviceInfo) bool {
	devices := struct {
		Devices []DeviceInfo `json:"devices"`
	}{
		Devices: deviceInfos,
	}

	requestEventService := DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "add_sub_device_request",
		EventTime: GetEventTimeStamp(),
		Paras:     devices,
	}

	request := Data{
		ObjectDeviceId: device.base.Id,
		Services:       []DataEntry{requestEventService},
	}

	if token := device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(request));
		token.Wait() && token.Error() != nil {
		glog.Warningf("gateway %s add sub devices request send failed", device.base.Id)
		return false
	}

	glog.Warningf("gateway %s add sub devices request send success", device.base.Id)
	return true
}

func (device *iotDevice) SetSubDevicesAddHandler(handler SubDevicesAddHandler) {
	device.base.subDevicesAddHandler = handler
}
func (device *iotDevice) SetSubDevicesDeleteHandler(handler SubDevicesDeleteHandler) {
	device.base.subDevicesDeleteHandler = handler
}

func (device *iotDevice) SyncAllVersionSubDevices() {
	dataEntry := DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "sub_device_sync_request",
		EventTime: GetEventTimeStamp(),
		Paras: struct {
		}{},
	}

	var dataEntries []DataEntry
	dataEntries = append(dataEntries, dataEntry)

	data := Data{
		Services: dataEntries,
	}

	if token := device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(data));
		token.Wait() && token.Error() != nil {
		glog.Errorf("send sub device sync request failed")
	}
}

func (device *iotDevice) SyncSubDevices(version int) {
	syncParas := struct {
		Version int `json:"version"`
	}{
		Version: version,
	}

	dataEntry := DataEntry{
		ServiceId: "$sub_device_manager",
		EventType: "sub_device_sync_request",
		EventTime: GetEventTimeStamp(),
		Paras:     syncParas,
	}

	var dataEntries []DataEntry
	dataEntries = append(dataEntries, dataEntry)

	data := Data{
		Services: dataEntries,
	}

	if token := device.base.Client.Publish(FormatTopic(DeviceToPlatformTopic, device.base.Id), device.base.qos, false, Interface2JsonString(data));
		token.Wait() && token.Error() != nil {
		glog.Errorf("send sync sub device request failed")
	}
}

func CreateIotDevice(id, password, servers string) Device {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      0,
	}

	return CreateIotDeviceWitConfig(config)
}

func CreateIotDeviceWithQos(id, password, servers string, qos byte) Device {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      qos,
	}

	return CreateIotDeviceWitConfig(config)
}

func CreateIotDeviceWitConfig(config DeviceConfig) Device {
	device := baseIotDevice{}
	device.Id = config.Id
	device.Password = config.Password
	device.Servers = config.Servers
	device.messageHandlers = []MessageHandler{}
	device.commandHandlers = []CommandHandler{}

	device.fileUrls = map[string]string{}

	device.qos = config.Qos
	device.batchSubDeviceSize = config.BatchSubDeviceSize

	result := &iotDevice{
		base: device,
	}
	return result
}
