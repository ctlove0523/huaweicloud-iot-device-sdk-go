package iot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type BootstrapClientConfig struct {
	id       string // 设备Id，平台又称为deviceId
	password string // 设备密码
	server   string // 设备发放平台地址：tls://iot-bs.cn-north-4.myhuaweicloud.com:8883
	caPath   string // 设备发放平台CA证书
}

func NewBootstrapClientConfig() *BootstrapClientConfig {
	return &BootstrapClientConfig{}
}

func (conf *BootstrapClientConfig) WithId(id string) *BootstrapClientConfig {
	conf.id = id
	return conf
}

func (conf *BootstrapClientConfig) WithPassword(password string) *BootstrapClientConfig {
	conf.password = password
	return conf
}

func (conf *BootstrapClientConfig) WithServer(server string) *BootstrapClientConfig {
	conf.server = server
	return conf
}

func (conf *BootstrapClientConfig) WithCaPath(caPath string) *BootstrapClientConfig {
	conf.caPath = caPath
	return conf
}

type BootstrapClient interface {
	Boot() string
	Close()
}

func NewBootstrapClient(conf *BootstrapClientConfig) (BootstrapClient, error) {
	client := &bsClient{
		conf:        conf,
		iotdaServer: newResult(),
	}

	res, err := client.init()
	if res {
		return client, nil
	}

	return nil, err
}

type bsClient struct {
	conf        *BootstrapClientConfig
	client      mqtt.Client // 使用的MQTT客户端
	iotdaServer *Result     // 设备接入平台地址
}

func (bs *bsClient) init() (bool, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(bs.conf.server)
	options.SetClientID(createClientId(bs.conf.id))
	options.SetUsername(bs.conf.id)
	options.SetPassword(hmacSha256(bs.conf.password, timeStamp()))
	options.SetKeepAlive(250 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)

	ca, err := ioutil.ReadFile(bs.conf.caPath)
	if err != nil {
		glog.Error("load server ca failed\n")
		return false, err
	}
	serverCaPool := x509.NewCertPool()
	serverCaPool.AppendCertsFromPEM(ca)

	tlsConfig := &tls.Config{
		RootCAs:            serverCaPool,
		InsecureSkipVerify: true,
		MaxVersion:         tls.VersionTLS12,
		MinVersion:         tls.VersionTLS12,
	}
	options.SetTLSConfig(tlsConfig)

	bs.client = mqtt.NewClient(options)
	if token := bs.client.Connect(); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s create bootstrap client failed,error = %v", bs.conf.id, token.Error())
		return false, token.Error()
	}

	downTopic := fmt.Sprintf("$oc/devices/%s/sys/bootstrap/down", bs.conf.id)
	subRes := bs.client.Subscribe(downTopic, 0, func(client mqtt.Client, message mqtt.Message) {
		go func() {
			fmt.Println("get message from bs server")
			serverResponse := &serverResponse{}
			err := json.Unmarshal(message.Payload(), serverResponse)
			if err != nil {
				fmt.Println(err)
				bs.iotdaServer.CompleteError(err)
			} else {
				bs.iotdaServer.Complete(serverResponse.Address)
			}
		}()
	})
	if subRes.Wait() && subRes.Error() != nil {
		fmt.Printf("sub topic %s failed,error is %s\n", downTopic, subRes.Error())
		return false, subRes.Error()
	} else {
		fmt.Printf("sub topic %s success\n", downTopic)
	}

	return true, nil
}

func (bs *bsClient) Boot() string {
	upTopic := fmt.Sprintf("$oc/devices/%s/sys/bootstrap/up", bs.conf.id)
	pubRes := bs.client.Publish(upTopic, 0, false, "")
	if pubRes.Wait() && pubRes.Error() != nil {
		fmt.Println(pubRes.Error())
		return ""
	}

	bs.iotdaServer.Wait()
	return "tls://" + bs.iotdaServer.Value()
}

func (bs *bsClient) Close() {
	bs.client.Disconnect(1000)
}

type serverResponse struct {
	Address string `json:"address"`
}

func createClientId(deviceId string) string {
	segments := make([]string, 4)
	segments[0] = deviceId
	segments[1] = "0"
	segments[2] = "0"
	segments[3] = timeStamp()

	return strings.Join(segments, "_")
}

type Result struct {
	Flag chan int

	err  error
	mErr sync.RWMutex

	res  string
	mRes sync.RWMutex
}

func (b *Result) Value() string {
	b.mRes.RLock()
	defer b.mRes.RUnlock()
	return b.res
}

// Wait implements the Token Wait method.
func (b *Result) Wait() bool {
	<-b.Flag
	return true
}

// WaitTimeout implements the Token WaitTimeout method.
func (b *Result) WaitTimeout(d time.Duration) bool {
	timer := time.NewTimer(d)
	select {
	case <-b.Flag:
		if !timer.Stop() {
			<-timer.C
		}
		return true
	case <-timer.C:
	}

	return false
}

func (b *Result) Complete(res string) {
	b.mRes.Lock()
	defer b.mRes.Unlock()
	b.res = res
	b.Flag <- 1
}

func (b *Result) CompleteError(err error) {
	b.mErr.Lock()
	defer b.mErr.Unlock()
	b.err = err
	b.Flag <- 1
}

func (b *Result) Error() error {
	b.mErr.RLock()
	defer b.mErr.RUnlock()
	return b.err
}

func newResult() *Result {
	return &Result{
		Flag: make(chan int),
	}
}
