package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
	"strconv"
	"time"
)

func main() {
	device := samples.CreateDevice()
	device.AddMessageHandler(func(message iot.Message) bool {
		fmt.Println(message)
		return true
	})
	device.SetSubDevicesAddHandler(func(devices iot.SubDeviceInfo) {
		fmt.Println(device)
	})
	device.SetSubDevicesDeleteHandler(func(devices iot.SubDeviceInfo) {
		fmt.Println(device)
	})
	device.SetDeviceStatusLogCollector(func(endTime string) []iot.DeviceLogEntry {
		fmt.Println("begin to collect log")
		entries := []iot.DeviceLogEntry{}

		for i := 0; i < 10; i++ {
			entry := iot.DeviceLogEntry{
				Type:      "DEVICE_MESSAGE",
				Timestamp: iot.GetEventTimeStamp(),
				Content:   "message hello " + strconv.Itoa(i),
			}
			entries = append(entries, entry)
		}
		return entries
	})
	device.Init()

	time.Sleep(1 * time.Minute)
}
