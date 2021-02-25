package main

import (
	"fmt"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"time"
)

// 处理平台下发的同步命令
func main() {
	// 创建一个设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_chentong", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()

	// 添加用于处理平台下发命令的callback
	device.AddCommandHandler(func(command iot.Command) bool {
		fmt.Println("I get command from platform")
		return true
	})
	time.Sleep(10 * time.Minute)
}
