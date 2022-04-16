package iot

import (
	"testing"
)

const deviceId = "611d13360ad1ed028658e089_device_cli"
const devicePwd = "123456789"
const server = "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883"
const qos = 1

func TestBaseIotDevice_Init(t *testing.T) {
	device := createBaseIotDevice()

	result := device.Init()

	if !result {
		t.Errorf("device init failed")
	}
}

func TestBaseIotDevice_IsConnected(t *testing.T) {
	device := createBaseIotDevice()
	device.Init()

	if !device.IsConnected() {
		t.Errorf("device connecte to server failed")
	}
}

func TestBaseIotDevice_DisConnect(t *testing.T) {
	device := createBaseIotDevice()
	device.Init()
	device.DisConnect()

	if device.IsConnected() {
		t.Errorf("device disconnect to server failed")
	}
}

func TestBaseIotDevice_AddMessageHandler(t *testing.T) {
	device := createBaseIotDevice()

	device.AddMessageHandler(func(message Message) bool {
		return true
	})

	if len(device.messageHandlers) == 0 {
		t.Errorf("add message handler failed")
	}
}

func TestBaseIotDevice_AddCommandHandler(t *testing.T) {
	device := createBaseIotDevice()

	device.AddCommandHandler(func(command Command) bool {
		return true
	})

	if len(device.commandHandler) == 0 {
		t.Errorf("add command handlers failed")
	}
}

func TestBaseIotDevice_AddPropertiesSetHandler(t *testing.T) {
	device := createBaseIotDevice()

	device.AddPropertiesSetHandler(func(message DevicePropertyDownRequest) bool {
		return true
	})

	if len(device.propertiesSetHandlers) == 0 {
		t.Errorf("add properties handler failed")
	}
}

func TestBaseIotDevice_SetPropertyQueryHandler(t *testing.T) {
	device := createBaseIotDevice()

	device.SetPropertyQueryHandler(func(query DevicePropertyQueryRequest) DevicePropertyEntry {
		return DevicePropertyEntry{}
	})

	if device.propertyQueryHandler == nil {
		t.Errorf("set property query handler failed")
	}
}

func TestBaseIotDevice_SetSwFwVersionReporter(t *testing.T) {
	device := createBaseIotDevice()

	device.SetSwFwVersionReporter(func() (string, string) {
		return "1.0", "2.0"
	})

	if device.swFwVersionReporter == nil {
		t.Errorf("set sw fw version reporter failed")
	}

}

func TestBaseIotDevice_SetDeviceUpgradeHandler(t *testing.T) {
	device := createBaseIotDevice()

	device.SetDeviceUpgradeHandler(func(upgradeType byte, info UpgradeInfo) UpgradeProgress {
		return UpgradeProgress{}
	})

	if device.deviceUpgradeHandler == nil {
		t.Errorf("set device upgrade handler failed")
	}
}

func createBaseIotDevice() baseIotDevice {
	device := baseIotDevice{}
	device.Id = deviceId
	device.Password = devicePwd
	device.Servers = server
	device.messageHandlers = []MessageHandler{}
	device.commandHandler = []CommandHandler{}

	device.fileUrls = map[string]string{}

	device.qos = qos
	device.batchSubDeviceSize = 10

	return device
}
