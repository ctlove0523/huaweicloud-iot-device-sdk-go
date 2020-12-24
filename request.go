package iot

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
