package bookUsecase

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/models"
	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/book/bookRepository"
)

type bookUsecaseImpl struct {
	bookRepository bookRepository.BookRepository
	cfg            *config.Config
}

// FindBookInIDs implements BookUsecase.
func (b *bookUsecaseImpl) FindBookInIDs(req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error) {
	ids := []uint{}
	for _, id := range req.Ids {
		ids = append(ids, uint(id))
	}

	books, _, err := b.bookRepository.GetBooksInIDs(ids)
	if err != nil {
		return nil, err
	}
	bookRes := []*bookPb.Book{}
	for _, v := range books {
		tags := []*bookPb.Tags{}
		for _, tag := range v.Tags {
			tags = append(tags, &bookPb.Tags{
				Id:   uint64(tag.ID),
				Name: tag.Name,
			})
		}
		bookRes = append(bookRes, &bookPb.Book{
			Id:             uint64(v.ID),
			Title:          v.Title,
			Price:          uint64(v.Price),
			FilePath:       v.FilePath,
			CoverImagePath: v.CoverImagePath,
			AuthorName:     v.AuthorName,
			Tags:           tags,
		})
	}
	return &bookPb.FindBooksInIdsRes{
		Book: bookRes,
	}, nil
}

// GetTags implements BookUsecase.
func (b *bookUsecaseImpl) GetTags() ([]book.BookTags, error) {
	tags, err := b.bookRepository.GetTags(&book.Tags{})
	if err != nil {
		return nil, err
	}

	res := []book.BookTags{}
	for _, t := range tags {
		res = append(res, book.BookTags{
			ID:   t.ID,
			Name: t.Name,
		})
	}
	return res, nil
}

// GetOneBook implements BookUsecase.
func (b *bookUsecaseImpl) GetOneBook(bookID uint) (*book.BookRes, error) {
	condition := &book.Books{}
	condition.ID = bookID
	result, err := b.bookRepository.GetOneBook(condition)
	if err != nil {
		return nil, err
	}

	return b.convertBooksToBookRes(result), nil
}

// SearchBooks implements BookUsecase.
func (b *bookUsecaseImpl) SearchBooks(req *book.SearchBooksReq) (*book.SearchBooksRes, error) {
	if req.Page == nil {
		var page uint = 1
		req.Page = &page
	}

	if req.Limit == nil {
		var limit int = 25
		req.Limit = &limit
	}

	if req.MinPrice == nil {
		var min uint = 0
		req.MinPrice = &min
	}

	tags := []*uint{}
	if req.TagIDs != "" {
		tagIDs := strings.Split(req.TagIDs, ",")
		for _, v := range tagIDs {
			tID, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, err
			}
			tIDUint := uint(tID)
			tIDPointer := &tIDUint
			tags = append(tags, tIDPointer)
		}
	}
	offset := (int(*req.Page) - 1) * *req.Limit

	books, count, err := b.bookRepository.SearchBook(
		*req.Limit,
		"created_at ASC",
		offset,
		req.Title,
		req.MaxPrice,
		req.MinPrice,
		req.AuthorName,
		tags,
	)
	if err != nil {
		return nil, err
	}

	numOfPage := count / int64(*req.Limit)
	if numOfPage == 0 {
		numOfPage = 1
	}

	res := &book.SearchBooksRes{
		Pagination: models.PaginatieRes{
			Limit:           *req.Limit,
			LastVisiblePage: numOfPage,
			HasNextPage:     (int64(*req.Page) < numOfPage),
			Total:           count,
		},
	}

	for _, v := range books {
		res.Books = append(res.Books, *b.convertBooksToBookRes(&v))
	}
	return res, nil
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

	for _, v := range notExistsTagIndex {
		req.Tags[v].ID = m[req.Tags[v].Name]
	}

	tags := []book.Tags{}
	for _, v := range req.Tags {
		t := book.Tags{}
		t.ID = v.ID
		tags = append(tags, t)
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

func (b *bookUsecaseImpl) convertBooksToBookRes(in *book.Books) *book.BookRes {
	splitedPath := strings.Split(in.CoverImagePath, "/")
	fileName := "/default/book-store_default_bookCover.png"
	if len(splitedPath) >= 1 {
		fileName = splitedPath[len(splitedPath)-1]
	}

	tags := []book.BookTags{}
	for _, tag := range in.Tags {
		tags = append(tags, book.BookTags{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	coverImgUrl := fmt.Sprintf("%s:%d/book/cover/%s", b.cfg.App.Host, b.cfg.App.Port, fileName)
	return &book.BookRes{
		Title:         in.Title,
		Price:         in.Price,
		CoverImageUrl: coverImgUrl,
		AuthorName:    in.AuthorName,
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
		Tags:          tags,
	}
}

func NewBookUsecaseImpl(cfg *config.Config, bookRepository bookRepository.BookRepository) BookUsecase {
	return &bookUsecaseImpl{cfg: cfg, bookRepository: bookRepository}
}
