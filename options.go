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