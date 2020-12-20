package iot

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

type Device interface {
	Init() bool
	IsConnected() bool
	SendMessage(message Message) bool
	ReportProperties(properties ServiceProperty) bool
	BatchReportSubDevicesProperties(service DevicesService)
	QueryDeviceShadow(query DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler)
	AddMessageHandler(handler MessageHandler)
	AddCommandHandler(handler CommandHandler)
	AddPropertiesSetHandler(handler DevicePropertiesSetHandler)
	SetPropertyQueryHandler(handler DevicePropertyQueryHandler)
}

type iotDevice struct {
	Id                             string
	Password                       string
	Servers                        string
	client                         mqtt.Client
	commandHandlers                []CommandHandler
	messageHandlers                []MessageHandler
	propertiesSetHandlers          []DevicePropertiesSetHandler
	propertyQueryHandler           DevicePropertyQueryHandler
	propertiesQueryResponseHandler DevicePropertyQueryResponseHandler
	topics                         map[string]string
}

func (device *iotDevice) createMessageMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	messageHandler := func(client mqtt.Client, message mqtt.Message) {
		msg := &Message{}
		if json.Unmarshal(message.Payload(), msg) != nil {
			glog.Warningf("unmarshal device message failed,device id = %s,message = %s", device.Id, message)
		}

		for _, handler := range device.messageHandlers {
			handler(*msg)
		}
	}

	return messageHandler
}

func (device *iotDevice) createCommandMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	commandHandler := func(client mqtt.Client, message mqtt.Message) {
		command := &Command{}
		if json.Unmarshal(message.Payload(), command) != nil {
			glog.Warningf("unmarshal platform command failed,device id = %s，message = %s", device.Id, message)
		}

		handleFlag := true
		for _, handler := range device.commandHandlers {
			handleFlag = handleFlag && handler(*command)
		}
		var res string
		if handleFlag {
			glog.Infof("device %s handle command success", device.Id)
			res = Interface2JsonString(SuccessIotCommandResponse())
		} else {
			glog.Warningf("device %s handle command failed", device.Id)
			res = Interface2JsonString(FailedIotCommandResponse())
		}
		if token := device.client.Publish(device.topics[CommandResponseTopicName]+GetTopicRequestId(message.Topic()), 1, false, res);
			token.Wait() && token.Error() != nil {
			glog.Infof("device %s send command response failed", device.Id)
		}
	}

	return commandHandler
}

func (device *iotDevice) createPropertiesSetMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesSetHandler := func(client mqtt.Client, message mqtt.Message) {
		propertiesSetRequest := &DevicePropertyDownRequest{}
		if json.Unmarshal(message.Payload(), propertiesSetRequest) != nil {
			glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", device.Id, message)
		}

		handleFlag := true
		for _, handler := range device.propertiesSetHandlers {
			handleFlag = handleFlag && handler(*propertiesSetRequest)
		}

		var res string
		if handleFlag {
			res = Interface2JsonString(SuccessPropertiesSetResponse())
		} else {
			res = Interface2JsonString(FailedPropertiesSetResponse())
		}
		if token := device.client.Publish(device.topics[PropertiesSetResponseTopicName]+GetTopicRequestId(message.Topic()), 1, false, res);
			token.Wait() && token.Error() != nil {
			glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", device.Id, message)
		}
	}

	return propertiesSetHandler
}

func (device *iotDevice) createPropertiesQueryMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesQueryHandler := func(client mqtt.Client, message mqtt.Message) {
		propertiesQueryRequest := &DevicePropertyQueryRequest{}
		if json.Unmarshal(message.Payload(), propertiesQueryRequest) != nil {
			glog.Warningf("device %s unmarshal properties query request failed %s", device.Id, message)
		}

		queryResult := device.propertyQueryHandler(*propertiesQueryRequest)
		responseToPlatform := Interface2JsonString(queryResult)
		if token := device.client.Publish(device.topics[PropertiesQueryResponseTopicName]+GetTopicRequestId(message.Topic()), 1, false, responseToPlatform);
			token.Wait() && token.Error() != nil {
			glog.Warningf("device %s send properties query response failed.", device.Id)
		}
	}

	return propertiesQueryHandler
}

func (device *iotDevice) createPropertiesQueryResponseMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesQueryResponseHandler := func(client mqtt.Client, message mqtt.Message) {
		propertiesQueryResponse := &DevicePropertyQueryResponse{}
		if json.Unmarshal(message.Payload(), propertiesQueryResponse) != nil {
			glog.Warningf("device %s unmarshal property response failed,message %s", device.Id, Interface2JsonString(message))
		}
		device.propertiesQueryResponseHandler(*propertiesQueryResponse)
	}

	return propertiesQueryResponseHandler
}

func (device *iotDevice) Init() bool {
	options := mqtt.NewClientOptions()
	options.AddBroker(device.Servers)
	options.SetClientID(assembleClientId(device))
	options.SetUsername(device.Id)
	options.SetPassword(HmacSha256(device.Password, TimeStamp()))

	device.client = mqtt.NewClient(options)

	if token := device.client.Connect(); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s init failed", device.Id)
		return false
	}

	device.subscribeDefaultTopics()

	go logFlush()

	return true

}

func (device *iotDevice) IsConnected() bool {
	if device.client != nil {
		return device.client.IsConnected()
	}
	return false
}

func (device *iotDevice) SendMessage(message Message) bool {
	messageData := Interface2JsonString(message)
	if token := device.client.Publish(device.topics[MessageUpTopicName], 2, false, messageData);
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s send messagefailed", device.Id)
		return false
	}

	return true
}

