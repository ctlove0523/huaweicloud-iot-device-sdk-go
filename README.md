# huaweicloud-iot-device-sdk-go

huaweicloud-iot-device-sdk-go提供设备接入华为云IoT物联网平台的Go版本的SDK，提供设备和平台之间通讯能力，以及设备服务、网关服务、OTA等高级服务。IoT设备开发者使用SDK可以大大简化开发复杂度，快速的接入平台。

支持如下功能：

* [设备连接鉴权](#设备连接鉴权)

* [设备命令](#设备命令) 

* [设备消息](#设备消息)

* [设备属性](#设备属性) 

* [文件上传/下载管理](#文件上传/下载管理)

  

## 版本说明

当前稳定版本：0.0.1，[下载地址](https://github.com/ctlove0523/huaweicloud-iot-device-sdk-go/releases/tag/v0.0.1)

开发中版本：0.0.2-alpha

## 安装和构建

安装和构建的过程取决于你是使用go的 [modules](https://golang.org/ref/mod)(推荐) 还是还是`GOPATH`

### Modules

如果你使用 [modules](https://golang.org/ref/mod) 只需要导入包"github.com/ctlove0523/huaweicloud-iot-device-sdk-go"即可使用。当你使用go
build命令构建项目时，依赖的包会自动被下载。注意使用go
build命令构建时会自动下载最新版本，最新版本还没有达到release的标准可能存在一些尚未修复的bug。如果想使用稳定的发布版本可以从[release](https://github.com/ctlove0523/huaweicloud-iot-device-sdk-go/releases)
获取最新稳定的版本号，并在go.mod文件中指定版本号。

~~~go
module example

go 1.15

require github.com/ctlove0523/huaweicloud-iot-device-sdk-go v0.0.1-alpha
~~~

### GOPATH

如果你使用GOPATH，下面的一条命令即可实现安装

~~~go
go get github.com/ctlove0523/huaweicloud-iot-device-sdk-go
~~~

## 使用API

### 设备连接鉴权

1、首先，在华为云IoT平台创建一个设备，设备的信息如下：

设备ID：5fdb75cccbfe2f02ce81d4bf_go-mqtt

设备密钥：123456789

2、使用SDK创建一个Device对象，并初始化Device。

~~~go
import (
	"fmt"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"time"
)

func main() {
	// 创建一个设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()
	if device.IsConnected() {
		fmt.Println("device connect huawei iot platform success")
	} else {
		fmt.Println("device connect huawei iot platform failed")
	}
}
~~~

> iot-mqtts.cn-north-4.myhuaweicloud.com为华为IoT平台（基础班）在华为云北京四的访问端点，如果你购买了标准版或企业版，请将iot-mqtts.cn-north-4.myhuaweicloud.com更换为对应的MQTT协议接入端点。

### 设备命令

1、首先，在华为云IoT平台创建一个设备，设备的信息如下：

设备ID：5fdb75cccbfe2f02ce81d4bf_go-mqtt

设备密钥：123456789

2、使用SDK创建一个Device对象，并初始化Device。

~~~go
// 创建一个设备并初始化
device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
device.Init()
if device.IsConnected() {
	fmt.Println("device connect huawei iot platform success")
} else {
	fmt.Println("device connect huawei iot platform failed")
}
~~~

3、注册命令处理handler，支持注册多个handler并且按照注册的顺序回调

~~~go
// 添加用于处理平台下发命令的callback
device.AddCommandHandler(func(command iot.Command) bool {
	fmt.Println("First command handler begin to process command.")
	return true
})

device.AddCommandHandler(func(command iot.Command) bool {
	fmt.Println("Second command handler begin to process command.")
	return true
})
~~~

4、通过应用侧API向设备下发一个命令，可以看到程序输出如下：

~~~
device connect huawei iot platform success
First command handler begin to process command.
Second command handler begin to process command.
~~~

#### 完整样例

~~~go
import (
	"fmt"
	"github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"time"
)

// 处理平台下发的同步命令
func main() {
	// 创建一个设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()
	if device.IsConnected() {
		fmt.Println("device connect huawei iot platform success")
	} else {
		fmt.Println("device connect huawei iot platform failed")
	}

	// 添加用于处理平台下发命令的callback
	device.AddCommandHandler(func(command iot.Command) bool {
		fmt.Println("First command handler begin to process command.")
		return true
	})

	device.AddCommandHandler(func(command iot.Command) bool {
		fmt.Println("Second command handler begin to process command.")
		return true
	})
	time.Sleep(1 * time.Minute)
}
~~~

> 设备支持的命令定义在产品中

### 设备消息

1、首先，在华为云IoT平台创建一个设备，设备的信息如下：

设备ID：5fdb75cccbfe2f02ce81d4bf_go-mqtt

设备密钥：123456789

2、使用SDK创建一个Device对象，并初始化Device。

~~~go
device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()
~~~

#### 设备消息上报

~~~go
message := iot.Message{
	ObjectDeviceId: uuid.NewV4().String(),
	Name:           "Fist send message to platform",
	Id:             uuid.NewV4().String(),
	Content:        "Hello Huawei IoT Platform",
}
device.SendMessage(message)
~~~

#### 平台消息下发

接收平台下发的消息，只需注册消息处理handler，支持注册多个handler并按照注册顺序回调。

~~~go
// 注册平台下发消息的callback，当收到平台下发的消息时，调用此callback.
// 支持注册多个callback，并且按照注册顺序调用
device.AddMessageHandler(func(message iot.Message) bool {
	fmt.Println("first handler called" + iot.Interface2JsonString(message))
	return true
})

device.AddMessageHandler(func(message iot.Message) bool {
	fmt.Println("second handler called" + iot.Interface2JsonString(message))
	return true
})
~~~

#### 完整样例

~~~go
import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

func main() {
	// 创建一个设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()

	// 注册平台下发消息的callback，当收到平台下发的消息时，调用此callback.
	// 支持注册多个callback，并且按照注册顺序调用
	device.AddMessageHandler(func(message iot.Message) bool {
		fmt.Println("first handler called" + iot.Interface2JsonString(message))
		return true
	})

	device.AddMessageHandler(func(message iot.Message) bool {
		fmt.Println("second handler called" + iot.Interface2JsonString(message))
		return true
	})

	//向平台发送消息
	message := iot.Message{
		ObjectDeviceId: uuid.NewV4().String(),
		Name:           "Fist send message to platform",
		Id:             uuid.NewV4().String(),
		Content:        "Hello Huawei IoT Platform",
	}
	device.SendMessage(message)
	time.Sleep(2 * time.Minute)

}
~~~

### 设备属性

1、首先，在华为云IoT平台创建一个设备，并在该设备下创建3个子设备，设备及子设备的信息如下：

设备ID：5fdb75cccbfe2f02ce81d4bf_go-mqtt

设备密钥：123456789

子设备ID：5fdb75cccbfe2f02ce81d4bf_sub-device-1

子设备ID：5fdb75cccbfe2f02ce81d4bf_sub-device-2

子设备ID：5fdb75cccbfe2f02ce81d4bf_sub-device-3

2、使用SDK创建一个Device对象，并初始化Device。

~~~go
// 创建设备并初始化
device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
device.Init()
fmt.Printf("device connected: %v\n", device.IsConnected())
~~~

#### 设备属性上报

使用`ReportProperties(properties ServiceProperty) bool` 上报设备属性

~~~go
// 设备上报属性
props := iot.ServicePropertyEntry{
	ServiceId: "value",
	EventTime: iot.DataCollectionTime(),
	Properties: DemoProperties{
		Value:   "chen tong",
		MsgType: "23",
	},
}

var content []iot.ServicePropertyEntry
content = append(content, props)
services := iot.ServiceProperty{
	Services: content,
}
device.ReportProperties(services)
~~~

#### 网关批量设备属性上报

使用`BatchReportSubDevicesProperties(service DevicesService)` 实现网关批量设备属性上报

~~~go
// 批量上报子设备属性
subDevice1 := iot.DeviceService{
	DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-1",
	Services: content,
}
subDevice2 := iot.DeviceService{
	DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-2",
	Services: content,
}

subDevice3 := iot.DeviceService{
	DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-3",
	Services: content,
}

var devices []iot.DeviceService
devices = append(devices, subDevice1, subDevice2, subDevice3)

device.BatchReportSubDevicesProperties(iot.DevicesService{
	Devices: devices,
})
~~~

#### 平台设置设备属性

使用`AddPropertiesSetHandler(handler DevicePropertiesSetHandler)` 注册平台设置设备属性handler，当接收到平台的命令时SDK回调。

~~~go
// 注册平台设置属性callback,当应用通过API设置设备属性时，会调用此callback，支持注册多个callback
device.AddPropertiesSetHandler(func(propertiesSetRequest iot.DevicePropertyDownRequest) bool {
	fmt.Println("I get property set command")
	fmt.Printf("request is %s", iot.Interface2JsonString(propertiesSetRequest))
	return true
})
~~~

#### 平台查询设备属性

使用`SetPropertyQueryHandler(handler DevicePropertyQueryHandler)`注册平台查询设备属性handler，当接收到平台的查询请求时SDK回调。

~~~go
// 注册平台查询设备属性callback，当平台查询设备属性时此callback被调用，仅支持设置一个callback
device.SetPropertyQueryHandler(func(query iot.DevicePropertyQueryRequest) iot.ServicePropertyEntry {
	return iot.ServicePropertyEntry{
		ServiceId: "value",
		Properties: DemoProperties{
			Value:   "QUERY RESPONSE",
			MsgType: "query property",
		},
		EventTime: "2020-12-19 02:23:24",
	}
})
~~~

#### 设备侧获取平台的设备影子数据

使用`QueryDeviceShadow(query DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler)`
可以查询平台的设备影子数据，当接收到平台的响应后SDK自动回调`DevicePropertyQueryResponseHandler`。

~~~go
// 设备查询设备影子数据
device.QueryDeviceShadow(iot.DevicePropertyQueryRequest{
	ServiceId: "value",
}, func(response iot.DevicePropertyQueryResponse) {
	fmt.Printf("query device shadow success.\n,device shadow data is %s\n", iot.Interface2JsonString(response))
})
~~~

#### 完整样例

~~~go
import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"time"
)

func main() {
	// 创建设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
	device.Init()
	fmt.Printf("device connected: %v\n", device.IsConnected())

	// 注册平台设置属性callback,当应用通过API设置设备属性时，会调用此callback，支持注册多个callback
	device.AddPropertiesSetHandler(func(propertiesSetRequest iot.DevicePropertyDownRequest) bool {
		fmt.Println("I get property set command")
		fmt.Printf("request is %s", iot.Interface2JsonString(propertiesSetRequest))
		return true
	})

	// 注册平台查询设备属性callback，当平台查询设备属性时此callback被调用，仅支持设置一个callback
	device.SetPropertyQueryHandler(func(query iot.DevicePropertyQueryRequest) iot.ServicePropertyEntry {
		return iot.ServicePropertyEntry{
			ServiceId: "value",
			Properties: DemoProperties{
				Value:   "QUERY RESPONSE",
				MsgType: "query property",
			},
			EventTime: "2020-12-19 02:23:24",
		}
	})

	// 设备上报属性
	props := iot.ServicePropertyEntry{
		ServiceId: "value",
		EventTime: iot.DataCollectionTime(),
		Properties: DemoProperties{
			Value:   "chen tong",
			MsgType: "23",
		},
	}

	var content []iot.ServicePropertyEntry
	content = append(content, props)
	services := iot.ServiceProperty{
		Services: content,
	}
	device.ReportProperties(services)

	// 设备查询设备影子数据
	device.QueryDeviceShadow(iot.DevicePropertyQueryRequest{
		ServiceId: "value",
	}, func(response iot.DevicePropertyQueryResponse) {
		fmt.Printf("query device shadow success.\n,device shadow data is %s\n", iot.Interface2JsonString(response))
	})

	// 批量上报子设备属性
	subDevice1 := iot.DeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-1",
		Services: content,
	}
	subDevice2 := iot.DeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-2",
		Services: content,
	}

	subDevice3 := iot.DeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-3",
		Services: content,
	}

	var devices []iot.DeviceService
	devices = append(devices, subDevice1, subDevice2, subDevice3)

	device.BatchReportSubDevicesProperties(iot.DevicesService{
		Devices: devices,
	})
	time.Sleep(1 * time.Minute)
}

type DemoProperties struct {
	Value   string `json:"value"`
	MsgType string `json:"msgType"`
}
~~~

### 文件上传/下载管理

#### 文件上传

~~~go
device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "xxx", "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883")
device.Init()

device.UploadFile("D/software/mqttfx/chentong.txt")
~~~



## 报告bugs

如果你在使用过程中遇到任何问题或bugs，请通过issue的方式上报问题或bug，我们将会在第一时间内答复。上报问题或bugs时请尽量提供以下内容：

* 使用的版本
* 使用场景
* 重现问题或bug的样例代码
* 错误信息
* ······

## 贡献

该项目欢迎来自所有人的pull request。