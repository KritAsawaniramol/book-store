package authHandler

import (
	"context"

	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
)

type (
	authGrpcHandlerImpl struct {
		authUsecase authUsecase.AuthUsecase
		authPb.UnimplementedAuthGrpcServiceServer
	}
)

// AccessTokenSearch implements authPb.AuthGrpcServiceServer.
func (a *authGrpcHandlerImpl) AccessTokenSearch(ctx context.Context, req *authPb.AccessTokenSearchReq) (*authPb.AccessTokenSearchRes, error) {
	return a.authUsecase.AccessTokenSearch(req.AccessToken)
}

func NewAuthGrpcHandlerImpl(authUsecase authUsecase.AuthUsecase) authPb.AuthGrpcServiceServer {
	return &authGrpcHandlerImpl{
		authUsecase: authUsecase,
	}
}
