package book

import "gorm.io/gorm"

type (
	Books struct {
		gorm.Model
		Title          string
		Price          uint
		FilePath       string
		CoverImagePath string
		AuthorName     string
		Tags           []BooksTags
	}

	BooksTags struct {
		BooksID uint  `gorm:"primaryKey"`
		Books   Books `gorm:"foreignKey:BooksID;references:ID"`
		TagsID  uint  `gorm:"primaryKey"`
		Tags    Tags  `gorm:"foreignKey:TagsID;references:ID"`
	}

	Tags struct {
		gorm.Model
		Name  string `gorm:"unique"`
		Books []BooksTags
	}
)
