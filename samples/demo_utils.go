package samples

import iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"

const deviceId = "5fdb75cccbfe2f02ce81d4bf_liqian"
const devicePassword = "123456789"
const Server = "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883"

func CreateAsyncDevice() iot.AsyncDevice {
	device := iot.CreateAsyncIotDevice(deviceId, devicePassword, Server)

	return device
}

func CreateDevice() iot.Device {
	device := iot.CreateIotDevice(deviceId, devicePassword, Server)

	return device
}
