package shelfUsecase

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfRepository"
	"github.com/kritAsawaniramol/book-store/util"
)

type shelfUsecaseImpl struct {
	shelfRepository shelfRepository.ShelfRepository
}

// GetMyShelves implements ShelfUsecase.
func (s *shelfUsecaseImpl) GetMyShelves(cfg *config.Config, userID uint, bookIDQuery uint) (*shelf.GetMyShelf, error) {

	bookIDsDQuery := []uint{}
	if bookIDQuery != 0 {
		bookIDsDQuery = append(bookIDsDQuery, bookIDQuery)
	} else {
		bookIDsDQuery = nil
	}
	shelves, err := s.shelfRepository.GetShelves(nil, []uint{userID}, bookIDsDQuery)
	if err != nil {
		return nil, err
	}

	bookIDs := []uint64{}
	uniqueBooks := map[uint]*bookPb.Book{}

	for _, s := range shelves {
		if _, ok := uniqueBooks[s.BooksID]; ok {
			continue
		} else {
			uniqueBooks[s.BooksID] = &bookPb.Book{}
			bookIDs = append(bookIDs, uint64(s.BooksID))
		}

	}
	findBooksInIdsRes, err := s.shelfRepository.FindBookInIds(cfg.Grpc.BookUrl, &bookPb.FindBooksInIdsReq{
		Ids: bookIDs,
	})
	if err != nil {
		return nil, err
	}
	for _, b := range findBooksInIdsRes.Book {
		uniqueBooks[uint(b.Id)] = b
	}

	shelvesRes := []shelf.ShelfRes{}

	for _, sh := range shelves {
		if v, ok := uniqueBooks[sh.BooksID]; ok && v != nil {

			shelvesRes = append(shelvesRes, shelf.ShelfRes{
				ID:     sh.ID,
				UserID: sh.UsersID,
				Book: shelf.BookRes{
					ID:            uint(v.Id),
					Title:         v.Title,
					AuthorName:    v.AuthorName,
					CoverImageUrl: v.CoverImagePath,
					CreatedAt:     sh.CreatedAt,
				},
			})
		}
	}
	util.PrintObjInJson(shelvesRes)
	return &shelf.GetMyShelf{Shelves: shelvesRes}, nil
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
