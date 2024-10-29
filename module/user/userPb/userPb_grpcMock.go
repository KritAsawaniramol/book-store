package userPb

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type UserGrpcServiceClientMock struct {
	mock.Mock
}

// FindOneUserProfileToRefresh implements UserGrpcServiceClient.
func (u *UserGrpcServiceClientMock) FindOneUserProfileToRefresh(ctx context.Context, in *FindOneUserProfileToRefreshReq, opts ...grpc.CallOption) (*UserProfile, error) {
	arguments := []interface{}{ctx, in}
	for _, v := range opts {
		arguments = append(arguments, v)
	}
	args := u.Called(arguments...)

	userProfile, ok := args.Get(0).(*UserProfile)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return userProfile, args.Error(1)
}

// FindUserProfileToLogin implements UserGrpcServiceClient.
func (u *UserGrpcServiceClientMock) FindUserProfileToLogin(ctx context.Context, in *FindUserProfileToLoginReq, opts ...grpc.CallOption) (*UserProfile, error) {
	arguments := []interface{}{ctx, in}
	for _, v := range opts {
		arguments = append(arguments, v)
	}
	args := u.Called(arguments...)
	userProfile, ok := args.Get(0).(*UserProfile)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return userProfile, args.Error(1)
}

func NewUserGrpcServiceClientMock() *UserGrpcServiceClientMock {
	return &UserGrpcServiceClientMock{}
}
