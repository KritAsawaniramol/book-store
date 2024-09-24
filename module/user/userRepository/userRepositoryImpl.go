package userRepository

import (
	"errors"
	"log"

	"github.com/kritAsawaniramol/book-store/module/user"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// createOneUser implements UserRepository.
func (u *userRepositoryImpl) CreateOneUser(in *user.User) (uint, error) {
	if err := u.db.Create(in).Error; err != nil {
		log.Printf("error: createOneUser: %s\n", err.Error())
		return 0, errors.New("create user failed")
	}
	return in.ID, nil
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}
