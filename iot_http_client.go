package iot

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

// 仅用于设备上传文件
type HttpClient interface {
	UploadFile(filename, uri string) bool
	DownloadFile(filename, uri string) bool
}

type httpClient struct {
	client *http.Client
}

func (client *httpClient) DownloadFile(fileName, downloadUrl string) bool {
	fmt.Println(downloadUrl)
	fileName = SmartFileName(fileName)
	out, err := os.Create(fileName)
	if err != nil {
		glog.Errorf("create file in os failed ,file name %s", fileName)
		return false
	}
	defer out.Close()

	bodyBuffer := &bytes.Buffer{}
	req, err := http.NewRequest("GET", downloadUrl, bodyBuffer)
	if err != nil {
		glog.Errorf("create request filed %v", err)
		return false
	}

	req.Header.Add("Content-Type", "text/plain")

	originalUri, err := url.ParseRequestURI(downloadUrl)
	if err != nil {
		glog.Errorf("parse request uri failed %v", err)
		return false
	}
	req.Header.Add("Host", originalUri.Host)

	res, err := client.client.Do(req)
	if err != nil {
		glog.Errorf("down load file error")
		return false
	}
	buf := make([]byte, 1024)
	for {
		size, _ := res.Body.Read(buf)
		if size == 0 {
			break
		} else {
			_, _ = out.Write(buf[:size])
		}
	}

	return true
}

func (client *httpClient) UploadFile(filename, uri string) bool {
	filename = SmartFileName(filename)
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, err := bodyWriter.CreateFormFile("files", filename)
	if err != nil {
		glog.Errorf("create form file failed %v", err)
		return false
	}

	file, err := os.Open(filename)
	if err != nil {
		glog.Errorf("open file failed %v", err)
		return false
	}

	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		glog.Errorf("copy file to writer failed %v", err)
	}

	//contentType := bodyWriter.FormDataContentType()
	defer bodyWriter.Close()

	req, err := http.NewRequest("PUT", uri, bodyBuffer)
	if err != nil {
		glog.Errorf("create request filed %v", err)
	}

	req.Header.Add("Content-Type", "text/plain")

	originalUri, err := url.ParseRequestURI(uri)
	if err != nil {
		glog.Errorf("parse request uri failed %v", err)
	}
	req.Header.Add("Host", originalUri.Host)
	resp, _ := client.client.Do(req)

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)

	return nil == err
}

func CreateHttpClient() HttpClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	innerClient := &http.Client{Transport: tr}

	httpClient := &httpClient{
		client: innerClient,
	}

	return httpClient

}
