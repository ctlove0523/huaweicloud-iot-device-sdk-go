package iot

type DeviceConfig struct {
	Id                 string
	Password           string
	Servers            string
	Qos                byte
	BatchSubDeviceSize int
}
