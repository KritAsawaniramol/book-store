package book

import "gorm.io/gorm"

type (
	Books struct {
		gorm.Model
		Title              string
		Price              uint
		FilePath           string
		Description        string
		CoverImagePath     string
		AuthorName         string
		IsAvailableInStore bool
		Tags               []Tags `gorm:"many2many:books_tags;"`
	}

	Tags struct {
		gorm.Model
		Name  string
		Books []Books `gorm:"many2many:books_tags;"`
	}
)
