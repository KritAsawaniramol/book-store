package shelfHandler

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
)

type (
	ShelfHttpHandler interface {
	}

	shelfHttpHandlerImpl struct {
		cfg          *config.Config
		shelfUsecase shelfUsecase.ShelfUsecase
	}
)

func NewShelfHttpHandlerImpl(cfg *config.Config, shelfUescase shelfUsecase.ShelfUsecase) ShelfHttpHandler {
	return &shelfHttpHandlerImpl{
		cfg:          cfg,
		shelfUsecase: shelfUescase,
	}
}
