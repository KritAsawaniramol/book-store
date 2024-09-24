package shelfUsecase

import (
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfRepository"
)

type shelfUsecaseImpl struct {
	shelfRepository shelfRepository.ShelfRepository
}

func NewShelfUsecaseImpl(shelfRepository shelfRepository.ShelfRepository) ShelfUsecase {
	return &shelfUsecaseImpl{
		shelfRepository: shelfRepository,
	}
}
