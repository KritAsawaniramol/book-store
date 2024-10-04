package orderUsecase

import (
	"fmt"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/order/orderRepository"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
)

const (
	pendding  = "pendding"
	failed    = "failed"
	completed = "completed"
)

type orderUsecaseImpl struct {
	orderRepository orderRepository.OrderRepository
}

// HandleAddBookRes implements OrderUsecase.
func (o *orderUsecaseImpl) HandleAddBookRes(cfg *config.Config, res *shelf.AddBooksRes) {
	if res.Error != "" {
		o.orderRepository.UpdateOrderByID(
			res.OrderID, &order.Orders{Status: failed, Note: res.Error},
		)
		o.orderRepository.RollbackUserTransaction(cfg, &user.RollbackUserTransactionReq{
			TransactionID: res.TransactionID,
		})
		return
	}

	if err := o.orderRepository.UpdateOrderByID(
		res.OrderID, &order.Orders{Status: completed},
	); err != nil {
		o.orderRepository.RollbackUserTransaction(cfg, &user.RollbackUserTransactionReq{
			TransactionID: res.TransactionID,
		})
		o.orderRepository.RollbackAddBooks(cfg, &shelf.RollbackAddBooks{
			ShelfIDs: res.ShelfIDs,
		})
	}
}

// HandleBuyBooksRes implements OrderUsecase.
func (o *orderUsecaseImpl) HandleBuyBooksRes(cfg *config.Config, res *user.BuyBookRes) {
	if res.Error != "" {
		o.orderRepository.UpdateOrderByID(
			res.OrderID, &order.Orders{Status: failed, Note: res.Error},
		)
		return
	}

	// add book to user shelf
	if err := o.orderRepository.AddBookToShelf(cfg, &shelf.AddBooksReq{
		OrderID:       res.OrderID,
		UserID:        res.UserID,
		TransactionID: res.TransactionID,
		BookIDs:       res.BookIDs,
	}); err != nil {
		//rollback transaction
		o.orderRepository.RollbackUserTransaction(cfg, &user.RollbackUserTransactionReq{
			TransactionID: res.TransactionID,
		})
	}
}

// BuyBooks implements OrderUsecase.
func (o *orderUsecaseImpl) BuyBooks(cfg *config.Config, req *order.BuyBooksReq) (*order.BuyBooksRes, error) {
	bookIDsUint := []uint{}
	bookIDsUint64 := []uint64{}
	for _, b := range req.Books {
		bookIDsUint = append(bookIDsUint, b.BookID)
		bookIDsUint64 = append(bookIDsUint64, uint64(b.BookID))
	}
	//Get book info
	fmt.Printf("bookIDsUint64: %v\n", bookIDsUint64)
	booksInfo, err := o.orderRepository.FindBookInIds(cfg.Grpc.BookUrl, &bookPb.FindBooksInIdsReq{Ids: bookIDsUint64})
	if err != nil {
		return nil, err
	}

	fmt.Printf("booksInfo: %v\n", booksInfo)

	orderBook := []order.OrderBook{}
	var totalPrice uint = 0
	for _, b := range booksInfo.Book {
		totalPrice += uint(b.Price)
		orderBook = append(orderBook, order.OrderBook{
			BookID: uint(b.Id),
			Price:  uint(b.Price),
		})
	}

	//create createOrderReq in database, status is "pendding"
	createOrderReq := &order.Orders{
		UserID: req.UserID,
		Status: pendding,
		Books:  orderBook,
		Total:  totalPrice,
	}
	if err := o.orderRepository.CreateOrder(createOrderReq); err != nil {
		return nil, err
	}

	//decress user money
	decressUserMoneyReq := &user.BuyBookReq{
		OrderID: createOrderReq.ID,
		UserID:  req.UserID,
		BookIDs: bookIDsUint,
		Total:   totalPrice,
	}
	if err := o.orderRepository.DecreaseUserMoney(cfg, decressUserMoneyReq); err != nil {
		o.orderRepository.UpdateOrderByID(
			createOrderReq.ID,
			&order.Orders{Status: failed, Note: err.Error()},
		)
		return nil, err
	}
	return nil, nil
}

func NewOrderUsecaseImpl(orderRepository orderRepository.OrderRepository) OrderUsecase {
	return &orderUsecaseImpl{orderRepository: orderRepository}
}
