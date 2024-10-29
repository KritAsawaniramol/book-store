package authUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
	"github.com/stretchr/testify/mock"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func NewAuthUsecaseMock() *AuthUsecaseMock {
	return &AuthUsecaseMock{}
}

// AccessTokenSearch implements AuthUsecase.
func (a *AuthUsecaseMock) AccessTokenSearch(accessToken string) (*authPb.AccessTokenSearchRes, error) {
	args := a.Called(accessToken)
	accessTokenSearchRes, ok := args.Get(0).(*authPb.AccessTokenSearchRes)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return accessTokenSearchRes, args.Error(1)
}

// Login implements AuthUsecase.
func (a *AuthUsecaseMock) Login(cfg *config.Config, req *auth.LoginReq) (*auth.LoginRes, error) {
	args := a.Called(cfg, req)
	loginRes, ok := args.Get(0).(*auth.LoginRes)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return loginRes, args.Error(1)
}

// Logout implements AuthUsecase.
func (a *AuthUsecaseMock) Logout(req *auth.LogoutReq) error {
	args := a.Called(req)
	return args.Error(0)
}

// RefreshToken implements AuthUsecase.
func (a *AuthUsecaseMock) RefreshToken(cfg *config.Config, req *auth.RefreshTokenReq) (*auth.CredentialRes, error) {
	args := a.Called(cfg, req)
	credentialRes, ok := args.Get(0).(*auth.CredentialRes)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return credentialRes, args.Error(1)
}
