package bookRepository

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kritAsawaniramol/book-store/module/book"
	"gorm.io/gorm"
)

type bookRepositoryImpl struct {
	db *gorm.DB
}

// GetBooksInIDs implements BookRepository.
func (b *bookRepositoryImpl) GetBooksInIDs(ids []uint) ([]book.Books, int64, error) {
	books := []book.Books{}
	result := b.db.Where(ids).Find(&books)
	if result.Error != nil {
		log.Printf("error: GetBooksInIDs: %s\n", result.Error.Error())
		return nil, 0, errors.New("error: get books failed")
	}
	var count int64
	result.Count(&count)
	return books, count, nil
}

// GetTags implements BookRepository.
func (b *bookRepositoryImpl) GetTags(in *book.Tags) ([]book.Tags, error) {
	tags := []book.Tags{}
	if err := b.db.Where(in).Find(&tags).Error; err != nil {
		log.Printf("error: GetTags: %s\n", err.Error())
		return nil, errors.New("error: get tags failed")
	}
	return tags, nil
}

// GetOneBook implements BookRepository.
func (b *bookRepositoryImpl) GetOneBook(in *book.Books) (*book.Books, error) {
	if err := b.db.Preload("Tags").First(&in).Error; err != nil {
		log.Printf("error: GetOneBook: %s\n", err.Error())
		return nil, errors.New("error: get book failed")
	}
	return in, nil
}

// SearchBook implements BookRepository.
func (b *bookRepositoryImpl) SearchBook(
	limit int,
	order string,
	offest int,
	title string,
	maxPrice *uint,
	minPrice *uint,
	authorName string,
	tagIDs []*uint,
) ([]book.Books, int64, error) {
	conditions := []string{}
	conditionsValue := []interface{}{}
	if title != "" {
		conditions = append(conditions, "title LIKE ?")
		conditionsValue = append(conditionsValue, fmt.Sprintf("%s%%", title))
	}

	if maxPrice != nil {
		conditions = append(conditions, "price <= ?")
		conditionsValue = append(conditionsValue, *maxPrice)
	}

	if minPrice != nil {
		conditions = append(conditions, "price >= ?")
		conditionsValue = append(conditionsValue, *minPrice)
	}

	if authorName != "" {
		conditions = append(conditions, "author_name LIKE ?")
		conditionsValue = append(conditionsValue, fmt.Sprintf("%s%%", authorName))
	}
	conditionsStr := strings.Join(conditions, " AND ")
	conds := make([]interface{}, 0)
	conds = append(conds, conditionsStr)
	conds = append(conds, conditionsValue...)

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

	result := tx.
		Order(order).
		Offset(int(offest)).
		Limit(int(limit)).
		Preload("Tags").
		// Find(&books, c...)
		Find(&books, conds...)
	if result.Error != nil {
		log.Printf("error: GetBooks: %s\n", result.Error.Error())
		return books, 0, errors.New("error: get books failed")
	}
	result.Count(&count)
	fmt.Printf("count: %v\n", count)
	return books, count, nil
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
