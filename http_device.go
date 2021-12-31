package iot

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"sync"
	"time"
)

// 使用HTTP协议的设备，当前使用HTTP协议的设备只支持上报消息和上报属性
type HttpDevice interface {

	// 上报消息
	SendMessage(message Message) bool

	// 上报属性
	ReportProperties(properties DeviceProperties) bool
}

type restyHttpDevice struct {
	Id       string
	Password string
	Servers  string
	client   *resty.Client

	lock        sync.RWMutex
	accessToken string
}

func (device *restyHttpDevice) SendMessage(message Message) bool {
	resp, err := device.client.R().
		SetBody(message).
		Post(fmt.Sprintf("%s/v5/devices/%s/sys/messages/up", device.Servers, device.Id))
	if err != nil {
		fmt.Printf("send message failed %s\n", err)
	}
	return err == nil && resp.StatusCode() == http.StatusOK
}

func (device *restyHttpDevice) ReportProperties(properties DeviceProperties) bool {
	response, err := device.client.R().
		SetBody(properties).
		Post(fmt.Sprintf("%s/v5/devices/%s/sys/properties/report", device.Servers, device.Id))
	if err != nil {
		fmt.Printf("report properties failed %s\n", err)
	}
	return err == nil && response.StatusCode() == http.StatusOK
}

func (device *restyHttpDevice) init() {
	accessTokenBody := accessTokenRequest{
		DeviceId:  device.Id,
		SignType:  0,
		Timestamp: "2019120219",
		Password:  hmacSha256(device.Password, "2019120219"),
	}

	response, err := device.client.R().
		SetBody(accessTokenBody).
		Post(fmt.Sprintf("%s%s", device.Servers, "/v5/device-auth"))
	if err != nil {
		fmt.Printf("get device access token failed %s\n", err)
		return
	}

	tokenResponse := &accessTokenResponse{}
	err = json.Unmarshal(response.Body(), tokenResponse)
	if err != nil {
		fmt.Printf("json unmarshal failed %v", err)
		return
	}

	device.lock.Lock()
	device.accessToken = tokenResponse.AccessToken
	device.lock.Unlock()
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type accessTokenRequest struct {
	DeviceId  string `json:"device_id"`
	SignType  int    `json:"sign_type"`
	Timestamp string `json:"timestamp"`
	Password  string `json:"password"`
}

func CreateHttpDevice(config HttpDeviceConfig) HttpDevice {
	c := resty.New()
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	c.SetTimeout(30 * time.Second)
	c.SetRetryCount(3)
	c.SetRetryWaitTime(10 * time.Second)
	c.AddRetryCondition(func(response *resty.Response, err error) bool {
		return response.StatusCode() == http.StatusForbidden
	})

	connsPerHost := 10
	if config.MaxConnsPerHost != 0 {
		connsPerHost = config.MaxConnsPerHost
	}
	c.SetTransport(&http.Transport{
		MaxConnsPerHost: connsPerHost,
	})

	device := &restyHttpDevice{
		Id:       config.Id,
		Password: config.Password,
		Servers:  config.Server,
		client:   c,
		lock:     sync.RWMutex{},
	}

	device.init()
	device.client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		device.lock.RLock()
		request.SetHeader("access_token", device.accessToken)
		device.lock.RUnlock()
		request.SetHeader("Content-Type", "application/json")
		return nil
	})
	device.client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		if response.StatusCode() == http.StatusForbidden {
			device.init()
		}
		return nil
	})

	return device
}

type HttpDeviceConfig struct {
	Id              string
	Password        string
	Server          string // https://iot-mqtts.cn-north-4.myhuaweicloud.com:443
	MaxConnsPerHost int
	MaxIdleConns    int
}
