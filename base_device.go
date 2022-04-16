package iot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

const (
	// MessageDownTopic 平台下发消息topic
	MessageDownTopic string = "$oc/devices/{device_id}/sys/messages/down"

	// MessageUpTopic 设备上报消息topic
	MessageUpTopic string = "$oc/devices/{device_id}/sys/messages/up"

	// CommandDownTopic 平台下发命令topic
	CommandDownTopic string = "$oc/devices/{device_id}/sys/commands/#"

	// CommandResponseTopic 设备响应平台命令
	CommandResponseTopic string = "$oc/devices/{device_id}/sys/commands/response/request_id="

	// PropertiesUpTopic 设备上报属性
	PropertiesUpTopic string = "$oc/devices/{device_id}/sys/properties/report"

	// PropertiesSetRequestTopic 平台设置属性topic
	PropertiesSetRequestTopic string = "$oc/devices/{device_id}/sys/properties/set/#"

	// PropertiesSetResponseTopic 设备响应平台属性设置topic
	PropertiesSetResponseTopic string = "$oc/devices/{device_id}/sys/properties/set/response/request_id="

	// PropertiesQueryRequestTopic 平台查询设备属性
	PropertiesQueryRequestTopic string = "$oc/devices/{device_id}/sys/properties/get/#"

	// PropertiesQueryResponseTopic 设备响应平台属性查询
	PropertiesQueryResponseTopic string = "$oc/devices/{device_id}/sys/properties/get/response/request_id="

	// DeviceShadowQueryRequestTopic 设备侧获取平台的设备影子数据
	DeviceShadowQueryRequestTopic string = "$oc/devices/{device_id}/sys/shadow/get/request_id="

	// DeviceShadowQueryResponseTopic 设备侧响应获取平台设备影子
	DeviceShadowQueryResponseTopic string = "$oc/devices/{device_id}/sys/shadow/get/response/#"

	// GatewayBatchReportSubDeviceTopic 网关批量上报子设备属性
	GatewayBatchReportSubDeviceTopic string = "$oc/devices/{device_id}/sys/gateway/sub_devices/properties/report"

	// FileActionUpload 平台下发文件上传和下载URL
	FileActionUpload   string = "upload"
	FileActionDownload string = "download"

	// DeviceToPlatformTopic 设备或网关向平台发送请求
	DeviceToPlatformTopic string = "$oc/devices/{device_id}/sys/events/up"

	// PlatformEventToDeviceTopic 平台向设备下发事件topic
	PlatformEventToDeviceTopic string = "$oc/devices/{device_id}/sys/events/down"
)

const (
	AuthTypePassword uint8 = 0
	AuthTypeX509     uint8 = 1
)

type DeviceConfig struct {
	Id                 string
	Password           string
	VerifyTimestamp    bool
	Servers            string
	Qos                byte
	BatchSubDeviceSize int
	AuthType           uint8
	ServerCaPath       string
	CertFilePath       string
	CertKeyFilePath    string
	UseBootstrap       bool // 使用设备引导功能开关，true-使用，false-不使用
}

type BaseDevice interface {
	Init() bool
	DisConnect()
	IsConnected() bool

	AddMessageHandler(handler MessageHandler)
	AddCommandHandler(handler CommandHandler)
	AddPropertiesSetHandler(handler DevicePropertiesSetHandler)
	SetPropertyQueryHandler(handler DevicePropertyQueryHandler)
	SetSwFwVersionReporter(handler SwFwVersionReporter)
	SetDeviceUpgradeHandler(handler DeviceUpgradeHandler)

	SetDeviceStatusLogCollector(collector DeviceStatusLogCollector)
	SetDevicePropertyLogCollector(collector DevicePropertyLogCollector)
	SetDeviceMessageLogCollector(collector DeviceMessageLogCollector)
	SetDeviceCommandLogCollector(collector DeviceCommandLogCollector)
}

type LogCollectionConfig struct {
	rw               sync.RWMutex
	logCollectSwitch bool   //on：开启设备侧日志收集功能 off：关闭设备侧日志收集开关
	endTime          string // format yyyy-MM-dd'T'HH:mm:ss'Z'
}

func (lcc *LogCollectionConfig) setLogCollectSwitch(switchFlag bool) {
	lcc.rw.Lock()
	defer lcc.rw.Unlock()
	lcc.logCollectSwitch = switchFlag
}

func (lcc *LogCollectionConfig) setEndTime(endTime string) {
	lcc.rw.Lock()
	defer lcc.rw.Unlock()
	lcc.endTime = endTime
}

