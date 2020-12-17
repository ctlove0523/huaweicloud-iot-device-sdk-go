package iotdevice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

func TimeStamp() string {
	strFormatTime :=time.Now().Format("2006-01-02 15:04:05")
	strFormatTime =strings.ReplaceAll(strFormatTime,"-","")
	strFormatTime = strings.ReplaceAll(strFormatTime," ","")
	strFormatTime = strFormatTime[0:10]
	return strFormatTime
}

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
