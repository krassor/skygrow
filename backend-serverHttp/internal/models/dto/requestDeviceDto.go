package dto

type RequestDeviceDto struct {
	DeviceVendor    string `json:"deviceVendor"`
	DeviceName      string `json:"deviceName"`
	DeviceSchema    string `json:"deviceSchema"`
	DeviceIpAddress string `json:"deviceIpAddress"`
	DevicePort      string `json:"devicePort"`
}
