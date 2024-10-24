package userHandler

import (
	"context"

	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
)

type userGrpcHandler struct {
	userUsecase userUsecase.UserUsecase
	userPb.UnimplementedUserGrpcServiceServer
}

// FindOneUserProfileToRefresh implements userPb.UserGrpcServiceServer.
func (u *userGrpcHandler) FindOneUserProfileToRefresh(ctx context.Context, req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
	return u.userUsecase.FindOneUserByID(uint(req.UserId))
}

// FindUserProfileToLogin implements userPb.UserGrpcServiceServer.
func (u *userGrpcHandler) FindUserProfileToLogin(ctx context.Context, req *userPb.FindUserProfileToLoginReq) (*userPb.UserProfile, error) {
	return u.userUsecase.FindOneUserByUsernameAndPassword(req.Username, req.Password)
}

func NewUserGrpcHandler(userUsecase userUsecase.UserUsecase) userPb.UserGrpcServiceServer {
	return &userGrpcHandler{userUsecase: userUsecase}
}
