package main

import (
	"encoding/json"
	"fmt"
	"huaweicloud-iot-device-sdk-go/handlers"
	"huaweicloud-iot-device-sdk-go/iotdevice"
	"time"
)

func main() {
	device := iotdevice.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()
	device.AddCommandHandler(func(command handlers.IotCommand) bool {
		fmt.Println("I get command from platform")
		return true
	})
	device.AddMessageHandler(func(message handlers.IotMessage) bool {
		fmt.Println("get message from platform")
		fmt.Println(message.Content)
		return true
	})
	device.AddPropertiesSetHandler(func(message handlers.IotDevicePropertyDownRequest) bool {
		fmt.Println("I get property set command")
		fmt.Println(message)
		return true
	})
	device.SetPropertyQueryHandler(func(query handlers.IotDevicePropertyQueryRequest) handlers.IotServicePropertyEntry {
		return handlers.IotServicePropertyEntry{
			ServiceId: "value",
			Properties: SelfProperties{
				Value:   "QUERY RESPONSE",
				MsgType: "query property",
			},
			EventTime: "2020-12-19 02:23:24",
		}
	})

	message := handlers.IotMessage{
		ObjectDeviceId: "chen tong",
		Name:           "chen tong send message",
		Id:             "id",
		Content:        "hello platform",
	}

	device.SendMessage(message)

	props := handlers.IotServicePropertyEntry{
		ServiceId: "value",
		EventTime: "2020-12-19 02:23:24",
		Properties: SelfProperties{
			Value:   "chen tong",
			MsgType: "10",
		},
	}

	var content []handlers.IotServicePropertyEntry
	content = append(content, props)
	services := handlers.IotServiceProperty{
		Services: content,
	}
	device.ReportProperties(services)

	device.QueryDeviceShadow(handlers.IotDevicePropertyQueryRequest{
		ServiceId: "value",
	}, func(response handlers.IotDevicePropertyQueryResponse) {
		fmt.Println(response.Shadow)
	})

	// 批量上报子设备属性
	subDevice1 := handlers.IotDeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-1",
		Services: content,
	}
	subDevice2 := handlers.IotDeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-2",
		Services: content,
	}

	subDevice3 := handlers.IotDeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-3",
		Services: content,
	}

	var devices []handlers.IotDeviceService
	devices = append(devices, subDevice1, subDevice2, subDevice3)

	fmt.Println("begin to batch report")
	fmt.Println(JsonString(handlers.IotDevicesService{
		Devices: devices,
	}))
	device.BatchReportSubDevicesProperties(handlers.IotDevicesService{
		Devices: devices,
	})
	time.Sleep(time.Hour)

}

func JsonString(v interface{}) string {
	byteData, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(byteData)
}

type SelfProperties struct {
	Value   string `json:"value"`
	MsgType string `json:"msgType"`
}
