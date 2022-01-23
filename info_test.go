package iot

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// 该测试用力仅能在Windows系统运行
func TestOsName(t *testing.T) {
	if !strings.Contains(OsName(), "windows") {
		t.Errorf(`OsName must be windwos`)
	}
}

func TestVersion(t *testing.T) {
	if SdkInfo()["sdk-version"] != "v2.0.0" {
		t.Errorf("sdk version must be v0.0.2")
	}

	if SdkInfo()["author"] != "chen tong" {
		t.Errorf("sdk author must be chen tong")
	}
}

func TestCreateFileUploadResultResponse(t *testing.T) {
	f, err := os.Open("sdk_info")
	if err != nil {
		fmt.Println(err.Error())
	}

	//建立缓冲区，把文件内容放到缓冲区中
	buf := bufio.NewReader(f)
	for {
		//遇到\n结束读取
		b, errR := buf.ReadBytes('\n')
		if errR != nil {
			if errR == io.EOF {
				break
			}
			fmt.Println(errR.Error())
		}
		fmt.Println(string(b))
	}
}
