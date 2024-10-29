package authRepository

import (
	"github.com/kritAsawaniramol/book-store/module/auth"
	userPb "github.com/kritAsawaniramol/book-store/module/user/userPb"
)

type AuthRepository interface {
	FindOneUserProfileToLogin(req *userPb.FindUserProfileToLoginReq) (*userPb.UserProfile, error)
	FindOneUserProfileToRefresh(req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error)

	CreateOneUserCredential(in *auth.Credential) (uint, error)
	GetOneUserCredential(in *auth.Credential) (*auth.Credential, error)
	DeleteOneUserCredentialByID(credentialID uint) error
	UpdateOneCredentialByID(credentialID uint, in *auth.Credential) error
}