func (lcc *LogCollectionConfig) getLogCollectSwitch() bool {
	lcc.rw.RLock()
	defer lcc.rw.RUnlock()
	return lcc.logCollectSwitch
}

func (lcc *LogCollectionConfig) getEndTime() string {
	lcc.rw.RLock()
	defer lcc.rw.RUnlock()
	return lcc.endTime
}

type baseIotDevice struct {
	Id                             string // 设备Id，平台又称为deviceId
	Password                       string // 设备密码
	VerifyTimestamp                bool
	AuthType                       uint8  // 鉴权类型，0：密码认证；1：x.509证书认证
	ServerCaPath                   string // 平台CA证书
	CertFilePath                   string // 设备证书路径
	CertKeyFilePath                string // 设备证书key路径
	Servers                        string
	Client                         mqtt.Client
	commandHandler                 CommandHandler
	messageHandlers                []MessageHandler
	propertiesSetHandlers          []DevicePropertiesSetHandler
	propertyQueryHandler           DevicePropertyQueryHandler
	propertiesQueryResponseHandler DevicePropertyQueryResponseHandler
	subDevicesAddHandler           SubDevicesAddHandler
	subDevicesDeleteHandler        SubDevicesDeleteHandler
	swFwVersionReporter            SwFwVersionReporter
	deviceUpgradeHandler           DeviceUpgradeHandler
	fileUrls                       map[string]string
	qos                            byte
	batchSubDeviceSize             int
	lcc                            *LogCollectionConfig
	deviceStatusLogCollector       DeviceStatusLogCollector
	devicePropertyLogCollector     DevicePropertyLogCollector
	deviceMessageLogCollector      DeviceMessageLogCollector
	deviceCommandLogCollector      DeviceCommandLogCollector
	useBootstrap                   bool
}

func (device *baseIotDevice) DisConnect() {
	if device.Client != nil {
		device.Client.Disconnect(0)
	}
}
func (device *baseIotDevice) IsConnected() bool {
	if device.Client != nil {
		return device.Client.IsConnectionOpen()
	}

	return false
}

func (device *baseIotDevice) Init() bool {
	options := mqtt.NewClientOptions()
	if device.useBootstrap {
		bootstracpClient, err := NewBootstrapClient(device.Id, device.Password)
		if err != nil {
			fmt.Printf("create bootstrap client failed,err %s\n", err)
			return false
		}

		serverAddress := bootstracpClient.Boot()
		if len(serverAddress) == 0 {
			fmt.Println("get server address from bootstrap server failed")
			return false
		}
		options.AddBroker(serverAddress)
	} else {
		options.AddBroker(device.Servers)
	}

	options.SetClientID(assembleClientId(device))
	options.SetUsername(device.Id)
	options.SetPassword(hmacSha256(device.Password, timeStamp()))
	options.SetKeepAlive(250 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)
	if strings.Contains(device.Servers, "tls") || strings.Contains(device.Servers, "ssl") {
		glog.Infof("server support tls connection")

		// 设备使用x.509证书认证
		if device.AuthType == AuthTypeX509 {
			if len(device.ServerCaPath) == 0 || len(device.CertFilePath) == 0 || len(device.CertKeyFilePath) == 0 {
				glog.Error("device use x.509 auth but not set cert")
				panic("not set cert")
			}

			ca, err := ioutil.ReadFile(device.ServerCaPath)
			if err != nil {
				glog.Error("load server ca failed\n")
				panic(err)
			}
			serverCaPool := x509.NewCertPool()
			serverCaPool.AppendCertsFromPEM(ca)

			deviceCert, err := tls.LoadX509KeyPair(device.CertFilePath, device.CertKeyFilePath)
			if err != nil {
				glog.Error("load device cert failed")
				panic("load device cert failed")
			}
			var clientCerts []tls.Certificate
			clientCerts = append(clientCerts, deviceCert)

			cipherSuites := []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			}
			tlsConfig := &tls.Config{
				RootCAs:            serverCaPool,
				Certificates:       clientCerts,
				InsecureSkipVerify: true,
				MaxVersion:         tls.VersionTLS12,
				MinVersion:         tls.VersionTLS12,
				CipherSuites:       cipherSuites,
			}
			options.SetTLSConfig(tlsConfig)

		}

		if device.AuthType == 0 {
			options.SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			})
		}

	} else {
		options.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	device.Client = mqtt.NewClient(options)
	if token := device.Client.Connect(); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s init failed,error = %v", device.Id, token.Error())
		return false
	}

	device.subscribeDefaultTopics()

	go logFlush()

	return true

}

