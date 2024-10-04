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
		UserID uint  `gorm:"not null"`
		Amount int64 `gorm:"not null"`
	}
)
