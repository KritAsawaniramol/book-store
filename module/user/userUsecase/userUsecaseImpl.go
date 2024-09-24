package userUsecase

import (
	"errors"
	"log"

	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"golang.org/x/crypto/bcrypt"
)

type userUsecaseImpl struct {
	userRepository userRepository.UserRepository
}

// Register implements UserUsecase.
func (u *userUsecaseImpl) Register(registReq *user.UserRegisterReq) (uint, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error: Register: %s\n", err.Error())
		return 0, errors.New("error: failed to hash password")
	}

	userID, err := u.userRepository.CreateOneUser(&user.User{
		Username: registReq.Username,
		Password: string(hashedPassword),
		RoleID:   2,
		Coin:     0,
	})
	if err != nil {
		return 0, err
	}
	return userID, err
}

func NewUserUsecaseImpl(userRepository userRepository.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepository: userRepository,
	}
}
