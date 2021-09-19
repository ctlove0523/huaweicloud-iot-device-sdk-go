package iot

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestIotDevice_SendMessage(t *testing.T) {
	device := createIotDevice()
	device.Init()

	message := Message{
		ObjectDeviceId: uuid.NewV4().String(),
		Name:           "Fist send message to platform",
		Id:             uuid.NewV4().String(),
		Content:        "Hello Huawei IoT Platform",
	}
	if !device.SendMessage(message) {
		t.Errorf("device send message failed")
	}
}

func TestIotDevice_ReportProperties(t *testing.T) {
	device := createIotDevice()
	device.Init()

	props := DevicePropertyEntry{
		ServiceId: "value",
		EventTime: GetEventTimeStamp(),
		Properties: struct {
			Value   string `json:"value"`
			MsgType string `json:"msgType"`
		}{
			Value:   "Test Report",
			MsgType: "123",
		},
	}

	var content []DevicePropertyEntry
	content = append(content, props)
	services := DeviceProperties{
		Services: content,
	}

	reportResult := device.ReportProperties(services)
	if !reportResult {
		t.Error("device report property failed")
	}
}

func createIotDevice() Device {
	return CreateIotDevice(deviceId, devicePwd, server)
}
