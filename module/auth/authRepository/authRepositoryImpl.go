package authRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"gorm.io/gorm"
)

type authRepositoryImpl struct {
	db *gorm.DB
}

// UpdateOneCredential implements AuthRepository.
func (a *authRepositoryImpl) UpdateOneCredentialByID(credentialID uint, in *auth.Credential) error {
	err := a.db.Model(&auth.Credential{}).Where("id = ?", credentialID).Updates(in).Error
	if err != nil {
		log.Printf("error: UpdateOneCredential: %s\n", err.Error())
		return errors.New("error: update credential failed")
	}
	return nil
}

// DeleteOneUserCredentialByID implements AuthRepository.
func (a *authRepositoryImpl) DeleteOneUserCredentialByID(credentialID uint) error {
	if err := a.db.Delete(&auth.Credential{}, credentialID).Error; err != nil {
		log.Printf("error: DeleteOneUserCredentialByID: %s\n", err.Error())
		return errors.New("error: delete user credential failed")
	}
	return nil
}

// GetOneUserCredential implements AuthRepository.
func (a *authRepositoryImpl) GetOneUserCredential(in *auth.Credential) (*auth.Credential, error) {
	credential := &auth.Credential{}
	if err := a.db.Where(&in).First(&credential).Error; err != nil {
		log.Printf("error: GetOneUserCredential: %s\n", err.Error())
		return nil, errors.New("error: credential not found")
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
		return 0, errors.New("error: create user credential failed")
	}
	return in.ID, nil
}

func (a *authRepositoryImpl) FindOneUserProfileToLogin(grpcUrl string, req *userPb.FindUserProfileToLoginReq) (*userPb.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}

	userProfile, err := conn.User().FindUserProfileToLogin(ctx, req)
	if err != nil {
		log.Printf("error: FindOneUserProfileToLogin: %s\n", err.Error())
		return nil, errors.New("error: email or password are incorrect")
	}
	return userProfile, nil
}

// FindOneUserProfileToRefresh implements AuthRepository.
func (a *authRepositoryImpl) FindOneUserProfileToRefresh(grpcUrl string, req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}

	userProfile, err := conn.User().FindOneUserProfileToRefresh(ctx, req)
	if err != nil {
		log.Printf("error: FindOneUserProfileToRefresh: %s\n", err.Error())
		return nil, errors.New("error: user not found")
	}
	return userProfile, nil
}

func NewAuthRepositoryImpl(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}
