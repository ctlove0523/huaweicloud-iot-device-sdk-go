package samples

import iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"

const deviceId = "611d13360ad1ed028658e089_http_report"
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

func CreateHttpDevice() iot.HttpDevice {
	config := iot.HttpDeviceConfig{
		Id:              deviceId,
		Password:        devicePassword,
		Server:          "https://iot-mqtts.cn-north-4.myhuaweicloud.com:443",
		MaxConnsPerHost: 2,
		MaxIdleConns:    0,
	}
	return iot.CreateHttpDevice(config)

}
