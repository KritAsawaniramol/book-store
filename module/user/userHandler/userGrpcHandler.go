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

// FindUserProfile implements userPb.UserGrpcServiceServer.
func (u *userGrpcHandler) FindUserProfile(ctx context.Context, req *userPb.FindUserProfileReq) (*userPb.UserProfile, error) {
	return u.userUsecase.FindOneUserByUsernameAndPassword(req.Username, req.Password)
}

func NewUserGrpcHandler(userUsecase userUsecase.UserUsecase) userPb.UserGrpcServiceServer {
	return &userGrpcHandler{userUsecase: userUsecase}
}
