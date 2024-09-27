package userUsecase

import (
	"errors"
	"log"

	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userUsecaseImpl struct {
	userRepository userRepository.UserRepository
}

// FindOneUserByID implements UserUsecase.
func (u *userUsecaseImpl) FindOneUserByID(userID uint) (*userPb.UserProfile, error) {
	condition := &user.User{}
	condition.ID = userID
	user, err := u.userRepository.GetOneUser(condition)
	if err != nil {
		return nil, err
	}
	return &userPb.UserProfile{
		Id:        uint64(user.ID),
		Username:  user.Username,
		RoleId:    uint32(user.RoleID),
		Coin:      user.Coin,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// FindOneUserByUsernameAndPassword implements UserUsecase.
func (u *userUsecaseImpl) FindOneUserByUsernameAndPassword(username string, password string) (*userPb.UserProfile, error) {
	user, err := u.userRepository.GetOneUser(&user.User{Username: username})
	if err != nil {
		return nil, err
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("error: FindOneUserByUsernameAndPassword: %s\n", err.Error())
		return nil, errors.New("error: password is incorrect")
	}
	return &userPb.UserProfile{
		Id:        uint64(user.ID),
		Username:  user.Username,
		RoleId:    uint32(user.RoleID),
		Coin:      user.Coin,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// Register implements UserUsecase.
func (u *userUsecaseImpl) Register(registReq *user.UserRegisterReq) (uint, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error: Register: %s\n", err.Error())
		return 0, errors.New("error: failed to hash password")
	}

	userID, err := u.userRepository.CreateOneUser(&user.User{
		Username: registReq.Username,
		Password: string(hashedPassword),
		RoleID:   2,
		Coin:     0,
	})
	if err != nil {
		return 0, err
	}
	return userID, err
}

func NewUserUsecaseImpl(userRepository userRepository.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepository: userRepository,
	}
}
