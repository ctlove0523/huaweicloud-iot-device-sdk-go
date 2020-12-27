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
type DevicePropertyQueryHandler func(query DevicePropertyQueryRequest) DevicePropertyEntry

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

// 定义平台和设备之间的数据交换结构体

type Data struct {
	ObjectDeviceId string      `json:"object_device_id"`
	Services       []DataEntry `json:"services"`
}

type DataEntry struct {
	ServiceId string      `json:"service_id"`
	EventType string      `json:"event_type"`
	EventTime string      `json:"event_time"`
	Paras     interface{} `json:"paras"` // 不同类型的请求paras使用的结构体不同
}

// 网关更新子设备状态
type SubDevicesStatus struct {
	DeviceStatuses []DeviceStatus `json:"device_statuses"`
}

type DeviceStatus struct {
	DeviceId string `json:"device_id"`
	Status   string `json:"status"` // 子设备状态。 OFFLINE：设备离线 ONLINE：设备上线
}

// 添加子设备
type SubDeviceInfo struct {
	Devices []DeviceInfo `json:"devices"`
	Version int          `json:"version"`
}

type DeviceInfo struct {
	ParentDeviceId string      `json:"parent_device_id,omitempty"`
	NodeId         string      `json:"node_id,omitempty"`
	DeviceId       string      `json:"device_id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	ManufacturerId string      `json:"manufacturer_id,omitempty"`
	Model          string      `json:"model,omitempty"`
	ProductId      string      `json:"product_id"`
	FwVersion      string      `json:"fw_version,omitempty"`
	SwVersion      string      `json:"sw_version,omitempty"`
	Status         string      `json:"status,omitempty"`
	ExtensionInfo  interface{} `json:"extension_info,omitempty"`
}

// 设备属性相关

// 设备属性
type DeviceProperties struct {
	Services []DevicePropertyEntry `json:"services"`
}

// 设备的一个属性
type DevicePropertyEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

// 平台设置设备属性==================================================
type DevicePropertyDownRequest struct {
	ObjectDeviceId string                           `json:"object_device_id"`
	Services       []DevicePropertyDownRequestEntry `json:"services"`
}

type DevicePropertyDownRequestEntry struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
}

// 平台设置设备属性==================================================
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
	DeviceId string                `json:"device_id"`
	Services []DevicePropertyEntry `json:"services"`
}

// 文件上传下载管理
func CreateFileUploadDownLoadResultResponse(filename, action string, result bool) FileResultResponse {
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
	if action == FileActionDownload {
		serviceEvent.EventType = "download_result_report"
	}
	if action == FileActionUpload {
		serviceEvent.EventType = "upload_result_report"
	}
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