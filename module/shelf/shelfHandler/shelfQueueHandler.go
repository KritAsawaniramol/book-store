package shelfHandler

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
)

type (
	ShelfQueueHandler interface {
	}

	shelfQueueHandlerImpl struct {
		cfg          *config.Config
		shelfUsecase shelfUsecase.ShelfUsecase
	}
)

func NewShelfQueueHandlerImpl(cfg *config.Config, shelfUescase shelfUsecase.ShelfUsecase) ShelfQueueHandler {
	return &shelfQueueHandlerImpl{
		cfg:          cfg,
		shelfUsecase: shelfUescase,
	}
}