package userHandler

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
)

type UserQueueHandler interface {
}

type userQueueHandlerImpl struct {
	userUsecase userUsecase.UserUsecase
	cfg         *config.Config
}

func NewUserQueueHandler(userUsecase userUsecase.UserUsecase, cfg *config.Config) UserQueueHandler {
	return &userQueueHandlerImpl{userUsecase: userUsecase, cfg: cfg}
}
