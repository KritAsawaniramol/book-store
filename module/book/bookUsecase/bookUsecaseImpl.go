package bookUsecase

import (
	"log"

	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/book/bookRepository"
	"github.com/kritAsawaniramol/book-store/util"
)

type bookUsecaseImpl struct {
	bookRepository bookRepository.BookRepository
}

// CreateOneBook implements BookUsecase.
func (b *bookUsecaseImpl) CreateOneBook(req *book.CreateBookReq) (uint, error) {
	newTags := []book.Tags{}
	notExistsTagIndex := []int{}
	for idx, v := range req.Tags {
		if v.ID == 0 {
			newTags = append(newTags, book.Tags{
				Name: v.Name,
			})
			notExistsTagIndex = append(notExistsTagIndex, idx)
		}
	}
	m := map[string]uint{}
	if len(newTags) > 0 {
		err := b.bookRepository.CreateTags(newTags)
		if err != nil {
			return 0, err
		}
		for _, v := range newTags {
			m[v.Name] = v.ID
		}
	}
	util.PrintObjInJson(newTags)

	for _, v := range notExistsTagIndex {
		req.Tags[v].ID = m[req.Tags[v].Name]
	}

	tags := []book.BooksTags{}
	for _, v := range req.Tags {
		t := book.Tags{}
		t.ID = v.ID
		tags = append(tags, book.BooksTags{
			TagsID: v.ID,
		})
	}

	newBook := &book.Books{
		Title:          req.Title,
		Price:          req.Price,
		FilePath:       req.FilePath,
		CoverImagePath: req.CoverImagePath,
		AuthorName:     req.AuthorName,
		Tags:           tags,
	}

	if err := b.bookRepository.CreateOneBook(newBook); err != nil {
		if len(newTags) > 0 {
			log.Println("rollBackCreateNewTags")
			b.rollBackCreateNewTags(newTags)
		}
		return 0, err
	}
	return newBook.ID, nil
}

func (b *bookUsecaseImpl) rollBackCreateNewTags(newTags []book.Tags) error {
	err := b.bookRepository.DeleteTags(newTags)
	if err != nil {
		return err
	}
	return nil
}

func NewBookUsecaseImpl(bookRepository bookRepository.BookRepository) BookUsecase {
	return &bookUsecaseImpl{bookRepository: bookRepository}
}
