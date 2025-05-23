package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string
type Role string

const (
	RoleSuperUser   Role = "SUPER_USER"
	RoleUser        Role = "USER"
	RoleAdmin       Role = "ADMIN"
	RoleContributor Role = "CONTRIBUTOR"
	RoleAccountant  Role = "ACCOUNTANT"
	RoleManager     Role = "MANAGER"
	RoleEditor      Role = "EDITOR"
)

type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	Password      string    `gorm:"type:varchar(255);not null"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Fullname      string    `gorm:"type:varchar(200);not null"`
	ContactTitle  string    `gorm:"type:text;not null"`
	Mobile        string    `gorm:"type:varchar(20);not null"`
	Enabled       bool      `gorm:"type:boolean;default:false;not null"`
	Status        string    `gorm:"column:status;type:enum('PENDING','VERIFIED', 'SUSPEND', 'PLAN_EXPIRED');default:'PENDING'" json:"status"`
	Image         string    `gorm:"type:text;default:'';not null"`
	Views         string    `gorm:"type:text;default:'';not null"`
	FacebookLink  string    `gorm:"type:text;default:'';not null"`
	InstagramLink string    `gorm:"type:text;default:'';not null"`
	TwitterLink   string    `gorm:"type:text;default:'';not null"`
	LinkedInLink  string    `gorm:"type:text;default:'';not null"`
	Role          string    `gorm:"column:role;type:enum('USER', 'ADMIN', 'SUPER_USER', 'CONTRIBUTOR', 'ACCOUNTANT', 'MANAGER', 'EDITOR');default:'USER'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return nil
}
