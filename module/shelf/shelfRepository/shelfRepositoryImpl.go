package shelfRepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
	"gorm.io/gorm"
)

type shelfRepositoryImpl struct {
	db *gorm.DB
}

// DeleteBookInIDs implements ShelfRepository.
func (s *shelfRepositoryImpl) DeleteUserBookInIDs(cfg *config.Config, ids []uint) error {
	if err := s.db.Delete(&shelf.Shelves{}, ids).Error; err != nil {
		log.Printf("error: DeleteBookInIDs: %s\n", err.Error())
		return errors.New("error: delete user's books failed")
	}
	return nil
}

// AddBookRes implements ShelfRepository.
func (s *shelfRepositoryImpl) AddBookRes(cfg *config.Config, res *shelf.AddBooksRes) {
	resInByte, err := json.Marshal(res)
	if err != nil {
		log.Printf("error: AddBookRes: %s\n", err.Error())
		return
	}

	if err := queue.PushMessageWithKeyToQueue([]string{cfg.Kafka.Url}, "order", "addbook", resInByte); err != nil {
		log.Printf("error: AddBookRes: %s\n", err.Error())
	}
}

// InsertUserBooks implements ShelfRepository.
func (s *shelfRepositoryImpl) InsertUserBooks(in []shelf.Shelves) error {
	fmt.Printf("in: %v\n", in)
	if err := s.db.Create(&in).Error; err != nil {
		log.Printf("error: InsertUserBooks: %s\n", err.Error())
		return errors.New("error: insert user books failed")
	}
	return nil
}

func NewUserRepositoryImpl(db *gorm.DB) ShelfRepository {
	return &shelfRepositoryImpl{
		db: db,
	}
}
