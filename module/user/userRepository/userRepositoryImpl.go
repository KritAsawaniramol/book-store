package userRepository

import (
	"errors"
	"log"
	"strings"

	"github.com/kritAsawaniramol/book-store/module/user"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// GetOneUser implements UserRepository.
func (u *userRepositoryImpl) GetOneUser(in *user.User) (*user.User, error) {
	user := &user.User{}
	if err := u.db.Where(&in).First(&user).Error; err != nil {
		log.Printf("error: GetOneUser: %s\n", err.Error())
		return nil, errors.New("error: user not found")
	}
	return user, nil
}

// createOneUser implements UserRepository.
func (u *userRepositoryImpl) CreateOneUser(in *user.User) (uint, error) {
	if err := u.db.Create(in).Error; err != nil {
		log.Printf("error: CreateOneUser: %s\n", err.Error())
		if strings.HasSuffix(err.Error(), "(SQLSTATE 23505)") {
			return 0, errors.New("error: username already in use")
		}
		return 0, errors.New("error: create user failed")
	}
	return in.ID, nil
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}
