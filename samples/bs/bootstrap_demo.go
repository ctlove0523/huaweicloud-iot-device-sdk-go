package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"time"
)

func main() {
	id := "611d13360ad1ed028658e089_zhou_sdk"
	pwd := "12345678901234567890"
	config := iot.DeviceConfig{
		Id:           id,
		Password:     pwd,
		UseBootstrap: true,
	}
	device := iot.CreateIotDeviceWitConfig(config)
	initRes := device.Init()
	fmt.Println(initRes)

	time.Sleep(1 * time.Minute)
}
