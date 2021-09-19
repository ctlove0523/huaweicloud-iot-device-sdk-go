package main

import (
	"fmt"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
	"time"
)

// 处理平台下发的同步命令
func main() {
	// 创建一个设备并初始化
	device := samples.CreateDevice()

	device.Init()

	// 添加用于处理平台下发命令的callback
	commandProcessResult := false
	device.AddCommandHandler(func(command iot.Command) bool {
		fmt.Println("I get command from platform")
		commandProcessResult = true
		return true
	})
	time.Sleep(10 * time.Minute)
}
