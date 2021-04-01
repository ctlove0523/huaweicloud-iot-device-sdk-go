package main

import (
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
)

func main() {
	device := samples.CreateDevice()
	device.SetDeviceStatusLogCollector(func(endTime string) []iot.DeviceLogEntry {
		return []iot.DeviceLogEntry{}
	})
	device.Init()
}
