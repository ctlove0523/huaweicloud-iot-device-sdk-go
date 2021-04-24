package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
	"time"
)

func main() {
	device := samples.CreateDevice()

	device.SetSubDevicesAddHandler(func(devices iot.SubDeviceInfo) {
		for _, info := range devices.Devices {
			fmt.Println("handle device add")
			fmt.Println(iot.Interface2JsonString(info))
		}
	})

	device.SetSubDevicesDeleteHandler(func(devices iot.SubDeviceInfo) {
		for _, info := range devices.Devices {
			fmt.Println("handle device delete")
			fmt.Println(iot.Interface2JsonString(info))
		}
	})

	device.Init()
	time.Sleep(200 * time.Second)

}
