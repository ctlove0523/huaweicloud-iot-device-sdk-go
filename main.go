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
		fmt.Println("I get message from platform")
		return true
	})

	time.Sleep(time.Hour)

}
