package bookRepository

import "github.com/kritAsawaniramol/book-store/module/book"

type BookRepository interface {
	CreateOneBook(in *book.Books) error
	CreateTags(in []book.Tags) error
	DeleteTags(in []book.Tags) error
	GetBooks(
		limit int,
		order interface{},
		offest uint,
		tagIDs []uint,
		conditions ...interface{},
	) ([]book.Books, int64, error)
	GetOneBook(in *book.Books) (*book.Books, error)
}
