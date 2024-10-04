package shelfUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfRepository"
)

type shelfUsecaseImpl struct {
	shelfRepository shelfRepository.ShelfRepository
}

// RollBacksAddBook implements ShelfUsecase.
func (s *shelfUsecaseImpl) RollBacksAddBook(cfg *config.Config, req *shelf.RollbackAddBooks) {
	s.shelfRepository.DeleteUserBookInIDs(cfg, req.ShelfIDs)
}

// AddBooks implements ShelfUsecase.
func (s *shelfUsecaseImpl) AddBooks(cfg *config.Config, req *shelf.AddBooksReq) {
	userBooks := []shelf.Shelves{}
	for _, id := range req.BookIDs {
		userBooks = append(userBooks, shelf.Shelves{
			UsersID: req.UserID,
			BooksID: id,
		})
	}
	res := &shelf.AddBooksRes{
		OrderID:       req.OrderID,
		UserID:        req.UserID,
		TransactionID: req.TransactionID,
		Error:         "",
	}

	if err := s.shelfRepository.InsertUserBooks(userBooks); err != nil {
		res.Error = err.Error()
		s.shelfRepository.AddBookRes(cfg, res)
		return
	}

	shelfIDs := []uint{}
	for _, userBook := range userBooks {
		shelfIDs = append(shelfIDs, userBook.ID)
	}
	res.ShelfIDs = shelfIDs
	s.shelfRepository.AddBookRes(cfg, res)
}

func NewShelfUsecaseImpl(shelfRepository shelfRepository.ShelfRepository) ShelfUsecase {
	return &shelfUsecaseImpl{
		shelfRepository: shelfRepository,
	}
}
