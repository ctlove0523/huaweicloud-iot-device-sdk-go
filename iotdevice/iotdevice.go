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
	SendMessage(message Message) bool
	AddMessageHandler(handler MessageHandler)
}

type MessageHandler interface {
	Handle(message Message)
}

type iotDevice struct {
	Id              string
	Password        string
	Servers         string
	client          mqtt.Client
	messageHandlers []MessageHandler
	messageDownTopic string
}



func handle(client mqtt.Client,message mqtt.Message)  {

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

	// todo subscribe default topic
	device.client.Subscribe(device.messageDownTopic,2,handle)
	return true

}

func (device *iotDevice) IsConnected() bool {
	if device.client != nil {
		return device.client.IsConnected()
	}
	return false
}

func (device *iotDevice) SendMessage(message Message) bool {
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

func (device *iotDevice) AddMessageHandler(handler MessageHandler) {
	if handler == nil {
		return
	}
	device.messageHandlers = append(device.messageHandlers, handler)
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
	device.messageHandlers = []MessageHandler{}
	device.messageDownTopic = strings.ReplaceAll("$oc/devices/{device_id}/sys/messages/down","{device_id}", "id")

	return device
}
