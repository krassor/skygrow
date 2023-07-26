package entities

import (
	"gorm.io/gorm"
)

type Devices struct {
	gorm.Model
	DeviceVendor    string `gorm:"column:deviceVendor"`
	DeviceName      string `gorm:"column:deviceName"`
	DeviceSchema    string `gorm:"column:deviceSchema"`
	DeviceIpAddress string `gorm:"column:deviceIpAddress"`
	DevicePort      string `gorm:"column:devicePort"`
	DeviceStatus    bool   `gorm:"column:deviceStatus;default:false"`
}

type Subscriber struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	ChatID   int64  `gorm:"column:chatid"`
	IsActive bool   `gorm:"column:isActive"`
}