func (device *baseIotDevice) AddMessageHandler(handler MessageHandler) {
	if handler == nil {
		return
	}
	device.messageHandlers = append(device.messageHandlers, handler)
}
func (device *baseIotDevice) AddCommandHandler(handler CommandHandler) {
	if handler == nil {
		return
	}

	device.commandHandler = handler
}
func (device *baseIotDevice) AddPropertiesSetHandler(handler DevicePropertiesSetHandler) {
	if handler == nil {
		return
	}
	device.propertiesSetHandlers = append(device.propertiesSetHandlers, handler)
}
func (device *baseIotDevice) SetSwFwVersionReporter(handler SwFwVersionReporter) {
	device.swFwVersionReporter = handler
}

func (device *baseIotDevice) SetDeviceUpgradeHandler(handler DeviceUpgradeHandler) {
	device.deviceUpgradeHandler = handler
}

func (device *baseIotDevice) SetPropertyQueryHandler(handler DevicePropertyQueryHandler) {
	device.propertyQueryHandler = handler
}

func (device *baseIotDevice) SetDeviceStatusLogCollector(collector DeviceStatusLogCollector) {
	device.deviceStatusLogCollector = collector
}

func (device *baseIotDevice) SetDevicePropertyLogCollector(collector DevicePropertyLogCollector) {
	device.devicePropertyLogCollector = collector
}

func (device *baseIotDevice) SetDeviceMessageLogCollector(collector DeviceMessageLogCollector) {
	device.deviceMessageLogCollector = collector
}

func (device *baseIotDevice) SetDeviceCommandLogCollector(collector DeviceCommandLogCollector) {
	device.deviceCommandLogCollector = collector
}

func assembleClientId(device *baseIotDevice) string {
	segments := make([]string, 4)
	segments[0] = device.Id
	segments[1] = "0"
	if device.VerifyTimestamp {
		segments[2] = "1"
	} else {
		segments[2] = "0"
	}
	segments[3] = timeStamp()

	return strings.Join(segments, "_")
}

func logFlush() {
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			glog.Flush()
		}
	}
}

func (device *baseIotDevice) createCommandMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	commandHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			command := &Command{}
			if json.Unmarshal(message.Payload(), command) != nil {
				glog.Warningf("unmarshal platform command failed,device id = %s，message = %s", device.Id, message)
			}

			flag, response := device.commandHandler(*command)
			var res string
			if flag {
				glog.Infof("device %s handle command success", device.Id)
				res = Interface2JsonString(CommandResponse{
					ResultCode: 0,
					Paras:      response,
				})
			} else {
				glog.Warningf("device %s handle command failed", device.Id)
				res = Interface2JsonString(CommandResponse{
					ResultCode: 1,
					Paras:      response,
				})
			}
			if token := device.Client.Publish(formatTopic(CommandResponseTopic, device.Id)+getTopicRequestId(message.Topic()), 1, false, res); token.Wait() && token.Error() != nil {
				glog.Infof("device %s send command response failed", device.Id)
			}
		}()

	}

	return commandHandler
}

func (device *baseIotDevice) createPropertiesSetMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesSetHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			propertiesSetRequest := &DevicePropertyDownRequest{}
			if json.Unmarshal(message.Payload(), propertiesSetRequest) != nil {
				glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", device.Id, message)
			}

			handleFlag := true
			for _, handler := range device.propertiesSetHandlers {
				handleFlag = handleFlag && handler(*propertiesSetRequest)
			}

			var res string
			response := struct {
				ResultCode byte   `json:"result_code"`
				ResultDesc string `json:"result_desc"`
			}{}
			if handleFlag {
				response.ResultCode = 0
				response.ResultDesc = "Set property success."
				res = Interface2JsonString(response)
			} else {
				response.ResultCode = 1
				response.ResultDesc = "Set properties failed."
				res = Interface2JsonString(response)
			}
			if token := device.Client.Publish(formatTopic(PropertiesSetResponseTopic, device.Id)+getTopicRequestId(message.Topic()), device.qos, false, res); token.Wait() && token.Error() != nil {
				glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", device.Id, message)
			}
		}()
	}

	return propertiesSetHandler
}

