package iotdevice

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

type IotDevice interface {
	Init() bool
	IsConnected() bool
	SendMessage(message IotMessage) bool
	AddMessageHandler(handler IotMessageHandler)
	AddCommandHandler(handler IotCommandHandler)
}

type IotMessageHandler interface {
	Handle(message IotMessage)
}

type IotCommandHandler interface {
	HandleCommand(message IotCommand) bool
}

type iotDevice struct {
	Id                   string
	Password             string
	Servers              string
	commandHandlers      []IotCommandHandler
	client               mqtt.Client
	messageHandlers      []IotMessageHandler
	messageDownTopic     string
	commandDownTopic     string
	commandResponseTopic string
}

func (device *iotDevice) createCommandMqttHandler() func(client mqtt.Client, message mqtt.Message)  {
	commandHandler := func(client mqtt.Client, message mqtt.Message) {
		command := &IotCommand{}
		if json.Unmarshal(message.Payload(), command) != nil {
			fmt.Println("unmarshal failed")
		}

		handleFlag := true
		for _, handler := range device.commandHandlers {
			handleFlag = handleFlag && handler.HandleCommand(*command)
		}
		var res string
		if handleFlag {
			res = JsonString(SuccessIotCommandResponse())
		} else {
			res = JsonString(FailedIotCommandResponse())
		}
		if token := device.client.Publish(device.commandResponseTopic+CommandRequestId(message.Topic()), 1, false, res);
			token.Wait() && token.Error() != nil {
			fmt.Println("send command response success")
		}
	}

	return commandHandler
}

func (device *iotDevice) Init() bool {
	options := mqtt.NewClientOptions()
	options.AddBroker(device.Servers)
	options.SetClientID(assembleClientId(device))
	options.SetUsername(device.Id)
	options.SetPassword(HmacSha256(device.Password, TimeStamp()))

	device.client = mqtt.NewClient(options)

	if token := device.client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("IoT device init failed,caulse %s\n", token.Error())
		return false
	}

	subscirbeToken := device.client.Subscribe(device.commandDownTopic, 2, device.createCommandMqttHandler())
	if subscirbeToken.Wait() && subscirbeToken.Error() != nil {
		fmt.Println(len(subscirbeToken.Error().Error()))
		fmt.Println("subscribe failed")
	} else {
		fmt.Println("subscribe success")
	}

	return true

}

func (device *iotDevice) IsConnected() bool {
	if device.client != nil {
		return device.client.IsConnected()
	}
	return false
}

func (device *iotDevice) SendMessage(message IotMessage) bool {
	topic := strings.Replace("$oc/devices/{device_id}/sys/messages/up", "{device_id}", device.Id, 1)
	messageData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("convert message to json format failed")
		return false
	}
	if token := device.client.Publish(topic, 2, false, string(messageData)); token.Wait() && token.Error() != nil {
		fmt.Println("send message failed")
		return false
	}

	return true
}

func (device *iotDevice) AddMessageHandler(handler IotMessageHandler) {
	if handler == nil {
		return
	}
	device.messageHandlers = append(device.messageHandlers, handler)
}

func (device *iotDevice) AddCommandHandler(handler IotCommandHandler) {
	if handler == nil {
		return
	}

	device.commandHandlers = append(device.commandHandlers, handler)
}

func assembleClientId(device *iotDevice) string {
	segments := make([]string, 4)
	segments[0] = device.Id
	segments[1] = "0"
	segments[2] = "0"
	segments[3] = TimeStamp()

	return strings.Join(segments, "_")
}

func CreateIotDevice(id, password, servers string) IotDevice {
	device := &iotDevice{}
	device.Id = id
	device.Password = password
	device.Servers = servers
	device.messageHandlers = []IotMessageHandler{}
	device.commandHandlers = []IotCommandHandler{}
	device.messageDownTopic = strings.ReplaceAll("$oc/devices/{device_id}/sys/messages/down", "{device_id}", id)
	device.commandDownTopic = strings.ReplaceAll("$oc/devices/{device_id}/sys/commands/#", "{device_id}", id)
	device.commandResponseTopic = strings.ReplaceAll("$oc/devices/{device_id}/sys/commands/response/request_id=", "{device_id}", id)

	return device
}
