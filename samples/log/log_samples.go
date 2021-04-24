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
	device.Init()
	var entries []iot.DeviceLogEntry

	for i := 0; i < 10; i++ {
		entry := iot.DeviceLogEntry{
			Type: "DEVICE_MESSAGE",
			//Timestamp: iot.GetEventTimeStamp(),
			Content: "message hello " + strconv.Itoa(i),
		}
		entries = append(entries, entry)
	}

	for i := 0; i < 100; i++ {
		result := device.ReportLogs(entries)
		fmt.Println(result)

	}

	time.Sleep(1 * time.Minute)
}
