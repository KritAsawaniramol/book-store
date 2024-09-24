package bookHandler

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
)

type (
	BookGrpcHandler interface {
	}

	bookGrpcHandlerImpl struct {
		cfg         *config.Config
		bookUsecase bookUsecase.BookUsecase
	}
)

func NewBookGrpcHandlerImpl(cfg *config.Config, bookUsecase bookUsecase.BookUsecase) BookGrpcHandler {
	return &bookGrpcHandlerImpl{cfg: cfg, bookUsecase: bookUsecase}
}
