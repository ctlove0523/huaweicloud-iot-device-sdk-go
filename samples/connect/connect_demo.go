package main

import (
	"fmt"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
)

func main() {
	device := samples.CreateDevice()

	initResult := device.Init()

	fmt.Printf("device init %v\n", initResult)

	fmt.Printf("device connected to server %v\n", device.IsConnected())

	device.DisConnect()

	fmt.Printf("device connected to server %v\n", device.IsConnected())

}
