package authRepository

import "gorm.io/gorm"

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryImpl(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}


