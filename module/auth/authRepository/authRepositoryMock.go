package authRepository

import (
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func NewAuthRepositoryMock() *AuthRepositoryMock {
	return &AuthRepositoryMock{}
}

// CreateOneUserCredential implements AuthRepository.
func (a *AuthRepositoryMock) CreateOneUserCredential(in *auth.Credential) (uint, error) {
	args := a.Called(in)
	return args.Get(0).(uint), args.Error(1)
}

// DeleteOneUserCredentialByID implements AuthRepository.
func (a *AuthRepositoryMock) DeleteOneUserCredentialByID(credentialID uint) error {
	args := a.Called(credentialID)
	return args.Error(0)
}

// FindOneUserProfileToLogin implements AuthRepository.
func (a *AuthRepositoryMock) FindOneUserProfileToLogin(req *userPb.FindUserProfileToLoginReq) (*userPb.UserProfile, error) {
	args := a.Called(req)
	userProfile, ok := args.Get(0).(*userPb.UserProfile)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return userProfile, args.Error(1)
}

// FindOneUserProfileToRefresh implements AuthRepository.
func (a *AuthRepositoryMock) FindOneUserProfileToRefresh(req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
	args := a.Called(req)
	userProfile, ok := args.Get(0).(*userPb.UserProfile)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return userProfile, args.Error(1)
}

// GetOneUserCredential implements AuthRepository.
func (a *AuthRepositoryMock) GetOneUserCredential(in *auth.Credential) (*auth.Credential, error) {
	args := a.Called(in)
	credential, ok := args.Get(0).(*auth.Credential)
	if !ok && args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return credential, args.Error(1)
}

// UpdateOneCredentialByID implements AuthRepository.
func (a *AuthRepositoryMock) UpdateOneCredentialByID(credentialID uint, in *auth.Credential) error {
	args := a.Called(credentialID, in)
	return args.Error(0)
}
