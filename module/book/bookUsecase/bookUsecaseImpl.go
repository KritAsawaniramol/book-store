package bookUsecase

import "github.com/kritAsawaniramol/book-store/module/book/bookRepository"

type bookUsecaseImpl struct {
	bookRepository bookRepository.BookRepository
}

func NewBookUsecaseImpl(bookRepository bookRepository.BookRepository) BookUsecase {
	return &bookUsecaseImpl{bookRepository: bookRepository}
}
