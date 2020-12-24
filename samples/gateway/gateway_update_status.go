package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
)

func main() {
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()

	TestDeleteSubDevices(device, []string{"5fdb75cccbfe2f02ce81d4bf_sub-device-3"})
}

func TestUpdateSubDeviceState(device iot.Device) {
	subDevice1 := iot.DeviceStatus{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-1",
		Status:   "OFFLINE",
	}
	subDevice2 := iot.DeviceStatus{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-2",
		Status:   "OFFLINE",
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

func TestDeleteSubDevices(device iot.Device, deviceIds []string) {
	ok := device.DeleteSubDevices([]string{"5fdb75cccbfe2f02ce81d4bf_sub-device-3"})
	if ok {
		fmt.Println("gateway send sub devices request success.")
	} else {
		fmt.Println("gateway send sub devices request failed.")
	}
}
