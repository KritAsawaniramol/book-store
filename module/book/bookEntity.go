package book

import "gorm.io/gorm"

type (
	Books struct {
		gorm.Model
		Title         string
		Price         int64
		FileUrl       string
		CoverImageUrl string
		AuthorName    string
		Genres        []BooksGenres
	}

	BooksGenres struct {
		BooksID  uint   `gorm:"primaryKey"`
		Books    Books  `gorm:"foreignKey:BooksID;references:ID"`
		GenresID uint   `gorm:"primaryKey"`
		Genres   Genres `gorm:"foreignKey:GenresID;references:ID"`
	}

	Genres struct {
		gorm.Model
		GenresTitle string
		Books       []BooksGenres
	}
)
