package shelfRepository

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/shelf"
)

type ShelfRepository interface {
	InsertUserBooks(in []shelf.Shelves) error
	AddBookRes(cfg *config.Config, res *shelf.AddBooksRes)
	DeleteUserBookInIDs(cfg *config.Config, ids []uint) error
	GetOneItemInShelf(in *shelf.Shelves) (*shelf.Shelves, error) 
	GetShelves(ids []uint, userIDs []uint, bookIDs []uint) ([]shelf.Shelves, error) 
	FindBookInIds(grpcUrl string, req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error)
}
