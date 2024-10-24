package user

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model
		Username         string `gorm:"unique;not null"`
		Password         string `gorm:"not null"`
		RoleID           uint   `gorm:"not null"`
		UserTransactions []UserTransactions
	}

	Role struct {
		gorm.Model
		RoleTitle string `gorm:"not null"`
	}

	UserTransactions struct {
		gorm.Model
		UserID       uint  `gorm:"not null"`
		Amount       int64 `gorm:"not null"`
		Note         string
		TopUpOrderID *uint
		TopUpOrder   *TopUpOrder
	}

	TopUpOrder struct {
		gorm.Model
		UserID    uint   `gorm:"not null"`
		Amount    int64  `gorm:"not null"`
		SessionID string `gorm:"not null"`
		Status    string `gorm:"not null"`
	}
)
