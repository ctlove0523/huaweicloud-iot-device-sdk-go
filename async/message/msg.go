package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

func main() {
	// 创建一个设备并初始化
	device := iot.CreateAsyncIotDevice("5fdb75cccbfe2f02ce81d4bf_liqian", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()

	// 注册平台下发消息的callback，当收到平台下发的消息时，调用此callback.
	// 支持注册多个callback，并且按照注册顺序调用
	device.AddMessageHandler(func(message iot.Message) bool {
		fmt.Println("first handler called" + iot.Interface2JsonString(message))
		return true
	})

	device.AddMessageHandler(func(message iot.Message) bool {
		fmt.Println("second handler called" + iot.Interface2JsonString(message))
		return true
	})

	//向平台发送消息
	message := iot.Message{
		ObjectDeviceId: uuid.NewV4().String(),
		Name:           "Fist send message to platform",
		Id:             uuid.NewV4().String(),
		Content:        "Hello Huawei IoT Platform",
	}
	asyncResult:=device.SendMessage(message)
	if asyncResult.Wait() && asyncResult.Error()!= nil {
		fmt.Println("async send message failed")
	} else {
		fmt.Println("async send message success")
	}
	time.Sleep(2 * time.Minute)

}
