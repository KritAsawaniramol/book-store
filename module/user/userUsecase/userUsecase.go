package userUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
)

type UserUsecase interface {
	Register(registReq *user.UserRegisterReq) (uint, error)
	FindOneUserByUsernameAndPassword(username string, password string) (*userPb.UserProfile, error)
	FindOneUserByID(userID uint) (*userPb.UserProfile, error)
	CreateUserTransaction(req *user.CreateUserTransactionReq) (*user.CreateUserTransactionRes, error)
	BuyBook(cfg *config.Config, req *user.BuyBookReq)
	GetUserBalance(userID uint) (*user.UserBalanceRes, error)
	RollbackUserTransaction(req *user.RollbackUserTransactionReq)
}
