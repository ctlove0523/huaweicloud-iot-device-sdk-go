package iot

type CommandHandler func(Command) bool

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

func SuccessIotCommandResponse() CommandResponse {
	return CommandResponse{
		ResultCode: 0,
	}
}

func FailedIotCommandResponse() CommandResponse {
	return CommandResponse{
		ResultCode: 1,
	}
}

// 设备消息
type MessageHandler func(message Message) bool

type Message struct {
	ObjectDeviceId string      `json:"object_device_id"`
	Name           string      `json:"name"`
	Id             string      `json:"id"`
	Content        interface{} `json:"content"`
}

// 设备属性上报
type ServiceProperty struct {
	Services []ServicePropertyEntry `json:"services"`
}

type ServicePropertyEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

// 平台设置设备属性==================================================
type DevicePropertiesSetHandler func(message DevicePropertyDownRequest) bool

type DevicePropertyDownRequest struct {
	ObjectDeviceId string                           `json:"object_device_id"`
	Services       []DevicePropertyDownRequestEntry `json:"services"`
}

type DevicePropertyDownRequestEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
}

type DevicePropertyDownResponse struct {
	ResultCode byte   `json:"result_code"`
	ResultDesc string `json:"result_desc"`
}

func SuccessPropertiesSetResponse() DevicePropertyDownResponse {
	return DevicePropertyDownResponse{
		ResultCode: 0,
		ResultDesc: "success set properties",
	}
}

func FailedPropertiesSetResponse() DevicePropertyDownResponse {
	return DevicePropertyDownResponse{
		ResultCode: 1,
		ResultDesc: "failed set properties",
	}
}

// 平台设置设备属性==================================================

// 平台查询设备属性
type DevicePropertyQueryHandler func(query DevicePropertyQueryRequest) ServicePropertyEntry
type DevicePropertyQueryRequest struct {
	ObjectDeviceId string `json:"object_device_id"`
	ServiceId      string `json:"service_id"`
}

// 设备获取设备影子数据
type DevicePropertyQueryResponseHandler func(response DevicePropertyQueryResponse)

type DevicePropertyQueryResponse struct {
	ObjectDeviceId string             `json:"object_device_id"`
	Shadow         []DeviceShadowData `json:"shadow"`
}

type DeviceShadowData struct {
	ServiceId string                     `json:"service_id"`
	Desired   DeviceShadowPropertiesData `json:"desired"`
	Reported  DeviceShadowPropertiesData `json:"reported"`
	Version   int                        `json:"version"`
}
type DeviceShadowPropertiesData struct {
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

// 网关批量上报子设备属性

type DevicesService struct {
	Devices []DeviceService `json:"devices"`
}

type DeviceService struct {
	DeviceId string                 `json:"device_id"`
	Services []ServicePropertyEntry `json:"services"`
}
