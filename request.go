package iot

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

//设备或网关请求的通用消息体
type Request struct {
	ObjectDeviceId string                `json:"object_device_id"`
	Services       []RequestEventService `json:"services"`
}
type RequestEventService struct {
	ServiceId string      `json:"service_id"`
	EventType string      `json:"event_type"`
	EventTime string      `json:"event_time"`
	Paras     interface{} `json:"paras"` // 不同类型的请求paras使用的结构体不同
}

// 1 网关更新子设备状态

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
	Version int        `json:"version"`
}

type DeviceInfo struct {
	ParentDeviceId string      `json:"parent_device_id"`
	NodeId         string      `json:"node_id"`
	DeviceId       string      `json:"device_id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	ManufacturerId string      `json:"manufacturer_id"`
	Model          string      `json:"model"`
	ProductId      string      `json:"product_id"`
	FwVersion      string      `json:"fw_version"`
	SwVersion      string      `json:"sw_version"`
	Status         string      `json:"status"`
	ExtensionInfo  interface{} `json:"extension_info"`
}
