package bookRepository

import "github.com/kritAsawaniramol/book-store/module/book"

type BookRepository interface {
	CreateOneBook(in *book.Books) error
	CreateTags(in []book.Tags) error
	DeleteTags(in []book.Tags) error
	SearchBook(
		limit int,
		order string,
		offest int,
		title string,
		maxPrice *uint,
		minPrice *uint,
		authorName string,
		isAvailableInStore *bool,
		tagIDs []*uint,
	) ([]book.Books, int64, error)
	GetBooksInIDs(ids []uint) ([]book.Books, int64, error)
	GetOneBook(in *book.Books) (*book.Books, error)
	GetTags(in *book.Tags) ([]book.Tags, error)
	UpdateOneBookDetail(bookID uint,in *book.Books) error
	UpdateNonZeroBookFields(bookID uint, in *book.Books) error 
}
