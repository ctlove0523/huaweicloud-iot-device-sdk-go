package main

import (
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
)

func main() {
	//创建一个设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()

	device.UploadFile("D/software/mqttfx/chentong.txt")
	device.DownloadFile("D/software/mqttfx/chentong.txt")
}

