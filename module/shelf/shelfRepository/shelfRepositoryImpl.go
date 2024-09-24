package shelfRepository

import "gorm.io/gorm"

type shelfRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) ShelfRepository {
	return &shelfRepositoryImpl{
		db: db,
	}
}
