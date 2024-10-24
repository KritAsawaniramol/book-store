package shelfHandler

import (
	"context"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfPb"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
)

type shelfGrpcHandlerImpl struct {
	cfg          *config.Config
	shelfUsecase shelfUsecase.ShelfUsecase
	shelfPb.UnimplementedShelfGrpcServiceServer
}

// SearchUserShelf implements shelfPb.ShelfGrpcServiceServer.
func (s *shelfGrpcHandlerImpl) SearchUserShelf(ctx context.Context, req *shelfPb.SearchUserShelfReq) (*shelfPb.SearchUserShelfRes, error) {
	myShelves, err := s.shelfUsecase.GetMyShelves(s.cfg, uint(req.UserId), uint(req.BookId))
	if err != nil {
		return &shelfPb.SearchUserShelfRes{IsValid: false}, err
	}
	if myShelves.Shelves == nil || len(myShelves.Shelves) < 1 {
		return &shelfPb.SearchUserShelfRes{IsValid: false}, nil
	}
	return &shelfPb.SearchUserShelfRes{IsValid: true}, nil
}

func NewShelfGrpcHandlerImpl(cfg *config.Config, shelfUsecase shelfUsecase.ShelfUsecase) shelfPb.ShelfGrpcServiceServer {
	return &shelfGrpcHandlerImpl{cfg: cfg, shelfUsecase: shelfUsecase}
}
