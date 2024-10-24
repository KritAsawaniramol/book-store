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

// GetOrders implements OrderRepository.
func (o *orderRepositoryImpl) GetOrders(ids []uint, userIDs []uint, status []string, preloadOrdersBooks bool) ([]order.Orders, error) {
	result := []order.Orders{}
	err := o.db.Transaction(func(tx *gorm.DB) error {
		// tx := o.db.Begin()
		if len(ids) > 0 && ids != nil {
			tx = tx.Where("id IN ?", ids)
		}
	
		if len(userIDs) > 0 && userIDs != nil {
			tx = tx.Where("user_id IN ?", userIDs)
		}
		if len(status) > 0 && status != nil {
			tx = tx.Where("status IN ?", status)
		}
		if preloadOrdersBooks  {
			tx = tx.Preload("OrdersBooks")
		}
	
		if err := tx.Find(&result).Error; err != nil {
			log.Printf("error: GetOrdersWithOrderBooks: %s\n", err.Error())
			return  errors.New("error: get orders failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindOneOrderWithBookDetail implements OrderRepository.
func (o *orderRepositoryImpl) FindOrdersWithBookDetail(in *order.Orders, status []string) ([]order.Orders, error) {
	result := []order.Orders{}
	tx := o.db.Where(&in)
	if len(status) > 0 || status == nil {
		tx = tx.Where("status IN ?", status)
	}

	if err := tx.Preload("OrdersBooks").Find(&result).Error; err != nil {
		log.Printf("error: FindOneOrderWithBookDetail: %s\n", err.Error())
		return nil, errors.New("error: get order failed")
	}
	return result, nil
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
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
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
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
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
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
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
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
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