func (device *iotDevice) ReportProperties(properties ServiceProperty) bool {
	propertiesData := Interface2JsonString(properties)
	if token := device.client.Publish(device.topics[PropertiesUpTopicName], 2, false, propertiesData);
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s report properties failed", device.Id)
		return false
	}
	return true
}

func (device *iotDevice) BatchReportSubDevicesProperties(service DevicesService) {
	if token := device.client.Publish(device.topics[GatewayBatchReportSubDeviceTopicName], 2, false, Interface2JsonString(service));
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s batch report sub device properties failed", device.Id)
	}
}

func (device *iotDevice) QueryDeviceShadow(query DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler) {
	device.propertiesQueryResponseHandler = handler
	requestId := uuid.NewV4()
	if token := device.client.Publish(device.topics[DeviceShadowQueryRequestTopicName]+requestId.String(), 2, false, Interface2JsonString(query));
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s query device shadow data failed,request id = %s", device.Id, requestId)
	}
}

func (device *iotDevice) AddMessageHandler(handler MessageHandler) {
	if handler == nil {
		return
	}
	device.messageHandlers = append(device.messageHandlers, handler)
}

func (device *iotDevice) AddCommandHandler(handler CommandHandler) {
	if handler == nil {
		return
	}

	device.commandHandlers = append(device.commandHandlers, handler)
}

func (device *iotDevice) AddPropertiesSetHandler(handler DevicePropertiesSetHandler) {
	if handler == nil {
		return
	}
	device.propertiesSetHandlers = append(device.propertiesSetHandlers, handler)
}

func (device *iotDevice) SetPropertyQueryHandler(handler DevicePropertyQueryHandler) {
	device.propertyQueryHandler = handler
}

func CreateIotDevice(id, password, servers string) Device {
	device := &iotDevice{}
	device.Id = id
	device.Password = password
	device.Servers = servers
	device.messageHandlers = []MessageHandler{}
	device.commandHandlers = []CommandHandler{}

	// 初始化设备相关的所有topic
	device.topics = make(map[string]string)
	device.topics[MessageDownTopicName] = FormatTopic(MessageDownTopic, id)
	device.topics[CommandDownTopicName] = FormatTopic(CommandDownTopic, id)
	device.topics[CommandResponseTopicName] = FormatTopic(CommandResponseTopic, id)
	device.topics[MessageUpTopicName] = FormatTopic(MessageUpTopic, id)
	device.topics[PropertiesUpTopicName] = FormatTopic(PropertiesUpTopic, id)
	device.topics[PropertiesSetRequestTopicName] = FormatTopic(PropertiesSetRequestTopic, id)
	device.topics[PropertiesSetResponseTopicName] = FormatTopic(PropertiesSetResponseTopic, id)
	device.topics[PropertiesQueryRequestTopicName] = FormatTopic(PropertiesQueryRequestTopic, id)
	device.topics[PropertiesQueryResponseTopicName] = FormatTopic(PropertiesQueryResponseTopic, id)
	device.topics[DeviceShadowQueryRequestTopicName] = FormatTopic(DeviceShadowQueryRequestTopic, id)
	device.topics[DeviceShadowQueryResponseTopicName] = FormatTopic(DeviceShadowQueryResponseTopic, id)
	device.topics[GatewayBatchReportSubDeviceTopicName] = FormatTopic(GatewayBatchReportSubDeviceTopic, id)
	return device
}

func assembleClientId(device *iotDevice) string {
	segments := make([]string, 4)
	segments[0] = device.Id
	segments[1] = "0"
	segments[2] = "0"
	segments[3] = TimeStamp()

	return strings.Join(segments, "_")
}

func logFlush() {
	ticker := time.Tick(5 * time.Second)
	for{
		select {
		case <-ticker:
			glog.Flush()
		}
	}
}

func (device *iotDevice) subscribeDefaultTopics() {
	// 订阅平台命令下发topic
	if token := device.client.Subscribe(device.topics[CommandDownTopicName], 2, device.createCommandMqttHandler());
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe platform send command topic %s failed", device.Id, device.topics[CommandDownTopicName])
		panic(0)
	}

	// 订阅平台消息下发的topic
	if token := device.client.Subscribe(device.topics[MessageDownTopicName], 2, device.createMessageMqttHandler());
		token.Wait() && token.Error() != nil {
		glog.Warningf("device % subscribe platform send message topic %s failed.", device.Id, device.topics[MessageDownTopicName])
		panic(0)
	}

	// 订阅平台设置设备属性的topic
	if token := device.client.Subscribe(device.topics[PropertiesSetRequestTopicName], 2, device.createPropertiesSetMqttHandler());
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe platform set properties topic %s failed", device.Id, device.topics[PropertiesSetRequestTopicName])
		panic(0)
	}

	// 订阅平台查询设备属性的topic
	if token := device.client.Subscribe(device.topics[PropertiesQueryRequestTopicName], 2, device.createPropertiesQueryMqttHandler());
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscriber platform query device properties topic failed %s", device.Id, device.topics[PropertiesQueryRequestTopicName])
		panic(0)
	}

	// 订阅查询设备影子响应的topic
	if token := device.client.Subscribe(device.topics[DeviceShadowQueryResponseTopicName], 2, device.createPropertiesQueryResponseMqttHandler());
		token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe query device shadow topic %s failed", device.Id, device.topics[DeviceShadowQueryResponseTopicName])
		panic(0)
	}

}
