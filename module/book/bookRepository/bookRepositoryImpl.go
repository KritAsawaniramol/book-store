package bookRepository

import (
	"errors"
	"fmt"
	"log"

	"github.com/kritAsawaniramol/book-store/module/book"
	"gorm.io/gorm"
)

type bookRepositoryImpl struct {
	db *gorm.DB
}

// GetOneBook implements BookRepository.
func (b *bookRepositoryImpl) GetOneBook(in *book.Books) (*book.Books, error) {
	if err := b.db.Preload("Tags").First(&in).Error; err != nil {
		log.Printf("error: GetOneBook: %s\n", err.Error())
		return nil, errors.New("error: get book failed")
	}
	return in, nil
}

// GetBooks implements BookRepository.
func (b *bookRepositoryImpl) GetBooks(
	limit int,
	order interface{},
	offest uint,
	tagIDs []uint,
	conditions ...interface{},
) ([]book.Books, int64, error) {
	books := []book.Books{}
	var count int64
	tx := b.db.Model(&book.Books{})
	for i, tagID := range tagIDs {
		alias := fmt.Sprintf("bt%d", i)
		tx = tx.Joins(
			fmt.Sprintf("JOIN books_tags %s ON books.id = %s.books_id AND %s.tags_id = %d",
				alias, alias, alias, tagID),
		)
	}
	result := tx.Order(order).Offset(int(offest)).Limit(int(limit)).Preload("Tags").Find(&books, conditions...)

	if result.Error != nil {
		log.Printf("error: GetBooks: %s\n", result.Error.Error())
		return books, 0, errors.New("error: get books failed")
	}

	result.Count(&count)
	fmt.Printf("count: %v\n", count)
	return books, count, nil
}

func (b *bookRepositoryImpl) Test() {
	books := []book.Books{}
	b.db.
		Joins("JOIN books_tags bt ON bt.book_id = books.id").
		Where("bt.tag_id IN ?", []uint{1, 3}).
		Find(books)
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
