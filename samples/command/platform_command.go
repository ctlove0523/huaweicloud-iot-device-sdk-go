package main

import (
	"fmt"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

// 处理平台下发的同步命令
func main() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	// 创建一个设备并初始化
	device := samples.CreateDevice()

	device.Init()

	// 添加用于处理平台下发命令的callback
	device.AddCommandHandler(func(command iot.Command) bool {
		fmt.Println("I get command from platform")
		time.Sleep(10 * time.Second)
		return true
	})
	time.Sleep(10 * time.Minute)
}
