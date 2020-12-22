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

func CreateFileUploadResultResponse(filename string, result bool) FileUploadResultResponse {
	code := 0
	if result {
		code = 0
	} else {
		code = 1
	}
	paras := UploadFileResultResponseServiceEventParas{
		ObjectName: filename,
		ResultCode: code,
	}

	serviceEvent := UploadResultResponseServiceEvent{
		Paras: paras,
	}
	serviceEvent.ServiceId = "$file_manager"
	serviceEvent.EventType = "upload_result_report"
	serviceEvent.EventTime = GetEventTimeStamp()

	var services []UploadResultResponseServiceEvent
	services = append(services, serviceEvent)

	response := FileUploadResultResponse{
		Services: services,
	}

	return response
}

type FileUploadUrlRequest struct {
	ObjectDeviceId string                      `json:"object_device_id"`
	Services       []UploadRequestServiceEvent `json:"services"`
}

type FileUploadUrlResponse struct {
	ObjectDeviceId string                       `json:"object_device_id"`
	Services       []UploadResponseServiceEvent `json:"services"`
}

type FileUploadResultResponse struct {
	ObjectDeviceId string                             `json:"object_device_id"`
	Services       []UploadResultResponseServiceEvent `json:"services"`
}

type BaseServiceEvent struct {
	ServiceId string `json:"service_id"`
	EventType string `json:"event_type"`
	EventTime string `json:"event_time"`
}

type UploadRequestServiceEvent struct {
	BaseServiceEvent
	Paras UploadRequestServiceEventParas `json:"paras"`
}

type UploadResponseServiceEvent struct {
	BaseServiceEvent
	Paras UploadResponseServiceEventParas `json:"paras"`
}

type UploadResultResponseServiceEvent struct {
	BaseServiceEvent
	Paras UploadFileResultResponseServiceEventParas `json:"paras"`
}

type UploadRequestServiceEventParas struct {
	FileName       string      `json:"file_name"`
	FileAttributes interface{} `json:"file_attributes"`
}

type UploadResponseServiceEventParas struct {
	Url            string      `json:"url"`
	BucketName     string      `json:"bucket_name"`
	ObjectName     string      `json:"object_name"`
	Expire         int         `json:"expire"`
	FileAttributes interface{} `json:"file_attributes"`
}

type UploadFileResultResponseServiceEventParas struct {
	ObjectName        string `json:"object_name"`
	ResultCode        int    `json:"result_code"`
	StatusCode        int    `json:"status_code"`
	StatusDescription string `json:"status_description"`
}
