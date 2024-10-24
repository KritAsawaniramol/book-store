package shelfRepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
	"gorm.io/gorm"
)

type shelfRepositoryImpl struct {
	db *gorm.DB
}

// FindBookInIds implements ShelfRepository.
func (s *shelfRepositoryImpl)  FindBookInIds(grpcUrl string, req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error) {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	books, err := conn.Book().FindBooksInIds(ctx, req)
	if err != nil {
		log.Printf("error: FindBookInIds: %s\n", err.Error())
		return nil, errors.New("error: books not found")
	}
	return books, nil
}

func (s *shelfRepositoryImpl) GetShelves(ids []uint, userIDs []uint, bookIDs []uint) ([]shelf.Shelves, error) {
	result := []shelf.Shelves{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if len(ids) > 0 && ids != nil {
			tx = tx.Where("id IN ?", ids)
		}
		if len(userIDs) > 0 && userIDs != nil {
			tx = tx.Where("users_id IN ?", userIDs)
		}
		if len(bookIDs) > 0 && bookIDs != nil {
			tx = tx.Where("books_id IN ?", bookIDs)
		}
		if err := tx.Find(&result).Error; err != nil {
			log.Printf("error: GetShelves: %s\n", err.Error())
			return errors.New("error: get shelves failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetOneItemInShelf implements ShelfRepository.
func (s *shelfRepositoryImpl) GetOneItemInShelf(in *shelf.Shelves) (*shelf.Shelves, error) {
	if err := s.db.Where(&in).First(&in).Error; err != nil {
		log.Printf("error: GetOneItemInShelf: %s\n", err.Error())
		return nil, errors.New("error: get item in shelf failed")
	}
	return in, nil
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

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret, 
		"order", 
		"addbook", 
		resInByte,
		); err != nil {
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
