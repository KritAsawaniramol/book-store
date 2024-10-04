package orderRepository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

// RollbackAddBooks implements OrderRepository.
func (o *orderRepositoryImpl) RollbackAddBooks(cfg *config.Config, req *shelf.RollbackAddBooks) error {
	reqInByte, err := json.Marshal(req)
	if err != nil {
		log.Printf("error: RollbackAddBooks: %s\n", err.Error())
		return errors.New("error: rollback add books failed")
	}
	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		"shelf",
		"rollbackaddbook",
		reqInByte,
	); err != nil {
		log.Printf("error: RollbackAddBooks: %s\n", err.Error())
		return errors.New("error: rollback add books failed")
	}
	return nil
}

// RollbackUserTransaction implements OrderRepository.
func (o *orderRepositoryImpl) RollbackUserTransaction(cfg *config.Config, req *user.RollbackUserTransactionReq) error {
	reqInByte, err := json.Marshal(req)
	if err != nil {
		log.Printf("error: RollbackUserTransaction: %s\n", err.Error())
		return errors.New("error: rollback usertransaction failed")
	}
	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		"user",
		"rollbackusertransaction",
		reqInByte,
	); err != nil {
		log.Printf("error: RollbackUserTransaction: %s\n", err.Error())
		return errors.New("error: rollback usertransaction failed")
	}
	return nil
}

// AddBookToShelf implements OrderRepository.
func (o *orderRepositoryImpl) AddBookToShelf(cfg *config.Config, req *shelf.AddBooksReq) error {
	reqInByte, err := json.Marshal(req)
	if err != nil {
		log.Printf("error: AddBookToShelf: %s\n", err.Error())
		return errors.New("error: add books to shelf failed")
	}
	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		"shelf",
		"addbook",
		reqInByte,
	); err != nil {
		log.Printf("error: AddBookToShelf: %s\n", err.Error())
		return errors.New("error: add books to shelf failed")
	}
	return nil
}

// UpdateOrderByID implements OrderRepository.
func (o *orderRepositoryImpl) UpdateOrderByID(id uint, in *order.Orders) error {
	if err := o.db.Model(&order.Orders{}).Where("id", id).Updates(in).Error; err != nil {
		log.Printf("error: UpdateOrderByID: %s\n", err.Error())
		return errors.New("error: update order failed")
	}
	return nil
}

// DecreaseUserMoney implements OrderRepository.
func (o *orderRepositoryImpl) DecreaseUserMoney(cfg *config.Config, req *user.BuyBookReq) error {
	reqInByte, err := json.Marshal(req)
	if err != nil {
		log.Printf("error: DecreaseUserMoney: %s\n", err.Error())
		return errors.New("error: decrease user money failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		"user",
		"buy",
		reqInByte,
	); err != nil {
		log.Printf("error: DecreaseUserMoney: %s\n", err.Error())
		return errors.New("error: decrease user money failed")
	}
	return nil
}

// CreateOrder implements OrderRepository.
func (o *orderRepositoryImpl) CreateOrder(in *order.Orders) error {
	err := o.db.Create(&in).Error
	if err != nil {
		log.Printf("error: CreateOrder: %s\n", err.Error())
		return errors.New("error: create order failed")
	}
	return nil
}

// FindBookInIds implements OrderRepository.
func (o *orderRepositoryImpl) FindBookInIds(grpcUrl string, req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error) {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	books, err := conn.Book().FindBooksInIds(ctx, req)
	if err != nil {
		log.Printf("error: FindBookInIds: %s\n", err.Error())
		return nil, errors.New("error: books not found")
	}
	return books, nil
}

func NewOrderRepositoryImpl(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}