func (device *baseIotDevice) createMessageMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	messageHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			msg := &Message{}
			if json.Unmarshal(message.Payload(), msg) != nil {
				glog.Warningf("unmarshal device message failed,device id = %s,message = %s", device.Id, message)
			}

			for _, handler := range device.messageHandlers {
				handler(*msg)
			}
		}()
	}

	return messageHandler
}

func (device *baseIotDevice) createPropertiesQueryMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesQueryHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			propertiesQueryRequest := &DevicePropertyQueryRequest{}
			if json.Unmarshal(message.Payload(), propertiesQueryRequest) != nil {
				glog.Warningf("device %s unmarshal properties query request failed %s", device.Id, message)
			}

			queryResult := device.propertyQueryHandler(*propertiesQueryRequest)
			responseToPlatform := Interface2JsonString(queryResult)
			if token := device.Client.Publish(formatTopic(PropertiesQueryResponseTopic, device.Id)+getTopicRequestId(message.Topic()), device.qos, false, responseToPlatform); token.Wait() && token.Error() != nil {
				glog.Warningf("device %s send properties query response failed.", device.Id)
			}
		}()
	}

	return propertiesQueryHandler
}

func (device *baseIotDevice) createPropertiesQueryResponseMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesQueryResponseHandler := func(client mqtt.Client, message mqtt.Message) {
		propertiesQueryResponse := &DevicePropertyQueryResponse{}
		if json.Unmarshal(message.Payload(), propertiesQueryResponse) != nil {
			glog.Warningf("device %s unmarshal property response failed,message %s", device.Id, Interface2JsonString(message))
		}
		device.propertiesQueryResponseHandler(*propertiesQueryResponse)
	}

	return propertiesQueryResponseHandler
}

func (device *baseIotDevice) subscribeDefaultTopics() {
	// 订阅平台命令下发topic
	topic := formatTopic(CommandDownTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createCommandMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe platform send command topic %s failed", device.Id, topic)
		panic(0)
	}

	// 订阅平台消息下发的topic
	topic = formatTopic(MessageDownTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createMessageMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device % subscribe platform send message topic %s failed.", device.Id, topic)
		panic(0)
	}

	// 订阅平台设置设备属性的topic
	topic = formatTopic(PropertiesSetRequestTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createPropertiesSetMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe platform set properties topic %s failed", device.Id, topic)
		panic(0)
	}

	// 订阅平台查询设备属性的topic
	topic = formatTopic(PropertiesQueryRequestTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createPropertiesQueryMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscriber platform query device properties topic failed %s", device.Id, topic)
		panic(0)
	}

	// 订阅查询设备影子响应的topic
	topic = formatTopic(DeviceShadowQueryResponseTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createPropertiesQueryResponseMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe query device shadow topic %s failed", device.Id, topic)
		panic(0)
	}

	// 订阅平台下发到设备的事件
	topic = formatTopic(PlatformEventToDeviceTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.handlePlatformToDeviceData()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe query device shadow topic %s failed", device.Id, topic)
		panic(0)
	}

}

// 平台向设备下发的事件callback
func (device *baseIotDevice) handlePlatformToDeviceData() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		data := &Data{}
		err := json.Unmarshal(message.Payload(), data)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, entry := range data.Services {
			eventType := entry.EventType
			switch eventType {
			case "add_sub_device_notify":
				// 子设备添加
				subDeviceInfo := &SubDeviceInfo{}
				if json.Unmarshal([]byte(Interface2JsonString(entry.Paras)), subDeviceInfo) != nil {
					continue
				}
				device.subDevicesAddHandler(*subDeviceInfo)
			case "delete_sub_device_notify":
				subDeviceInfo := &SubDeviceInfo{}
				if json.Unmarshal([]byte(Interface2JsonString(entry.Paras)), subDeviceInfo) != nil {
					continue
				}
				device.subDevicesDeleteHandler(*subDeviceInfo)

			case "get_upload_url_response":
				//获取文件上传URL
				fileResponse := &FileResponseServiceEventParas{}
				if json.Unmarshal([]byte(Interface2JsonString(entry.Paras)), fileResponse) != nil {
					continue
				}
				device.fileUrls[fileResponse.ObjectName+FileActionUpload] = fileResponse.Url
			case "get_download_url_response":
				fileResponse := &FileResponseServiceEventParas{}
				if json.Unmarshal([]byte(Interface2JsonString(entry.Paras)), fileResponse) != nil {
					continue
				}
				device.fileUrls[fileResponse.ObjectName+FileActionDownload] = fileResponse.Url
			case "version_query":
				// 查询软固件版本
				device.reportVersion()

			case "firmware_upgrade":
				upgradeInfo := &UpgradeInfo{}
				if json.Unmarshal([]byte(Interface2JsonString(entry.Paras)), upgradeInfo) != nil {
					continue
				}
				device.upgradeDevice(1, upgradeInfo)

			case "software_upgrade":
				upgradeInfo := &UpgradeInfo{}
				if json.Unmarshal([]byte(Interface2JsonString(entry.Paras)), upgradeInfo) != nil {
					continue
				}
				device.upgradeDevice(0, upgradeInfo)

			case "log_config":
				// 平台下发日志收集通知
				fmt.Println("platform send log collect command")
				logConfig := &LogCollectionConfig{}
				if json.Unmarshal([]byte(Interface2JsonString(entry.Paras)), logConfig) != nil {
					continue
				}

				lcc := &LogCollectionConfig{
					logCollectSwitch: logConfig.logCollectSwitch,
					endTime:          logConfig.endTime,
				}
				device.lcc = lcc
				device.reportLogsWorker()
			}
		}

	}

	return handler
}

