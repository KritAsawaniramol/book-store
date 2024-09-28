package bookRepository

import "github.com/kritAsawaniramol/book-store/module/book"

type BookRepository interface {
	CreateOneBook(in *book.Books) error
	CreateTags(in []book.Tags)  error
	DeleteTags(in []book.Tags) error
}
