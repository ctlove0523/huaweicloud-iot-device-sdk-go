package iot

// 子设备添加回调函数
type SubDevicesAddHandler func(devices SubDeviceInfo)

type SubDevicesDeleteHandler func(devices SubDeviceInfo)