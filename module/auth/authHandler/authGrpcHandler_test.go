package authHandler

import (
	"context"
	"testing"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenSearch(t *testing.T) {
	cfg := config.LoadConfig("../../../env/test/.env")
	usecaseMock := authUsecase.NewAuthUsecaseMock()
	grpcHandler := NewAuthGrpcHandlerImpl(usecaseMock)
	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(cfg.Grpc.AuthUrl)
		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)
		grpcServer.Serve(listener)
	}()

	conn, err := grpccon.NewGrpcClient(cfg.Grpc.AuthUrl)
	if err != nil {
		panic(err)
	}

	req := &authPb.AccessTokenSearchReq{
		AccessToken: "access_token",
	}
	t.Run("success", func(t *testing.T) {
		usecaseMock.On("AccessTokenSearch", req.AccessToken).Return(&authPb.AccessTokenSearchRes{IsValid: true}, nil)
		res, err := conn.Auth().AccessTokenSearch(context.Background(), req)
		assert.True(t, res.IsValid)
		assert.NoError(t, err)
	})
}
