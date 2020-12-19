package handlers

type IotCommandHandler func(IotCommand) bool

// 设备命令
type IotCommand struct {
	ObjectDeviceId string      `json:"object_device_id"`
	ServiceId      string      `json:"service_id""`
	CommandName    string      `json:"command_name"`
	Paras          interface{} `json:"paras"`
}

type IotCommandResponse struct {
	ResultCode   byte        `json:"result_code"`
	ResponseName string      `json:"response_name"`
	Paras        interface{} `json:"paras"`
}

func SuccessIotCommandResponse() IotCommandResponse {
	return IotCommandResponse{
		ResultCode: 0,
	}
}

func FailedIotCommandResponse() IotCommandResponse {
	return IotCommandResponse{
		ResultCode: 1,
	}
}

// 设备消息
type IotMessageHandler func(message IotMessage) bool

type IotMessage struct {
	ObjectDeviceId string      `json:"object_device_id"`
	Name           string      `json:"name"`
	Id             string      `json:"id"`
	Content        interface{} `json:"content"`
}

// 设备属性上报
type IotServiceProperty struct {
	Services []IotServicePropertyEntry `json:"services"`
}

type IotServicePropertyEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

// 平台设置设备属性==================================================
type IotDevicePropertiesSetHandler func(message IotDevicePropertyDownRequest) bool

type IotDevicePropertyDownRequest struct {
	ObjectDeviceId string                              `json:"object_device_id"`
	Services       []IotDevicePropertyDownRequestEntry `json:"services"`
}

type IotDevicePropertyDownRequestEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
}

type IotDevicePropertyDownResponse struct {
	ResultCode byte   `json:"result_code"`
	ResultDesc string `json:"result_desc"`
}

func SuccessPropertiesSetResponse() IotDevicePropertyDownResponse {
	return IotDevicePropertyDownResponse{
		ResultCode: 0,
		ResultDesc: "success set properties",
	}
}

func FailedPropertiesSetResponse() IotDevicePropertyDownResponse {
	return IotDevicePropertyDownResponse{
		ResultCode: 1,
		ResultDesc: "failed set properties",
	}
}

// 平台设置设备属性==================================================

// 平台查询设备属性
type IotDevicePropertyQueryHandler func(query IotDevicePropertyQueryRequest) IotServicePropertyEntry
type IotDevicePropertyQueryRequest struct {
	ObjectDeviceId string `json:"object_device_id"`
	ServiceId      string `json:"service_id"`
}

// 设备获取设备影子数据
type IotDevicePropertyQueryResponseHandler func(response IotDevicePropertyQueryResponse)

type IotDevicePropertyQueryResponse struct {
	ObjectDeviceId string                `json:"object_device_id"`
	Shadow         []IotDeviceShadowData `json:"shadow"`
}

type IotDeviceShadowData struct {
	ServiceId string                        `json:"service_id"`
	Desired   IotDeviceShadowPropertiesData `json:"desired"`
	Reported  IotDeviceShadowPropertiesData `json:"reported"`
	Version   int                           `json:"version"`
}
type IotDeviceShadowPropertiesData struct {
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}
