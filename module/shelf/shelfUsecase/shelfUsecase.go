package shelfUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf"
)

type ShelfUsecase interface {
	AddBooks(cfg *config.Config, req *shelf.AddBooksReq) 
	RollBacksAddBook(cfg *config.Config, req *shelf.RollbackAddBooks)
GetMyShelves(cfg *config.Config, userID uint, bookIDQuery uint) (*shelf.GetMyShelf, error)
}