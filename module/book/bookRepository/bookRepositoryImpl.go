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

// UpdateNonZeroBookFields implements BookRepository.
func (b *bookRepositoryImpl) UpdateNonZeroBookFields(bookID uint, in *book.Books) error {
	condition := &book.Books{}
	condition.ID = bookID
	if err := b.db.Model(condition).Where(condition).Updates(&in).Error; err != nil {
		log.Printf("error: UpdateNonZeroBookFields: %s\n", err.Error())
		return errors.New("error: update book failed")
	}
	return nil
}

// UpdateOneBookDetail implements BookRepository.
func (b *bookRepositoryImpl) UpdateOneBookDetail(bookID uint, in *book.Books) error {
	condition := &book.Books{}
	condition.ID = bookID
	err := b.db.Transaction(func(tx *gorm.DB) error {

		// update book
		if err := tx.
			Model(&book.Books{}).
			Where(condition).Select("Title", "Price", "Description", "AuthorName", "IsAvailableInStore").Updates(in).Error; err != nil {
			return err
		}

		//Replace Associations
		if err := tx.Model(condition).Association("Tags").Replace(in.Tags); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("error: UpdateOneBookDetail: %s\n", err.Error())
		return errors.New("error: update book failed")
	}
	return nil
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
	if err := b.db.Preload("Tags").Where(&in).First(&in).Error; err != nil {
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
	isAvailable *bool,
	tagIDs []*uint,
) ([]book.Books, int64, error) {
	// conditions := []string{}
	// conditionsValue := []interface{}{}
	var count int64
	books := []book.Books{}
	err := b.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Model(&book.Books{})
		if title != "" {
			tx = tx.Where("title LIKE ?", fmt.Sprintf("%s%%", title))
			// conditions = append(conditions, "title LIKE ?")
			// conditionsValue = append(conditionsValue, fmt.Sprintf("%s%%", title))
		}

		if maxPrice != nil {
			tx = tx.Where("price <= ?", *maxPrice)

			// conditions = append(conditions, "price <= ?")
			// conditionsValue = append(conditionsValue, *maxPrice)
		}

		if minPrice != nil {
			tx = tx.Where("price >= ?", *minPrice)

			// conditions = append(conditions, "price >= ?")
			// conditionsValue = append(conditionsValue, *minPrice)
		}

		if authorName != "" {
			tx = tx.Where("author_name LIKE ?", fmt.Sprintf("%s%%", authorName))

			// conditions = append(conditions, "author_name LIKE ?")
			// conditionsValue = append(conditionsValue, fmt.Sprintf("%s%%", authorName))
		}
		if isAvailable != nil {
			tx = tx.Where("is_available_in_store = ?", *isAvailable)
		}
		for i, tagID := range tagIDs {

			alias := fmt.Sprintf("bt%d", i)
			tx = tx.Joins(
				fmt.Sprintf("JOIN books_tags %s ON books.id = %s.books_id AND %s.tags_id = %d",
					alias, alias, alias, *tagID),
			)
		}

		result := tx.
			Preload("Tags").
			Order(order).
			Count(&count).
			Offset(int(offest)).
			Limit(int(limit)).
			Find(&books)
		if result.Error != nil {
			log.Printf("error: GetBooks: %s\n", result.Error.Error())
			return errors.New("error: get books failed")
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}
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
