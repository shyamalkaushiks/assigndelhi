package model

import "time"

type UserRegisterwe struct {
	ID           uint   `gorm:"primaryKey"`
	PhoneNumber  string `gorm:"unique;not null"`
	OTP          string `gorm:"size:6"`
	OTPExpiresAt time.Time
	DeviceID     string `gorm:"not null"`
	Token        string `json:"token"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
