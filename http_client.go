package iot

import (
	"github.com/go-resty/resty/v2"
	"github.com/golang/glog"
	"io/ioutil"
	"net/url"
	"os"
)

// 仅用于设备上传文件
type HttpClient interface {
	UploadFile(filename, uri string) bool
	DownloadFile(filename, uri string) bool
}

type httpClient struct {
	client *resty.Client
}

func (client *httpClient) DownloadFile(fileName, downloadUrl string) bool {
	glog.Infof("begin to download file %s, url = %s", fileName, downloadUrl)
	fileName = smartFileName(fileName)
	out, err := os.Create(fileName)
	if err != nil {
		glog.Errorf("create file in os failed ,file name %s", fileName)
		return false
	}

	originalUri, err := url.ParseRequestURI(downloadUrl)
	if err != nil {
		glog.Errorf("parse request uri failed %v", err)
		return false
	}

	resp, err := client.client.R().
		SetHeader("Content-Type", "text/plain").
		SetHeader("Host", originalUri.Host).
		Get(downloadUrl)

	if err != nil {
		glog.Errorf("download file request failed %v", err)
		return false
	}

	_, err = out.Write(resp.Body())
	if err != nil {
		glog.Errorf("write file failed")
		return false
	}

	return true
}

func (client *httpClient) UploadFile(filename, uri string) bool {
	filename = smartFileName(filename)
	fileBytes, err := ioutil.ReadFile(filename)

	if err != nil {
		glog.Errorf("read file failed %v", err)
		return false
	}

	originalUri, err := url.ParseRequestURI(uri)
	if err != nil {
		glog.Errorf("parse request uri failed %v", err)
		return false
	}

	resp, err := client.client.R().
		SetHeader("Content-Type", "text/plain").
		SetHeader("Host", originalUri.Host).
		SetBody(fileBytes).
		Put(uri)

	if err != nil {
		glog.Errorf("upload request failed %v", err)
	}

	return resp.StatusCode() == 200
}

func CreateHttpClient() HttpClient {
	client := resty.New()

	client.SetRetryCount(3)

	httpClient := &httpClient{
		client: client,
	}

	return httpClient

}
