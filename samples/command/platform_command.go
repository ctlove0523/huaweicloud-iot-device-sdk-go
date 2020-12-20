package main

import (
	"fmt"
	"huaweicloud-iot-device-sdk-go"
	"time"
)

// 处理平台下发的同步命令
func main() {
	// 创建一个设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()

	// 添加用于处理平台下发命令的callback
	device.AddCommandHandler(func(command iot.Command) bool {
		fmt.Println("I get command from platform")
		return true
	})
	time.Sleep(1 * time.Minute)
}
