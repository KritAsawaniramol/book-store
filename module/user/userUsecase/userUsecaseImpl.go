package userUsecase

import (
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
)

type userUsecaseImpl struct {
	userRepository userRepository.UserRepository
}

// Register implements UserUsecase.
func (u *userUsecaseImpl) Register(registReq *user.UserRegisterReq) (uint, error) {
	userID, err := u.userRepository.CreateOneUser(&user.User{
		Username: registReq.Username,
		Password: registReq.Password,
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
