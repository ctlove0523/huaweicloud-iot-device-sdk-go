package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
	uuid "github.com/satori/go.uuid"
)

func main() {
	// 创建一个设备并初始化
	device := samples.CreateDevice()
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

	for i := 0; i < 100; i++ {
		sendMsgResult := device.SendMessage(message)
		fmt.Printf("send message %v", sendMsgResult)
	}

}
