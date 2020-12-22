package iot

import (
	"bufio"
	"github.com/golang/glog"
	"io"
	"os"
	"runtime"
	"strings"
)

func OsName() string {
	return runtime.GOOS
}

func SdkInfo() map[string]string {
	f, err := os.Open("sdk_info")
	if err != nil {
		glog.Warning("read sdk info failed")
		return map[string]string{}
	}

	// 文件很小
	info := make(map[string]string)
	buf := bufio.NewReader(f)
	for {
		b, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			glog.Warningf("read sdk info failed or end")
			break
		}
		line := string(b)
		if len(line) != 0 {
			parts := strings.Split(line, "=")
			info[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
	}

	return info
}
