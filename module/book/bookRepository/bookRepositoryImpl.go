package bookRepository

import "gorm.io/gorm"

type bookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepositoryImpl(db *gorm.DB) BookRepository {
	return &bookRepositoryImpl{db: db}
}
