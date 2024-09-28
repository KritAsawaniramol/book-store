package bookUsecase

import "github.com/kritAsawaniramol/book-store/module/book"

type BookUsecase interface {
	CreateOneBook(req *book.CreateBookReq) (uint, error)
}
