package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
)

func main() {
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()

	subDevice1 := iot.DeviceStatus{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-1",
		Status:   "ONLINE",
	}
	subDevice2 := iot.DeviceStatus{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-2",
		Status:   "ONLINE",
	}

	subDevice3 := iot.DeviceStatus{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-3",
		Status:   "ONLINE",
	}

	devicesStatus := []iot.DeviceStatus{subDevice1, subDevice2, subDevice3}

	ok := device.UpdateSubDeviceState(iot.SubDevicesStatus{
		DeviceStatuses: devicesStatus,
	})
	if ok {
		fmt.Println("gateway update sub devices status success")
	} else {
		fmt.Println("gateway update sub devices status failed")
	}
}
