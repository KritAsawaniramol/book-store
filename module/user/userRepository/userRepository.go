package userRepository

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
)

type UserRepository interface {
	CreateOneUser(in *user.User) (uint, error)
	GetOneUser(in *user.User) (*user.User, error)
	CreateUserTransaction(in *user.UserTransactions) error
	GetUserTransactions(in *user.UserTransactions) ([]user.UserTransactions, error)
	DeleteUserTransaction(transactionID uint) error
	BuyBookRes(res *user.BuyBookRes, cfg *config.Config) 
}
