package bookUsecase

import (
	"fmt"
	"log"
	"math"
	"os"
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

// UpdateOneBookFile implements BookUsecase.
func (b *bookUsecaseImpl) UpdateOneBookFile(bookID uint, newFilePath string) error {
	condition := &book.Books{}
	condition.ID = bookID
	result, err := b.bookRepository.GetOneBook(condition)
	if err != nil {
		return err
	}

	if err := b.bookRepository.UpdateNonZeroBookFields(
		bookID, &book.Books{FilePath: newFilePath}); err != nil {
		return err
	}

	if err := os.Remove(result.FilePath); err != nil {
		log.Printf(
			"error: UpdateOneBookFile: fail to remove old book file at path: %s: %s\n",
			result.CoverImagePath, err.Error())
	}

	return nil
}

// UpdateOneBookCover implements BookUsecase.
func (b *bookUsecaseImpl) UpdateOneBookCover(bookID uint, newImagePath string) error {
	condition := &book.Books{}
	condition.ID = bookID
	result, err := b.bookRepository.GetOneBook(condition)
	if err != nil {
		return err
	}

	if err := b.bookRepository.UpdateNonZeroBookFields(bookID, &book.Books{CoverImagePath: newImagePath}); err != nil {
		return err
	}

	if result.CoverImagePath != "asset/image/bookCover/default/book-store_default_bookCover.png" {
		if err := os.Remove(result.CoverImagePath); err != nil {
			log.Printf("error: UpdateOneBookCover: fail to remove old image at path: %s: %s\n", result.CoverImagePath, err.Error())
		}
	}

	return nil
}

// UpdateOneBookDetail implements BookUsecase.
func (b *bookUsecaseImpl) UpdateOneBookDetail(req *book.UpdateBookDetailReq) error {

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
			return err
		}
		for _, v := range newTags {
			m[v.Name] = v.ID
		}
	}

	fmt.Printf("newTags: %v\n", newTags)

	for _, v := range notExistsTagIndex {
		req.Tags[v].ID = m[req.Tags[v].Name]
	}

	tags := []book.Tags{}
	for _, v := range req.Tags {
		t := book.Tags{}
		t.ID = v.ID
		tags = append(tags, t)
	}

	updateBook := &book.Books{
		Title:              req.Title,
		Price:              req.Price,
		Description:        req.Description,
		AuthorName:         req.AuthorName,
		Tags:               tags,
		IsAvailableInStore: req.IsAvailable,
	}

	if err := b.bookRepository.UpdateOneBookDetail(req.BookID, updateBook); err != nil {
		if len(newTags) > 0 {
			log.Println("rollBackCreateNewTags")
			b.rollBackCreateNewTags(newTags)
		}

		return err
	}
	return nil
}

// GetOneBookFileUrl implements BookUsecase.
func (b *bookUsecaseImpl) GetOneBookFilePath(bookID uint) (string, error) {
	condition := &book.Books{}
	condition.ID = bookID
	result, err := b.bookRepository.GetOneBook(condition)
	if err != nil {
		return "", err
	}
	return result.FilePath, nil
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
			CoverImagePath: b.convertBookCoverPathToUrl(v.CoverImagePath),
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
func (b *bookUsecaseImpl) GetOneBook(bookID uint, roleID uint) (*book.BookRes, error) {
	condition := &book.Books{}
	condition.ID = bookID
	if roleID != 1 {
		condition.IsAvailableInStore = true
	}
	result, err := b.bookRepository.GetOneBook(condition)
	if err != nil {
		return nil, err
	}

	return b.convertBooksToBookRes(result), nil
}

// SearchBooks implements BookUsecase.
func (b *bookUsecaseImpl) SearchBooks(req *book.SearchBooksReq, roleID uint) (*book.SearchBooksRes, error) {
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

	var isAvailable *bool
	if roleID == 1 {
		isAvailable = nil
	} else {
		b := true
		isAvailable = &b
	}

	books, count, err := b.bookRepository.SearchBook(
		*req.Limit,
		"created_at ASC",
		offset,
		req.Title,
		req.MaxPrice,
		req.MinPrice,
		req.AuthorName,
		isAvailable,
		tags,
	)
	if err != nil {
		return nil, err
	}

	// numOfPage := 1
	// if count < int64(*req.Limit) {
	// 	num
	// }
	numOfPage := float64(count) / float64(*req.Limit)
	numOfPage = math.Round(numOfPage)
	if numOfPage == 0 {
		numOfPage = 1
	}

	res := &book.SearchBooksRes{
		Pagination: models.PaginatieRes{
			Limit:           *req.Limit,
			LastVisiblePage: int64(numOfPage),
			HasNextPage:     (int64(*req.Page) < int64(numOfPage)),
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
		Title:              req.Title,
		Price:              req.Price,
		FilePath:           req.FilePath,
		Description:        req.Description,
		CoverImagePath:     req.CoverImagePath,
		AuthorName:         req.AuthorName,
		Tags:               tags,
		IsAvailableInStore: true,
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

func (b *bookUsecaseImpl) convertBookCoverPathToUrl(path string) string {
	splitedPath := strings.Split(path, "/")
	fileName := "/default/book-store_default_bookCover.png"
	if len(splitedPath) >= 1 {
		fileName = splitedPath[len(splitedPath)-1]
	}
	coverImgUrl := fmt.Sprintf("/book/cover/%s", fileName)
	return coverImgUrl
}

func (b *bookUsecaseImpl) convertBooksToBookRes(in *book.Books) *book.BookRes {
	coverImgUrl := b.convertBookCoverPathToUrl(in.CoverImagePath)
	tags := []book.BookTags{}
	for _, tag := range in.Tags {
		tags = append(tags, book.BookTags{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return &book.BookRes{
		ID:            in.ID,
		Title:         in.Title,
		Price:         in.Price,
		CoverImageUrl: coverImgUrl,
		AuthorName:    in.AuthorName,
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
		Description:   in.Description,
		Tags:          tags,
		IsAvailable:   in.IsAvailableInStore,
	}
}

func NewBookUsecaseImpl(cfg *config.Config, bookRepository bookRepository.BookRepository) BookUsecase {
	return &bookUsecaseImpl{cfg: cfg, bookRepository: bookRepository}
}
