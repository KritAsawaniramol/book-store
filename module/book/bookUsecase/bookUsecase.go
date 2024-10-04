package bookUsecase

import (
	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
)

type BookUsecase interface {
	CreateOneBook(req *book.CreateBookReq) (uint, error)
	SearchBooks(req *book.SearchBooksReq) (*book.SearchBooksRes, error)
	FindBookInIDs(req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error)
	GetOneBook(bookID uint) (*book.BookRes, error)
	GetTags() ([]book.BookTags, error)
}
