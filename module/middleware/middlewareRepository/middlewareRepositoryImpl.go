package middlewareRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfPb"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"gorm.io/gorm"
)

type middlewareRepositoryImpl struct {
	db *gorm.DB
}

// BookShelfSearch implements MiddlewareRepository.
func (m *middlewareRepositoryImpl) BookShelfSearch(grpcUrl string, userID uint, bookID uint) error {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchUserShelfRes, err := conn.Shelf().SearchUserShelf(ctx, &shelfPb.SearchUserShelfReq{
		UserId: uint64(userID),
		BookId: uint64(bookID),
	})
	if err != nil {
		log.Printf("error: BookShelfSearch: %s\n", err.Error())
		return errors.New("error: failed to search your book in your shelf")
	}

	if searchUserShelfRes.IsValid == false {
		return errors.New("error: you don't have permission to read this book")
	}
	return nil
}

// AccessTokenSearch implements MiddlewareRepository.
func (m *middlewareRepositoryImpl) AccessTokenSearch(grpcUrl string, accessToken string) error {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	accessTokenSearchRes, err := conn.Auth().AccessTokenSearch(ctx, &authPb.AccessTokenSearchReq{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Printf("error: AccessTokenSearch: %s\n", err.Error())
		return errors.New("error: access token is invalid")
	}

	if accessTokenSearchRes == nil {
		return errors.New("error: access token is invalid")
	}

	if !accessTokenSearchRes.IsValid {
		return errors.New("error: access token is invalid")
	}

	return nil
}

func NewMiddlewareRepositoryImpl(db *gorm.DB) MiddlewareRepository {
	return &middlewareRepositoryImpl{db: db}
}
