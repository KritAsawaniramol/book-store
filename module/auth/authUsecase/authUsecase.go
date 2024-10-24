package authUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
)

type AuthUsecase interface {
	Login(cfg *config.Config, req *auth.LoginReq) (*auth.LoginRes, error)
	Logout(req *auth.LogoutReq) error
	RefreshToken(cfg *config.Config, req *auth.RefreshTokenReq) (*auth.CredentialRes, error)
	AccessTokenSearch(accessToken string) (*authPb.AccessTokenSearchRes, error)
}
