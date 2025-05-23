package models

import (
	"time"
)

type AccountVerification struct {
	ID         uint           	`gorm:"primaryKey;autoIncrement"`
	UserId 		uint      		`gorm:"not null;"` 
	OTP      	string    		`gorm:"size:255;not null"`
	ExpiredAt time.Time 		`gorm:"autoCreateTime"`
	CreatedAt  time.Time      	`gorm:"autoCreateTime"`
}
