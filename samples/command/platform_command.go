package main

import (
	"fmt"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
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
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_chentong", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()

	// 添加用于处理平台下发命令的callback
	device.AddCommandHandler(func(command iot.Command) bool {
		fmt.Println("I get command from platform")
		time.Sleep(10 * time.Second)
		return true
	})
	time.Sleep(10 * time.Minute)
}
