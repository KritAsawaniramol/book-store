package authRepository

import (
	"context"
	"log"
	"time"

	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"gorm.io/gorm"
)

type authRepositoryImpl struct {
	db                    *gorm.DB
	userServiceGrpcClient userPb.UserGrpcServiceClient
}

// UpdateOneCredential implements AuthRepository.
func (a *authRepositoryImpl) UpdateOneCredentialByID(credentialID uint, in *auth.Credential) error {
	err := a.db.Model(&auth.Credential{}).Where("id = ?", credentialID).Updates(in).Error
	if err != nil {
		log.Printf("error: UpdateOneCredential: %s\n", err.Error())
		return auth.ErrUpdateCredential
	}
	return nil
}

// DeleteOneUserCredentialByID implements AuthRepository.
func (a *authRepositoryImpl) DeleteOneUserCredentialByID(credentialID uint) error {
	if err := a.db.Delete(&auth.Credential{}, credentialID).Error; err != nil {
		log.Printf("error: DeleteOneUserCredentialByID: %s\n", err.Error())
		return auth.ErrDeleteUserCredentialFail
	}
	return nil
}

// GetOneUserCredential implements AuthRepository.
func (a *authRepositoryImpl) GetOneUserCredential(in *auth.Credential) (*auth.Credential, error) {
	credential := &auth.Credential{}
	if err := a.db.Where(&in).First(&credential).Error; err != nil {
		log.Printf("error: GetOneUserCredential: %s\n", err.Error())
		return nil, auth.ErrCredentialNotFound
	}
	return credential, nil
}

// CreateOneUserCredential implements AuthRepository.
func (a *authRepositoryImpl) CreateOneUserCredential(in *auth.Credential) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := a.db.WithContext(ctx).Create(in).Error
	if err != nil {
		log.Printf("error: CreateOneUserCredential: %s\n", err.Error())
		return 0, auth.ErrCreateUserCredential
	}
	return in.ID, nil
}

func (a *authRepositoryImpl) FindOneUserProfileToLogin(req *userPb.FindUserProfileToLoginReq) (*userPb.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userProfile, err := a.userServiceGrpcClient.FindUserProfileToLogin(ctx, req)
	if err != nil {
		log.Printf("error: FindOneUserProfileToLogin: %s\n", err.Error())
		return nil, auth.ErrEmailOrPasswordIncorrect
	}
	return userProfile, nil
}

// FindOneUserProfileToRefresh implements AuthRepository.
func (a *authRepositoryImpl) FindOneUserProfileToRefresh(req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userProfile, err := a.userServiceGrpcClient.FindOneUserProfileToRefresh(ctx, req)
	if err != nil {
		log.Printf("error: FindOneUserProfileToRefresh: %s\n", err.Error())
		return nil, auth.ErrUserNotFound
	}
	return userProfile, nil
}

func NewAuthRepositoryImpl(db *gorm.DB, userServiceGrpcClient userPb.UserGrpcServiceClient) AuthRepository {
	return &authRepositoryImpl{db: db, userServiceGrpcClient: userServiceGrpcClient}
}
