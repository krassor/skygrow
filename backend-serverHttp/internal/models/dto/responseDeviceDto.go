package dto

type ResponseDeviceDto struct {
	DeviceVendor    string `json:"deviceVendor"`
	DeviceName      string `json:"deviceModel"`
	DeviceSchema    string `json:"deviceSchema"`
	DeviceIpAddress string `json:"deviceIpAddress"`
	DevicePort      string `json:"devicePort"`
}

type ResponseDeviceParams struct {
	DeviceId uint `json:"deviceId"`
}