func (device *baseIotDevice) reportLogsWorker() {
	go func() {
		for {
			if !device.lcc.getLogCollectSwitch() {
				break
			}
			logs := device.deviceStatusLogCollector(device.lcc.getEndTime())
			if len(logs) == 0 {
				fmt.Println("no log about device status")
				break
			}
			device.reportLogs(logs)
		}

	}()

	go func() {
		for {
			if !device.lcc.getLogCollectSwitch() {
				break
			}
			logs := device.devicePropertyLogCollector(device.lcc.getEndTime())
			if len(logs) == 0 {
				fmt.Println("no log about device property")
				break
			}
			device.reportLogs(logs)
		}

	}()

	go func() {
		for {
			if !device.lcc.getLogCollectSwitch() {
				break
			}
			logs := device.deviceMessageLogCollector(device.lcc.getEndTime())
			if len(logs) == 0 {
				fmt.Println("no log about device message")
				break
			}
			device.reportLogs(logs)
		}

	}()

	go func() {
		for {
			if !device.lcc.getLogCollectSwitch() {
				break
			}
			logs := device.deviceCommandLogCollector(device.lcc.getEndTime())
			if len(logs) == 0 {
				fmt.Println("no log about device command")
				break
			}
			device.reportLogs(logs)
		}

	}()

}

func (device *baseIotDevice) reportLogs(logs []DeviceLogEntry) {
	var dataEntries []DataEntry
	for _, log := range logs {
		dataEntry := DataEntry{
			ServiceId: "$log",
			EventType: "log_report",
			EventTime: GetEventTimeStamp(),
			Paras:     log,
		}
		dataEntries = append(dataEntries, dataEntry)
	}
	data := Data{
		Services: dataEntries,
	}

	reportedLog := Interface2JsonString(data)
	device.Client.Publish(formatTopic(DeviceToPlatformTopic, device.Id), 0, false, reportedLog)
}

func (device *baseIotDevice) reportVersion() {
	sw, fw := device.swFwVersionReporter()
	dataEntry := DataEntry{
		ServiceId: "$ota",
		EventType: "version_report",
		EventTime: GetEventTimeStamp(),
		Paras: struct {
			SwVersion string `json:"sw_version"`
			FwVersion string `json:"fw_version"`
		}{
			SwVersion: sw,
			FwVersion: fw,
		},
	}
	data := Data{
		ObjectDeviceId: device.Id,
		Services:       []DataEntry{dataEntry},
	}

	device.Client.Publish(formatTopic(DeviceToPlatformTopic, device.Id), device.qos, false, Interface2JsonString(data))
}

func (device *baseIotDevice) upgradeDevice(upgradeType byte, upgradeInfo *UpgradeInfo) {
	progress := device.deviceUpgradeHandler(upgradeType, *upgradeInfo)
	dataEntry := DataEntry{
		ServiceId: "$ota",
		EventType: "upgrade_progress_report",
		EventTime: GetEventTimeStamp(),
		Paras:     progress,
	}
	data := Data{
		ObjectDeviceId: device.Id,
		Services:       []DataEntry{dataEntry},
	}

	if token := device.Client.Publish(formatTopic(DeviceToPlatformTopic, device.Id), device.qos, false, Interface2JsonString(data)); token.Wait() && token.Error() != nil {
		glog.Errorf("device %s upgrade failed,type %d", device.Id, upgradeType)
	}
}
