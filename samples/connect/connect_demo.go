package main

import (
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"strconv"
	"time"
)

const password = "123456789"
const server = "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883"
const produceId = "5fdb75cccbfe2f02ce81d4bf"


func main() {
	content := iot.Message{
		Content: "test content",
	}
	for i := 1; i <= 110; i++ {
		go SendMessage(strconv.Itoa(i), content)
	}

	time.Sleep(time.Hour)
}

func SendMessage(id string, message iot.Message) {
	device := iot.CreateIotDevice(produceId+"_test-"+id, password, server)
	device.Init()

	for i := 0; i < 100; i++ {
		device.SendMessage(message)
	}

	time.Sleep(time.Hour)

}
