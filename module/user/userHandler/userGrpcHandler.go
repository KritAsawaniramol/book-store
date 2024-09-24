package userHandler

import "github.com/kritAsawaniramol/book-store/module/user/userUsecase"

type UserGrpcHandler interface {
}

type userGrpcHandlerImpl struct {
	userUsecase userUsecase.UserUsecase

}

func NewUserGrpcHandler(userUsecase userUsecase.UserUsecase) UserGrpcHandler {
	return &userGrpcHandlerImpl{userUsecase: userUsecase}
}
