package shelf

import "gorm.io/gorm"

type (
	Shelf struct {
		gorm.Model
		BooksID uint `gorm:"not null"`
		Users   uint `gorm:"not null"`
	}
)
