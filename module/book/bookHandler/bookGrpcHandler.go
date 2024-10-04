package bookHandler

import (
	"context"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/book/bookUsecase"
)

type (
	BookGrpcHandler interface {
	}

	bookGrpcHandlerImpl struct {
		cfg         *config.Config
		bookUsecase bookUsecase.BookUsecase
		bookPb.UnimplementedBookGrpcServiceServer
	}
)

// FindBooksInIds implements bookPb.BookGrpcServiceServer.
func (b *bookGrpcHandlerImpl) FindBooksInIds(ctx context.Context, req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error) {
	return b.bookUsecase.FindBookInIDs(req)
}

func NewBookGrpcHandlerImpl(cfg *config.Config, bookUsecase bookUsecase.BookUsecase) bookPb.BookGrpcServiceServer {
	return &bookGrpcHandlerImpl{cfg: cfg, bookUsecase: bookUsecase}
}
