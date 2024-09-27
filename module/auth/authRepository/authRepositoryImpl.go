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

// findUserProfile implements AuthRepository.
func (a *authRepositoryImpl) FindOneUserProfile( grpcUrl string, req *userPb.FindUserProfileReq) (*userPb.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}

	userProfile, err := conn.User().FindUserProfile(ctx, req)
	if err != nil {
		log.Printf("error: FindOneUserProfile: %s\n", err.Error())
		return nil, errors.New("error: email or password are incorrect")
	}
	return userProfile, nil
}

func NewAuthRepositoryImpl(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}
