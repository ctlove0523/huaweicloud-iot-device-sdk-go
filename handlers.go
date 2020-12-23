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

// 文件上传下载管理

func CreateFileUploadResultResponse(filename string, result bool) FileResultResponse {
	code := 0
	if !result {
		code = 1
	}

	paras := FileResultServiceEventParas{
		ObjectName: filename,
		ResultCode: code,
	}

	serviceEvent := FileResultResponseServiceEvent{
		Paras: paras,
	}
	serviceEvent.ServiceId = "$file_manager"
	serviceEvent.EventType = "upload_result_report"
	serviceEvent.EventTime = GetEventTimeStamp()

	var services []FileResultResponseServiceEvent
	services = append(services, serviceEvent)

	response := FileResultResponse{
		Services: services,
	}

	return response
}

// 设备获取文件上传下载请求体
type FileRequest struct {
	ObjectDeviceId string                    `json:"object_device_id"`
	Services       []FileRequestServiceEvent `json:"services"`
}

// 平台下发文件上传和下载URL响应
type FileResponse struct {
	ObjectDeviceId string                     `json:"object_device_id"`
	Services       []FileResponseServiceEvent `json:"services"`
}

type FileResultResponse struct {
	ObjectDeviceId string                           `json:"object_device_id"`
	Services       []FileResultResponseServiceEvent `json:"services"`
}

type BaseServiceEvent struct {
	ServiceId string `json:"service_id"`
	EventType string `json:"event_type"`
	EventTime string `json:"event_time"`
}

type FileRequestServiceEvent struct {
	BaseServiceEvent
	Paras FileRequestServiceEventParas `json:"paras"`
}

type FileResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResponseServiceEventParas `json:"paras"`
}

type FileResultResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResultServiceEventParas `json:"paras"`
}

// 设备获取文件上传下载URL参数
type FileRequestServiceEventParas struct {
	FileName       string      `json:"file_name"`
	FileAttributes interface{} `json:"file_attributes"`
}

// 平台下发响应参数
type FileResponseServiceEventParas struct {
	Url            string      `json:"url"`
	BucketName     string      `json:"bucket_name"`
	ObjectName     string      `json:"object_name"`
	Expire         int         `json:"expire"`
	FileAttributes interface{} `json:"file_attributes"`
}

// 上报文件上传下载结果参数
type FileResultServiceEventParas struct {
	ObjectName        string `json:"object_name"`
	ResultCode        int    `json:"result_code"`
	StatusCode        int    `json:"status_code"`
	StatusDescription string `json:"status_description"`
}
