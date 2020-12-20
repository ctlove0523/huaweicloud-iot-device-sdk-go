package iot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

// 时间戳：为设备连接平台时的UTC时间，格式为YYYYMMDDHH，如UTC 时间2018/7/24 17:56:20 则应表示为2018072417。
func TimeStamp() string {
	strFormatTime := time.Now().Format("2006-01-02 15:04:05")
	strFormatTime = strings.ReplaceAll(strFormatTime, "-", "")
	strFormatTime = strings.ReplaceAll(strFormatTime, " ", "")
	strFormatTime = strFormatTime[0:10]
	return strFormatTime
}

// 设备采集数据UTC时间（格式：yyyyMMdd'T'HHmmss'Z'），如：20161219T114920Z。
//设备上报数据不带该参数或参数格式错误时，则数据上报时间以平台时间为准。
func DataCollectionTime() string {
	now := time.Now().UTC()
	return now.Format("20060102T150405Z")
}

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Interface2JsonString(v interface{}) string {
	byteData, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(byteData)
}

func GetTopicRequestId(topic string) string {
	return strings.Split(topic, "=")[1]
}

func FormatTopic(topic, deviceId string) string {
	return strings.ReplaceAll(topic, "{device_id}", deviceId)
}
