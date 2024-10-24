package userRepository

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/stripe/stripe-go/v80"
)

type UserRepository interface {
	CreateOneUser(in *user.User) (uint, error)
	GetOneUser(in *user.User) (*user.User, error)
	GetUserTransactions(in *user.UserTransactions) ([]user.UserTransactions, error)
	DeleteUserTransaction(transactionID uint) error
	BuyBookRes(res *user.BuyBookRes, cfg *config.Config)
	CreateUserTransaction(in *user.UserTransactions) error
	GetUserInIDs(ids []uint) ([]user.User, int64, error)
	CreateOneTopUpOrder(in *user.TopUpOrder) error
	GetOneTopUpOrder(in *user.TopUpOrder) (*user.TopUpOrder, error)
	CheckOutTopUp(cfg *config.Config,amount int64) (*stripe.CheckoutSession, error)
	UpdateOneTopUpOrderStatusBySessionID(sessionID string, newStatus string) error
}
