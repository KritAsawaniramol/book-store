package middlewareRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"gorm.io/gorm"
)

type middlewareRepositoryImpl struct {
	db *gorm.DB
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
