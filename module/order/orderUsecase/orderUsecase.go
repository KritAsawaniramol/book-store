package orderUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
)

type OrderUsecase interface {
	BuyBooks(cfg *config.Config, req *order.BuyBooksReq) (*order.BuyBooksRes, error)
	HandleBuyBooksRes(cfg *config.Config, res *user.BuyBookRes)
	HandleAddBookRes(cfg *config.Config, res *shelf.AddBooksRes)
}
