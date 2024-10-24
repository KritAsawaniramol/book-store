package userRepository

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// UpdateOneTopUpOrderStatusBySessionID implements UserRepository.
func (u *userRepositoryImpl) UpdateOneTopUpOrderStatusBySessionID(sessionID string, newStatus string) error {
	condition := &user.TopUpOrder{SessionID: sessionID}
	err := u.db.Model(&user.TopUpOrder{}).Where(condition).Update("status", newStatus).Error
	if err != nil {
		log.Printf("error: UpdateOneTopUpOrderStatusBySessionID: %s\n", err.Error())
		return errors.New("error: update top-up order status failed")
	}
	return nil
}

// CheckOut implements UserRepository.
func (u *userRepositoryImpl) CheckOutTopUp(cfg *config.Config, amount int64) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: []*string{stripe.String("card")},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("thb"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Top-Up"),
					},
					UnitAmount: stripe.Int64(amount * 100),
				},
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(cfg.Client.URL + "/top-up/success"),
		CancelURL:  stripe.String(cfg.Client.URL + "/top-up/cancel"),
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("error: TopUp: session.New: %v", err)
		return nil, errors.New("error: top-up failed")
	}
	return s, nil
}

// CreateOneTopUpOrder implements UserRepository.
func (u *userRepositoryImpl) CreateOneTopUpOrder(in *user.TopUpOrder) error {
	if err := u.db.Create(in).Error; err != nil {
		log.Printf("error: CreateOneTopUpOrder: %s\n", err.Error())
		return errors.New("error: create top-up order failed")
	}
	return nil
}

func (u *userRepositoryImpl) GetUserInIDs(ids []uint) ([]user.User, int64, error) {
	users := []user.User{}
	result := u.db.Where(ids).Find(&users)
	if result.Error != nil {
		log.Printf("error: GetUserInIDs: %s\n", result.Error.Error())
		return nil, 0, errors.New("error: get users failed")
	}
	var count int64
	result.Count(&count)
	return users, count, nil
}

// BuyBookRes implements UserRepository.
func (u *userRepositoryImpl) BuyBookRes(res *user.BuyBookRes, cfg *config.Config) {
	resInByte, err := json.Marshal(res)
	if err != nil {
		log.Printf("error: BuyBookRes: %s\n", err.Error())
		return
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"order",
		"usertransaction",
		resInByte,
	); err != nil {
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
	if err := u.db.Where(in).Order("created_at desc").Find(&transactions).Error; err != nil {
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

// GetOneTopUpOrder implements UserRepository.
func (u *userRepositoryImpl) GetOneTopUpOrder(in *user.TopUpOrder) (*user.TopUpOrder, error) {
	t := &user.TopUpOrder{}
	if err := u.db.Where(&in).First(&t).Error; err != nil {
		log.Printf("error: GetOneTopUpOrder: %s\n", err.Error())
		return t, errors.New("error: top-up order not found")
	}
	return t, nil
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
