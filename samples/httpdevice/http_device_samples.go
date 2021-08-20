package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go/samples"
)

func main() {
	httpDevice := samples.CreateHttpDevice()

	s := SelfHttpResponse{
		HttpCode:    404,
		HttpMessage: "test http device",
		ReportTime:  iot.GetEventTimeStamp(),
	}
	entry := iot.DevicePropertyEntry{
		ServiceId:  "http_api",
		Properties: s,
		EventTime:  iot.GetEventTimeStamp(),
	}

	var entries []iot.DevicePropertyEntry
	entries = append(entries, entry)
	properties := iot.DeviceProperties{
		Services: entries,
	}

	fmt.Println(httpDevice.ReportProperties(properties))
}

type SelfHttpResponse struct {
	HttpCode    int    `json:"http_code"`
	HttpMessage string `json:"http_message"`
	ReportTime  string `json:"report_time"`
}
