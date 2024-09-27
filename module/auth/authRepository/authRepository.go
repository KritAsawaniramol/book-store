package authRepository

import (
	"github.com/kritAsawaniramol/book-store/module/auth"
	userPb "github.com/kritAsawaniramol/book-store/module/user/userPb"
)

type AuthRepository interface {
	FindOneUserProfile(grpcUrl string, req *userPb.FindUserProfileReq) (*userPb.UserProfile, error)
	CreateOneUserCredential(in *auth.Credential) (uint, error)
	GetOneUserCredential(in *auth.Credential) (*auth.Credential, error)
}
