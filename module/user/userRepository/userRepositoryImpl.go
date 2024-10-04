package userRepository

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// BuyBookRes implements UserRepository.
func (u *userRepositoryImpl) BuyBookRes(res *user.BuyBookRes, cfg *config.Config) {
	resInByte, err := json.Marshal(res)
	if err != nil {
		log.Printf("error: BuyBookRes: %s\n", err.Error())
		return
	}

	if err := queue.PushMessageWithKeyToQueue([]string{cfg.Kafka.Url}, "order", "usertransaction", resInByte); err != nil {
		log.Printf("error: BuyBookRes: %s\n", err.Error())
	}
}

// DeleteUserTransaction implements UserRepository.
func (u *userRepositoryImpl) DeleteUserTransaction(transactionID uint) error {
	if err := u.db.Delete(&user.UserTransactions{}, transactionID).Error; err != nil {
		return errors.New("error: delete user transaction")
	}
	return nil
}

// GetUserTransactions implements UserRepository.
func (u *userRepositoryImpl) GetUserTransactions(in *user.UserTransactions) ([]user.UserTransactions, error) {
	transactions := []user.UserTransactions{}
	if err := u.db.Where(in).Find(&transactions).Error; err != nil {
		log.Printf("error: GetUserTransactions: %s\n", err.Error())
		return nil, errors.New("error: get user transactions failed")
	}
	return transactions, nil
}

// CreateUserTransacton implements UserRepository.
func (u *userRepositoryImpl) CreateUserTransaction(in *user.UserTransactions) error {
	if err := u.db.Create(&in).Error; err != nil {
		log.Printf("error: CreateUserTransaction: %s\n", err.Error())
		return errors.New("error: create user transaction ")
	}
	return nil
}

// GetOneUser implements UserRepository.
func (u *userRepositoryImpl) GetOneUser(in *user.User) (*user.User, error) {
	user := &user.User{}
	if err := u.db.Where(&in).First(&user).Error; err != nil {
		log.Printf("error: GetOneUser: %s\n", err.Error())
		return nil, errors.New("error: user not found")
	}
	return user, nil
}

// createOneUser implements UserRepository.
func (u *userRepositoryImpl) CreateOneUser(in *user.User) (uint, error) {
	if err := u.db.Create(in).Error; err != nil {
		log.Printf("error: CreateOneUser: %s\n", err.Error())
		if strings.HasSuffix(err.Error(), "(SQLSTATE 23505)") {
			return 0, errors.New("error: username already in use")
		}
		return 0, errors.New("error: create user failed")
	}
	return in.ID, nil
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}
