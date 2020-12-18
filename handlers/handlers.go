package handlers

type IotCommandHandler func(IotCommand) bool

type IotCommand struct {
	ObjectDeviceId string `json:"object_device_id"`
	ServiceId string `json:"service_id""`
	CommandName string  `json:"command_name"`
	Paras interface{} `json:"paras"`
}

type IotCommandResponse struct {
	ResultCode   byte      `json:"result_code"`
	ResponseName string      `json:"response_name"`
	Paras         interface{} `json:"paras"`
}

func SuccessIotCommandResponse() IotCommandResponse {
	return IotCommandResponse{
		ResultCode: 0,
	}
}

func FailedIotCommandResponse() IotCommandResponse {
	return IotCommandResponse{
		ResultCode: 1,
	}
}
