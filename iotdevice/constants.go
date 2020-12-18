package iotdevice

const (
	MessageDownTopic     string = "$oc/devices/{device_id}/sys/messages/down"
	MessageDownTopicName string = "messageDownTopicName"

	CommandDownTopicName string = "commandDownTopicName"
	CommandDownTopic     string = "$oc/devices/{device_id}/sys/commands/#"

	CommandResponseTopicName string = "commandResponseTopicName"
	CommandResponseTopic     string = "$oc/devices/{device_id}/sys/commands/response/request_id="
)
