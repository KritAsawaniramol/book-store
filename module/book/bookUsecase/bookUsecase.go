package bookUsecase

import (
	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
)

type BookUsecase interface {
	CreateOneBook(req *book.CreateBookReq) (uint, error)
	SearchBooks(req *book.SearchBooksReq, roleID uint) (*book.SearchBooksRes, error)
	FindBookInIDs(req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error)
	GetOneBook(bookID uint, roleID uint) (*book.BookRes, error)
	GetOneBookFilePath(bookID uint) (string, error)
	UpdateOneBookDetail(req *book.UpdateBookDetailReq) error
	GetTags() ([]book.BookTags, error)
	UpdateOneBookCover(bookID uint, newImagePath string) error
	UpdateOneBookFile(bookID uint, newFilePath string) error
}
