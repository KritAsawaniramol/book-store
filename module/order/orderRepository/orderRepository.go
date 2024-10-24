package orderRepository

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
)

type OrderRepository interface {
	FindBookInIds(grpcUrl string, req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error)
	CreateOrder(in *order.Orders) error
	UpdateOrderByID(id uint, in *order.Orders) error
	DecreaseUserMoney(cfg *config.Config, req *user.BuyBookReq) error
	AddBookToShelf(cfg *config.Config, req *shelf.AddBooksReq) error
	RollbackUserTransaction(cfg *config.Config, req *user.RollbackUserTransactionReq) error
	RollbackAddBooks(cfg *config.Config, req *shelf.RollbackAddBooks) error
	FindOrdersWithBookDetail(in *order.Orders, status []string) ([]order.Orders, error)
	GetOrders(ids []uint, userIDs []uint, status []string, preloadOrdersBooks bool) ([]order.Orders, error) 
}