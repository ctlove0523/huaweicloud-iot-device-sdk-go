package main

import (
	"fmt"
	"huaweicloud-iot-device-sdk-go/iotdevice"
	"time"
)

func main() {
	device := iotdevice.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()
	commandHandler := SelfMessageHandler{}
	device.AddCommandHandler(commandHandler)

	time.Sleep(time.Hour)

}

type SelfMessageHandler struct {
}

func (handler SelfMessageHandler) HandleCommand(message iotdevice.IotCommand) bool {
	fmt.Println("I get message from platform")
	return true
}
