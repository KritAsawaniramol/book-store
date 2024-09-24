package auth

import "gorm.io/gorm"

type (
	Credential struct {
		gorm.Model
		UserID       uint   `gorm:"not null"`
		RoleID       uint   `gorm:"not null"`
		AccessToken  string `gorm:"not null"`
		RefreshToken string `gorm:"not null"`
	}
)
