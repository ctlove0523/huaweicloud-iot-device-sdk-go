package iotdevice

const (
	MessageDownTopic     string = "$oc/devices/{device_id}/sys/messages/down"
	MessageDownTopicName string = "messageDownTopicName"

	MessageUpTopic     string = "$oc/devices/{device_id}/sys/messages/up"
	MessageUpTopicName string = "messageUpTopicName"

	CommandDownTopicName string = "commandDownTopicName"
	CommandDownTopic     string = "$oc/devices/{device_id}/sys/commands/#"

	CommandResponseTopicName string = "commandResponseTopicName"
	CommandResponseTopic     string = "$oc/devices/{device_id}/sys/commands/response/request_id="

	PropertiesUpTopicName string = "propertiesUpTopicName"
	PropertiesUpTopic string = "$oc/devices/{device_id}/sys/properties/report"
)
