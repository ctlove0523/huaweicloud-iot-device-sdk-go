package main

import (
	"fmt"
	"huaweicloud-iot-device-sdk-go/handlers"
	"huaweicloud-iot-device-sdk-go/iotdevice"
	"time"
)

func main() {
	device := iotdevice.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()
	device.AddCommandHandler(func(command handlers.IotCommand) bool {
		fmt.Println("I get command from platform")
		return true
	})
	device.AddMessageHandler(func(message handlers.IotMessage) bool {
		fmt.Println("get message from platform")
		fmt.Println(message.Content)
		return true
	})

	message:=handlers.IotMessage{
		ObjectDeviceId: "chen tong",
		Name: "chen tong send message",
		Id: "id",
		Content: "hello platform",
	}

	device.SendMessage(message)
	time.Sleep(time.Hour)

}
