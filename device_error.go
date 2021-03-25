package iot

type DeviceError struct {
	errorMsg string
}

func (err *DeviceError) Error() string {
	return err.errorMsg
}
