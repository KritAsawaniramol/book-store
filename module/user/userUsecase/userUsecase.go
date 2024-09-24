package userUsecase

import "github.com/kritAsawaniramol/book-store/module/user"

type UserUsecase interface {
	Register(registReq *user.UserRegisterReq) (uint, error)
}
