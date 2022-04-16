package iot

// 子设备添加回调函数
type SubDevicesAddHandler func(devices SubDeviceInfo)

//子设备删除糊掉函数
type SubDevicesDeleteHandler func(devices SubDeviceInfo)

// 处理平台下发的命令
type CommandHandler func(Command) (bool, interface{})

// 设备消息
type MessageHandler func(message Message) bool

// 平台设置设备属性
type DevicePropertiesSetHandler func(message DevicePropertyDownRequest) bool

// 平台查询设备属性
type DevicePropertyQueryHandler func(query DevicePropertyQueryRequest) DevicePropertyEntry

// 设备执行软件/固件升级.upgradeType = 0 软件升级，upgradeType = 1 固件升级
type DeviceUpgradeHandler func(upgradeType byte, info UpgradeInfo) UpgradeProgress

// 设备上报软固件版本,第一个返回值为软件版本，第二个返回值为固件版本
type SwFwVersionReporter func() (string, string)

// 平台下发的升级信息
type UpgradeInfo struct {
	Version     string `json:"version"`      //软固件包版本号
	Url         string `json:"url"`          //软固件包下载地址
	FileSize    int    `json:"file_size"`    //软固件包文件大小
	AccessToken string `json:"access_token"` //软固件包url下载地址的临时token
	Expires     string `json:"expires"`      //access_token的超期时间
	Sign        string `json:"sign"`         //软固件包MD5值
}

// 设备升级状态响应，用于设备向平台反馈进度，错误信息等
// ResultCode： 设备的升级状态，结果码定义如下：
// 0：处理成功
// 1：设备使用中
// 2：信号质量差
// 3：已经是最新版本
// 4：电量不足
// 5：剩余空间不足
// 6：下载超时
// 7：升级包校验失败
// 8：升级包类型不支持
// 9：内存不足
// 10：安装升级包失败
// 255： 内部异常
type UpgradeProgress struct {
	ResultCode  int    `json:"result_code"`
	Progress    int    `json:"progress"`    // 设备的升级进度，范围：0到100
	Version     string `json:"version"`     // 设备当前版本号
	Description string `json:"description"` // 升级状态描述信息，可以返回具体升级失败原因。
}

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
	ObjectDeviceId string `json:"object_device_id"`
	Name           string `json:"name"`
	Id             string `json:"id"`
	Content        string `json:"content"`
}

// 定义平台和设备之间的数据交换结构体

type Data struct {
	ObjectDeviceId string      `json:"object_device_id,omitempty"`
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
	EventTime string `json:"event_time,omitempty"`
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

// 上报设备信息请求
type ReportDeviceInfoRequest struct {
	ObjectDeviceId string                         `json:"object_device_id,omitempty"`
	Services       []ReportDeviceInfoServiceEvent `json:"services,omitempty"`
}

type ReportDeviceInfoServiceEvent struct {
	BaseServiceEvent
	Paras ReportDeviceInfoEventParas `json:"paras,omitempty"`
}

// 设备信息上报请求参数
type ReportDeviceInfoEventParas struct {
	DeviceSdkVersion string `json:"device_sdk_version,omitempty"`
	SwVersion        string `json:"sw_version,omitempty"`
	FwVersion        string `json:"fw_version,omitempty"`
}

// 上报设备日志请求
type ReportDeviceLogRequest struct {
	Services []ReportDeviceLogServiceEvent `json:"services,omitempty"`
}

type ReportDeviceLogServiceEvent struct {
	BaseServiceEvent
	Paras DeviceLogEntry `json:"paras,omitempty"`
}

// 设备状态日志收集器
type DeviceStatusLogCollector func(endTime string) []DeviceLogEntry

// 设备属性日志收集器
type DevicePropertyLogCollector func(endTime string) []DeviceLogEntry

// 设备消息日志收集器
type DeviceMessageLogCollector func(endTime string) []DeviceLogEntry

// 设备命令日志收集器
type DeviceCommandLogCollector func(endTime string) []DeviceLogEntry

type DeviceLogEntry struct {
	Timestamp string `json:"timestamp"` // 日志产生时间
	Type      string `json:"type"`      // 日志类型：DEVICE_STATUS，DEVICE_PROPERTY ，DEVICE_MESSAGE ，DEVICE_COMMAND
	Content   string `json:"content"`   // 日志内容
}
