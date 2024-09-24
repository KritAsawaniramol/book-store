package bookHandler

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
)

type (
	BookHttpHandler interface {
	}

	bookHttpHandlerImpl struct {
		cfg         *config.Config
		bookUsecase bookUsecase.BookUsecase
	}
)

func NewBookHttpHandlerImpl(cfg *config.Config, bookUsecase bookUsecase.BookUsecase) BookHttpHandler {
	return &bookHttpHandlerImpl{cfg: cfg, bookUsecase: bookUsecase}
}
