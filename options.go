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

// 设备命令
type Command struct {
	ObjectDeviceId string      `json:"object_device_id"`
	ServiceId      string      `json:"service_id""`
	CommandName    string      `json:"command_name"`
	Paras          interface{} `json:"paras"`
}

type CommandResponse struct {
	ResultCode   byte        `json:"result_code"`
	ResponseName string      `json:"response_name"`
	Paras        interface{} `json:"paras"`
}

// 消息
type Message struct {
	ObjectDeviceId string      `json:"object_device_id"`
	Name           string      `json:"name"`
	Id             string      `json:"id"`
	Content        string `json:"content"`
}