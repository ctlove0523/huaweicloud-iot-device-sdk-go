package iot


const (
	// 平台下发消息topic
	MessageDownTopic     string = "$oc/devices/{device_id}/sys/messages/down"

	// 设备上报消息topic
	MessageUpTopic     string = "$oc/devices/{device_id}/sys/messages/up"

	// 平台下发命令topic
	CommandDownTopic string = "$oc/devices/{device_id}/sys/commands/#"

	// 设备响应平台命令
	CommandResponseTopic     string = "$oc/devices/{device_id}/sys/commands/response/request_id="

	// 设备上报属性
	PropertiesUpTopic     string = "$oc/devices/{device_id}/sys/properties/report"

	//平台设置属性topic
	//PropertiesSetRequestTopicName  string = "propertiesSetRequestTopicName"
	PropertiesSetRequestTopic      string = "$oc/devices/{device_id}/sys/properties/set/#"

	// 设备响应平台属性设置topic
	//PropertiesSetResponseTopicName string = "propertiesSetResponseTopicName"
	PropertiesSetResponseTopic     string = "$oc/devices/{device_id}/sys/properties/set/response/request_id="

	// 平台查询设备属性
	PropertiesQueryRequestTopicName  string = "propertiesQueryRequestTopicName"
	PropertiesQueryRequestTopic      string = "$oc/devices/{device_id}/sys/properties/get/#"
	PropertiesQueryResponseTopicName string = "propertiesQueryResponseTopicName"
	PropertiesQueryResponseTopic     string = "$oc/devices/{device_id}/sys/properties/get/response/request_id="

	// 设备侧获取平台的设备影子数据
	DeviceShadowQueryRequestTopicName  string = "deviceShadowQueryRequestTopicName"
	DeviceShadowQueryRequestTopic      string = "$oc/devices/{device_id}/sys/shadow/get/request_id="
	DeviceShadowQueryResponseTopicName string = "deviceShadowQueryResponseTopicName"
	DeviceShadowQueryResponseTopic     string = "$oc/devices/{device_id}/sys/shadow/get/response/#"

	// 网关批量上报子设备属性
	GatewayBatchReportSubDeviceTopicName string = "gatewayBatchReportSubDeviceTopicName"
	GatewayBatchReportSubDeviceTopic     string = "$oc/devices/{device_id}/sys/gateway/sub_devices/properties/report"

	// 文件上传请求：获取上传和下载URL，上报结果
	FileRequestTopicName string = "fileUploadRequestTopicName"
	FileRequestTopic     string = "$oc/devices/{device_id}/sys/events/up"

	// 平台下发文件上传和下载URL
	FileResponseTopicName string = "FileUploadResultTopic"

	FileActionUpload   string = "upload"
	FileActionDownload string = "download"

	// 设备或网关向平台发送请求
	DeviceToPlatformTopic string = "$oc/devices/{device_id}/sys/events/up"

	// 平台向设备下发事件topic
	PlatformEventToDeviceTopic string = "$oc/devices/{device_id}/sys/events/down"
)
