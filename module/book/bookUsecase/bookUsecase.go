package bookUsecase

import "github.com/kritAsawaniramol/book-store/module/book"

type BookUsecase interface {
	CreateOneBook(req *book.CreateBookReq) (uint, error)
	SearchBooks(req *book.SearchBooksReq) (*book.SearchBooksRes, error)
	GetOneBook(bookID uint) (*book.BookRes, error)
	GetTags() ([]book.BookTags, error)
}
