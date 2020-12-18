package main

import (
	"fmt"
	"huaweicloud-iot-device-sdk-go/iotdevice"
)

func main()  {
	device:=iotdevice.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt","123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()

	handler:=SelfMessageHandler{}
	device.AddMessageHandler(handler)

}

type SelfMessageHandler struct {

}

func (handler SelfMessageHandler) Handle(message iotdevice.Message)  {
	fmt.Println("I get message from platform")
}
