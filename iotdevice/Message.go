package iotdevice

type IotMessage struct {
	ObjectDeviceId string      `json:"object_device_id"`
	Name           string      `json:"name"`
	Id             string      `json:"id"`
	Content        interface{} `json:"content"`
}

