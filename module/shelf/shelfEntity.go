package shelf

import "gorm.io/gorm"

type (
	Shelves struct {
		gorm.Model
		UsersID uint `gorm:"not null"`
		BooksID uint `gorm:"not null"`
	}
)
