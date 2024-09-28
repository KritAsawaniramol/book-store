package bookRepository

import (
	"errors"
	"log"

	"github.com/kritAsawaniramol/book-store/module/book"
	"gorm.io/gorm"
)

type bookRepositoryImpl struct {
	db *gorm.DB
}

// DeleteTags implements BookRepository.
func (b *bookRepositoryImpl) DeleteTags(in []book.Tags) error {
	if err := b.db.Delete(&in).Error; err != nil {
		log.Printf("error: DeleteTags: %s\n", err.Error())
		return errors.New("error: delete tags failed")
	}
	return nil
}

// CreateTags implements BookRepository.
func (b *bookRepositoryImpl) CreateTags(in []book.Tags) error {
	if err := b.db.Create(&in).Error; err != nil {
		log.Printf("error: CreateTags: %s\n", err.Error())
		return errors.New("error: create tags failed")
	}
	return nil
}

// CreateOneBook implements BookRepository.
func (b *bookRepositoryImpl) CreateOneBook(in *book.Books) error {
	if err := b.db.Create(&in).Error; err != nil {
		log.Printf("error: CreateOneBook: %s\n", err.Error())
		return errors.New("error: create book failed")
	}
	return nil
}

func NewBookRepositoryImpl(db *gorm.DB) BookRepository {
	return &bookRepositoryImpl{db: db}
}
