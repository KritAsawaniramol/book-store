package shelfRepository

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf"
)

type ShelfRepository interface {
	InsertUserBooks(in []shelf.Shelves) error
	AddBookRes(cfg *config.Config, res *shelf.AddBooksRes)
	DeleteUserBookInIDs(cfg *config.Config, ids []uint) error
}
