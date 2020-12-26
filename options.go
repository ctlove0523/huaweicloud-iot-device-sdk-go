package iot

// 子设备添加回调函数
type SubDevicesAddHandler func(devices SubDeviceInfo)

//子设备删除糊掉函数
type SubDevicesDeleteHandler func(devices SubDeviceInfo)

// 处理平台下发的命令
type CommandHandler func(Command) bool

// 设备消息
type MessageHandler func(message Message) bool

// 平台设置设备属性
type DevicePropertiesSetHandler func(message DevicePropertyDownRequest) bool

// 平台查询设备属性
type DevicePropertyQueryHandler func(query DevicePropertyQueryRequest) ServicePropertyEntry