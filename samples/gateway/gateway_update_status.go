package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"strconv"
	"time"
)

func main() {
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-sdk", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
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
	TestUpdateSubDeviceState(device)
	time.Sleep(200 * time.Second)

}

func TestUpdateSubDeviceState(device iot.Device) {

	var devicesStatus []iot.DeviceStatus
	for i := 0; i < 200; i++ {
		subDevice := iot.DeviceStatus{
			DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-" + strconv.Itoa(i),
			Status:   "ONLINE",
		}

		devicesStatus = append(devicesStatus, subDevice)
	}

	ok := device.UpdateSubDeviceState(iot.SubDevicesStatus{
		DeviceStatuses: devicesStatus,
	})
	if ok {
		fmt.Println("gateway update sub devices status success")
	} else {
		fmt.Println("gateway update sub devices status failed")
	}
}

func TestDeleteSubDevices(device iot.Device) {
	ok := device.DeleteSubDevices([]string{"5fdb75cccbfe2f02ce81d4bf_sub-device-3"})
	if ok {
		fmt.Println("gateway send sub devices request success.")
	} else {
		fmt.Println("gateway send sub devices request failed.")
	}
}

func TestAddSubDevices(device iot.Device) {
	subDevices := []iot.DeviceInfo{{
		NodeId:    "sub-device-3",
		ProductId: "5fdb75cccbfe2f02ce81d4bf",
	}, {
		NodeId:    "sub-device-4",
		ProductId: "5fdb75cccbfe2f02ce81d4bf",
	}, {
		NodeId:    "sub-device-5",
		ProductId: "5fdb75cccbfe2f02ce81d4bf",
	}}

	ok := device.AddSubDevices(subDevices)
	if ok {
		fmt.Println("gateway add sub-devices success")
	} else {
		fmt.Println("gateway add sub-devices failed")
	}

}
