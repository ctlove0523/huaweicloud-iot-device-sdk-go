package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

func main1() {
	res := &Result{
		complete: make(chan struct{}),
	}
	go func() {
		time.Sleep(5 * time.Second)
		res.SetResult(false)
	}()
	if res.Wait() {
		fmt.Println(res.Result())
	}
}

type Result struct {
	m        sync.RWMutex
	complete chan struct{}
	value    bool
}

func (r *Result) Wait() bool {
	fmt.Println("begin to wait")
	<-r.complete

	fmt.Println("end to wait")
	return true
}

func (r *Result) SetResult(v bool) {
	fmt.Println("begin to set value")
	r.m.Lock()
	defer r.m.Unlock()
	r.value = v
	close(r.complete)
	fmt.Println("end to set value")
}

func (r *Result) Result() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.value
}

func main() {
	// 创建一个设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-sdk", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()

	//向平台发送消息
	message := iot.Message{
		ObjectDeviceId: uuid.NewV4().String(),
		Name:           "Fist send message to platform",
		Id:             uuid.NewV4().String(),
		Content:        "Hello Huawei IoT Platform",
	}
	result := device.SendMessage(message)

	fmt.Println("send message async")

	if result.Wait() && result.Error() != nil {
		fmt.Println("send message failed")
	} else {
		fmt.Println("send message success")
	}
	time.Sleep(1 * time.Hour)
}
