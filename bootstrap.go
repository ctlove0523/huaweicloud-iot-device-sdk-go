package iot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"sync"
	"time"
)

const (
	bsServer   = "tls://iot-bs.cn-north-4.myhuaweicloud.com:8883"
	bsServerCa = "-----BEGIN CERTIFICATE-----\n" +
		"MIIETjCCAzagAwIBAgINAe5fIh38YjvUMzqFVzANBgkqhkiG9w0BAQsFADBMMSAw\n" +
		"HgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMzETMBEGA1UEChMKR2xvYmFs\n" +
		"U2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjAeFw0xODExMjEwMDAwMDBaFw0yODEx\n" +
		"MjEwMDAwMDBaMFAxCzAJBgNVBAYTAkJFMRkwFwYDVQQKExBHbG9iYWxTaWduIG52\n" +
		"LXNhMSYwJAYDVQQDEx1HbG9iYWxTaWduIFJTQSBPViBTU0wgQ0EgMjAxODCCASIw\n" +
		"DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKdaydUMGCEAI9WXD+uu3Vxoa2uP\n" +
		"UGATeoHLl+6OimGUSyZ59gSnKvuk2la77qCk8HuKf1UfR5NhDW5xUTolJAgvjOH3\n" +
		"idaSz6+zpz8w7bXfIa7+9UQX/dhj2S/TgVprX9NHsKzyqzskeU8fxy7quRU6fBhM\n" +
		"abO1IFkJXinDY+YuRluqlJBJDrnw9UqhCS98NE3QvADFBlV5Bs6i0BDxSEPouVq1\n" +
		"lVW9MdIbPYa+oewNEtssmSStR8JvA+Z6cLVwzM0nLKWMjsIYPJLJLnNvBhBWk0Cq\n" +
		"o8VS++XFBdZpaFwGue5RieGKDkFNm5KQConpFmvv73W+eka440eKHRwup08CAwEA\n" +
		"AaOCASkwggElMA4GA1UdDwEB/wQEAwIBhjASBgNVHRMBAf8ECDAGAQH/AgEAMB0G\n" +
		"A1UdDgQWBBT473/yzXhnqN5vjySNiPGHAwKz6zAfBgNVHSMEGDAWgBSP8Et/qC5F\n" +
		"JK5NUPpjmove4t0bvDA+BggrBgEFBQcBAQQyMDAwLgYIKwYBBQUHMAGGImh0dHA6\n" +
		"Ly9vY3NwMi5nbG9iYWxzaWduLmNvbS9yb290cjMwNgYDVR0fBC8wLTAroCmgJ4Yl\n" +
		"aHR0cDovL2NybC5nbG9iYWxzaWduLmNvbS9yb290LXIzLmNybDBHBgNVHSAEQDA+\n" +
		"MDwGBFUdIAAwNDAyBggrBgEFBQcCARYmaHR0cHM6Ly93d3cuZ2xvYmFsc2lnbi5j\n" +
		"b20vcmVwb3NpdG9yeS8wDQYJKoZIhvcNAQELBQADggEBAJmQyC1fQorUC2bbmANz\n" +
		"EdSIhlIoU4r7rd/9c446ZwTbw1MUcBQJfMPg+NccmBqixD7b6QDjynCy8SIwIVbb\n" +
		"0615XoFYC20UgDX1b10d65pHBf9ZjQCxQNqQmJYaumxtf4z1s4DfjGRzNpZ5eWl0\n" +
		"6r/4ngGPoJVpjemEuunl1Ig423g7mNA2eymw0lIYkN5SQwCuaifIFJ6GlazhgDEw\n" +
		"fpolu4usBCOmmQDo8dIm7A9+O4orkjgTHY+GzYZSR+Y0fFukAj6KYXwidlNalFMz\n" +
		"hriSqHKvoflShx8xpfywgVcvzfTO3PYkz6fiNJBonf6q8amaEsybwMbDqKWwIX7eSPY=\n" +
		"-----END CERTIFICATE-----"
)

type BootstrapClient interface {
	Boot() string
	Close()
}

func NewBootstrapClient(id, password string) (BootstrapClient, error) {
	client := &bsClient{
		id:          id,
		password:    password,
		iotdaServer: newResult(),
	}

	res, err := client.init()
	if res {
		return client, nil
	}

	return nil, err
}

type bsClient struct {
	id          string
	password    string
	client      mqtt.Client // 使用的MQTT客户端
	iotdaServer *Result     // 设备接入平台地址
}

func (bs *bsClient) init() (bool, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(bsServer)
	options.SetClientID(CreateMqttClientId(bs.id))
	options.SetUsername(bs.id)
	options.SetPassword(hmacSha256(bs.password, timeStamp()))
	options.SetKeepAlive(250 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)

	ca := []byte(bsServerCa)
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
		glog.Warningf("device %s create bootstrap client failed,error = %v", bs.id, token.Error())
		return false, token.Error()
	}

	downTopic := fmt.Sprintf("$oc/devices/%s/sys/bootstrap/down", bs.id)
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
	upTopic := fmt.Sprintf("$oc/devices/%s/sys/bootstrap/up", bs.id)
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

func (b *Result) Wait() bool {
	<-b.Flag
	return true
}

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
