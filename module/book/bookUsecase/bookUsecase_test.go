package bookUsecase

import (
	"errors"
	"testing"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/models" // Added back models import
	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm" // Added gorm import
)

// MockBookRepository is a mock implementation of BookRepository
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) CreateOneBook(b *book.Books) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockBookRepository) GetOneBook(condition *book.Books) (*book.Books, error) {
	args := m.Called(condition)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*book.Books), args.Error(1)
}

func (m *MockBookRepository) GetBooksInIDs(ids []uint) ([]book.Books, int64, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]book.Books), args.Get(1).(int64), args.Error(2) // Corrected to use Get(1).(int64)
}

// Corrected SearchBook signature
func (m *MockBookRepository) SearchBook(limit int, order string, offset int, title string, maxPrice *uint, minPrice *uint, authorName string, isAvailable *bool, tagIDs []*uint) ([]book.Books, int64, error) {
	args := m.Called(limit, order, offset, title, maxPrice, minPrice, authorName, isAvailable, tagIDs)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]book.Books), args.Get(1).(int64), args.Error(2) // Corrected to use Get(1).(int64)
}

func (m *MockBookRepository) UpdateOneBookDetail(bookID uint, updatedBook *book.Books) error {
	args := m.Called(bookID, updatedBook)
	return args.Error(0)
}

func (m *MockBookRepository) UpdateNonZeroBookFields(bookID uint, updatedFields *book.Books) error {
	args := m.Called(bookID, updatedFields)
	return args.Error(0)
}

func (m *MockBookRepository) DeleteOneBook(bookID uint) error {
	args := m.Called(bookID)
	return args.Error(0)
}

func (m *MockBookRepository) CreateTags(tags []book.Tags) error {
	args := m.Called(tags)
	// ID assignment should be handled by the test's .Run() block if specific IDs are needed.
	// The .Run() block for CreateTags in TestUpdateOneBookDetail_Success and TestCreateOneBook_Success
	// will set the IDs as expected for those tests.
	return args.Error(0)
}

func (m *MockBookRepository) GetTags(condition *book.Tags) ([]book.Tags, error) {
	args := m.Called(condition)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]book.Tags), args.Error(1)
}

func (m *MockBookRepository) DeleteTags(tags []book.Tags) error {
	args := m.Called(tags)
	return args.Error(0)
}

func TestNewBookUsecaseImpl(t *testing.T) {
	mockRepo := new(MockBookRepository)
	mockCfg := &config.Config{}
	uc := NewBookUsecaseImpl(mockCfg, mockRepo)
	assert.NotNil(t, uc)
}

func TestBookUsecaseImpl_CreateOneBook_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.CreateBookReq{
		Title:       "Test Book",
		Price:       100,
		FilePath:    "/path/to/book.pdf",
		Description: "A great book",
		CoverImagePath: "/path/to/cover.jpg",
		AuthorName:  "Test Author",
		Tags: []book.BookTags{
			{Name: "Fiction"}, // New tag
			{ID: 2, Name: "Existing Tag"}, // This is correct for book.BookTags
		},
	}

	// Mock expectations
	mockRepo.On("CreateTags", mock.AnythingOfType("[]book.Tags")).Return(nil).Once()
	mockRepo.On("CreateOneBook", mock.AnythingOfType("*book.Books")).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*book.Books)
		arg.ID = 1 // Simulate database assigning an ID
	}).Return(nil).Once()

	bookID, err := uc.CreateOneBook(req)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), bookID)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_CreateOneBook_CreateTagsError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.CreateBookReq{
		Title: "Test Book",
		Tags:  []book.BookTags{{Name: "New Tag"}}, // Corrected type to book.BookTags
	}

	expectedErr := errors.New("create tags error")
	mockRepo.On("CreateTags", mock.AnythingOfType("[]book.Tags")).Return(expectedErr).Once()

	bookID, err := uc.CreateOneBook(req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, uint(0), bookID)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_CreateOneBook_CreateBookError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.CreateBookReq{
		Title: "Test Book",
		Tags:  []book.BookTags{{ID: 1, Name: "Existing Tag"}}, // Corrected to use book.BookTags with direct ID
	}

	expectedErr := errors.New("create book error")
	// No call to CreateTags expected
	mockRepo.On("CreateOneBook", mock.AnythingOfType("*book.Books")).Return(expectedErr).Once()

	bookID, err := uc.CreateOneBook(req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, uint(0), bookID)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_CreateOneBook_CreateBookErrorWithRollback(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.CreateBookReq{
		Title: "Test Book",
		Tags:  []book.BookTags{{Name: "New Tag For Rollback"}}, // Corrected type to book.BookTags
	}

	createBookErr := errors.New("create book error")
	// Expect CreateTags to be called
	mockRepo.On("CreateTags", mock.MatchedBy(func(tags []book.Tags) bool {
		return len(tags) == 1 && tags[0].Name == "New Tag For Rollback"
	})).Return(nil).Once()
	// Expect CreateOneBook to be called and fail
	mockRepo.On("CreateOneBook", mock.AnythingOfType("*book.Books")).Return(createBookErr).Once()
	// Expect DeleteTags (rollback) to be called
	mockRepo.On("DeleteTags", mock.MatchedBy(func(tags []book.Tags) bool {
		return len(tags) == 1 && tags[0].Name == "New Tag For Rollback" // ID would have been set by mock CreateTags
	})).Return(nil).Once()


	bookID, err := uc.CreateOneBook(req)

	assert.Error(t, err)
	assert.Equal(t, createBookErr, err)
	assert.Equal(t, uint(0), bookID)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_GetOneBook_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo) // Simplified config

	now := time.Now()
	expectedBookEntity := &book.Books{
		Model:              gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		Title:              "Test Book",
		Price:              100,
		CoverImagePath:     "asset/image/bookCover/default/book-store_default_bookCover.png",
		AuthorName:         "Test Author",
		Description:        "A great book",
		Tags:               []book.Tags{{Model: gorm.Model{ID: 1}, Name: "Fiction"}},
		IsAvailableInStore: true,
	}

	// Mock for GetOneBook should match how it's called in usecase
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:1}, IsAvailableInStore: true}).Return(expectedBookEntity, nil).Once()
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:1}}).Return(expectedBookEntity, nil).Once()


	// Test for non-admin user (roleID = 2)
	res, err := uc.GetOneBook(1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedBookEntity.ID, res.ID) // ID is from gorm.Model
	assert.Equal(t, expectedBookEntity.Title, res.Title)
	assert.Equal(t, "/book/cover/book-store_default_bookCover.png", res.CoverImageUrl) // Expect path


	// Test for admin user (roleID = 1)
	resAdmin, errAdmin := uc.GetOneBook(1, 1)
	assert.NoError(t, errAdmin)
	assert.NotNil(t, resAdmin)
	assert.Equal(t, expectedBookEntity.ID, resAdmin.ID)


	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_GetOneBook_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	expectedErr := errors.New("book not found")
	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:1}, IsAvailableInStore: true}).Return(nil, expectedErr).Once()

	res, err := uc.GetOneBook(1, 2) // roleID = 2 (non-admin)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_GetTags_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	expectedRepoTags := []book.Tags{ // Tags returned by repository
		{Model: gorm.Model{ID: 1}, Name: "Fiction"},
		{Model: gorm.Model{ID: 2}, Name: "Science"},
	}
	mockRepo.On("GetTags", &book.Tags{}).Return(expectedRepoTags, nil).Once()

	res, err := uc.GetTags() // res here is []book.BookTags

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 2)
	assert.Equal(t, expectedRepoTags[0].ID, res[0].ID) // Comparing ID from gorm.Model
	assert.Equal(t, expectedRepoTags[0].Name, res[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_GetTags_Error(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	expectedErr := errors.New("failed to get tags")
	mockRepo.On("GetTags", &book.Tags{}).Return(nil, expectedErr).Once()

	res, err := uc.GetTags()

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_FindBookInIDs_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo) // Simplified config

	req := &bookPb.FindBooksInIdsReq{Ids: []uint64{1, 2}}
	repoResBooks := []book.Books{ // Entities from repo
		{Model: gorm.Model{ID: 1}, Title: "Book 1", CoverImagePath: "path/to/cover1.jpg", Tags: []book.Tags{{Model: gorm.Model{ID:1}, Name: "Tag1"}}},
		{Model: gorm.Model{ID: 2}, Title: "Book 2", CoverImagePath: "path/to/cover2.jpg", Tags: []book.Tags{{Model: gorm.Model{ID:2}, Name: "Tag2"}}},
	}

	mockRepo.On("GetBooksInIDs", []uint{1, 2}).Return(repoResBooks, int64(2), nil).Once()

	res, err := uc.FindBookInIDs(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Book, 2)
	assert.Equal(t, uint64(1), res.Book[0].Id)
	assert.Equal(t, "/book/cover/cover1.jpg", res.Book[0].CoverImagePath) // Expect path
	assert.Equal(t, uint64(2), res.Book[1].Id)
	assert.Equal(t, "/book/cover/cover2.jpg", res.Book[1].CoverImagePath) // Expect path
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_FindBookInIDs_Error(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &bookPb.FindBooksInIdsReq{Ids: []uint64{1}}
	expectedErr := errors.New("db error")

	mockRepo.On("GetBooksInIDs", []uint{1}).Return(nil, int64(0), expectedErr).Once()

	res, err := uc.FindBookInIDs(req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}


func TestBookUsecaseImpl_GetOneBookFilePath_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	expectedPath := "/path/to/book.pdf"
	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:1}}).Return(&book.Books{Model: gorm.Model{ID:1}, FilePath: expectedPath}, nil).Once()

	filePath, err := uc.GetOneBookFilePath(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, filePath)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_GetOneBookFilePath_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	expectedErr := errors.New("book not found")
	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:1}}).Return(nil, expectedErr).Once()

	filePath, err := uc.GetOneBookFilePath(1)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, filePath)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_UpdateOneBookDetail_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.UpdateBookDetailReq{
		BookID:      1,
		Title:       "Updated Title",
		Price:       150,
		Description: "Updated Description",
		AuthorName:  "Updated Author",
		Tags: []book.BookTags{ // Corrected type to book.BookTags
			{Name: "New Tag"},         // New tag
			{ID: 2, Name: "Existing Tag"}, // Existing tag
		},
		IsAvailable: true,
	}

	// Mock CreateTags for the new tag
	mockRepo.On("CreateTags", mock.MatchedBy(func(tags []book.Tags) bool {
		return len(tags) == 1 && tags[0].Name == "New Tag"
	})).Run(func(args mock.Arguments) {
		tags := args.Get(0).([]book.Tags)
		tags[0].ID = 3 // Simulate assigning ID to new tag
	}).Return(nil).Once()

	// Mock UpdateOneBookDetail
	mockRepo.On("UpdateOneBookDetail", uint(1), mock.MatchedBy(func(b *book.Books) bool {
		return b.Title == req.Title && len(b.Tags) == 2 && b.Tags[0].ID == 3 && b.Tags[1].ID == 2 // Verify new and existing tags are included with correct IDs
	})).Return(nil).Once()

	err := uc.UpdateOneBookDetail(req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_UpdateOneBookDetail_CreateTagsError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.UpdateBookDetailReq{
		BookID: 1,
		Tags:   []book.BookTags{{Name: "New Tag"}}, // Corrected type to book.BookTags
	}
	expectedErr := errors.New("create tags error")
	mockRepo.On("CreateTags", mock.AnythingOfType("[]book.Tags")).Return(expectedErr).Once()

	err := uc.UpdateOneBookDetail(req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_UpdateOneBookDetail_UpdateErrorWithRollback(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.UpdateBookDetailReq{
		BookID: 1,
		Title:  "Updated Title",
		Tags:   []book.BookTags{{Name: "Another New Tag"}}, // Corrected type to book.BookTags
	}

	updateErr := errors.New("update book detail error")

	// Expect CreateTags to be called for "Another New Tag"
	mockRepo.On("CreateTags", mock.MatchedBy(func(tags []book.Tags) bool {
		return len(tags) == 1 && tags[0].Name == "Another New Tag"
	})).Run(func(args mock.Arguments) {
		tags := args.Get(0).([]book.Tags)
		tags[0].ID = 5 // Simulate assigning ID
	}).Return(nil).Once()

	// Expect UpdateOneBookDetail to fail
	mockRepo.On("UpdateOneBookDetail", uint(1), mock.AnythingOfType("*book.Books")).Return(updateErr).Once()

	// Expect DeleteTags (rollback) to be called for "Another New Tag"
	mockRepo.On("DeleteTags", mock.MatchedBy(func(tags []book.Tags) bool {
		// The mock CreateTags should have set an ID, e.g., 5
		return len(tags) == 1 && tags[0].Name == "Another New Tag" && tags[0].ID == 5
	})).Return(nil).Once()

	err := uc.UpdateOneBookDetail(req)

	assert.Error(t, err)
	assert.Equal(t, updateErr, err)
	mockRepo.AssertExpectations(t)
}


func TestBookUsecaseImpl_UpdateOneBookCover_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	bookID := uint(1)
	newImagePath := "/path/to/new_cover.jpg"
	oldImagePath := "/path/to/old_cover.jpg" // Simulate an existing, non-default image

	// Mock GetOneBook to return a book with an old image path
	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:bookID}}).Return(&book.Books{Model: gorm.Model{ID:bookID}, CoverImagePath: oldImagePath}, nil).Once()
	// Mock UpdateNonZeroBookFields to succeed
	mockRepo.On("UpdateNonZeroBookFields", bookID, &book.Books{CoverImagePath: newImagePath}).Return(nil).Once()
	// os.Remove will be called, we are not mocking it here but in a real scenario with file system interaction, you might.
	// For now, we assume it works or handle its potential error if critical.

	err := uc.UpdateOneBookCover(bookID, newImagePath)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	// Add assertion for os.Remove if you mock it.
	// For now, we check that no error was returned from the usecase.
}

func TestBookUsecaseImpl_UpdateOneBookCover_GetBookError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	bookID := uint(1)
	newImagePath := "/path/to/new_cover.jpg"
	expectedErr := errors.New("get book error")

	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:bookID}}).Return(nil, expectedErr).Once()

	err := uc.UpdateOneBookCover(bookID, newImagePath)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_UpdateOneBookCover_UpdateError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	bookID := uint(1)
	newImagePath := "/path/to/new_cover.jpg"
	oldImagePath := "asset/image/bookCover/default/book-store_default_bookCover.png" // Default image
	expectedErr := errors.New("update error")

	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:bookID}}).Return(&book.Books{Model: gorm.Model{ID:bookID}, CoverImagePath: oldImagePath}, nil).Once()
	mockRepo.On("UpdateNonZeroBookFields", bookID, &book.Books{CoverImagePath: newImagePath}).Return(expectedErr).Once()

	err := uc.UpdateOneBookCover(bookID, newImagePath)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
	// os.Remove should not be called if update fails or if old path is default
}


func TestBookUsecaseImpl_UpdateOneBookFile_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	bookID := uint(1)
	newFilePath := "/path/to/new_book.pdf"
	oldFilePath := "/path/to/old_book.pdf"

	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:bookID}}).Return(&book.Books{Model: gorm.Model{ID:bookID}, FilePath: oldFilePath}, nil).Once()
	mockRepo.On("UpdateNonZeroBookFields", bookID, &book.Books{FilePath: newFilePath}).Return(nil).Once()
	// os.Remove will be called for oldFilePath.

	err := uc.UpdateOneBookFile(bookID, newFilePath)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_UpdateOneBookFile_GetBookError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	bookID := uint(1)
	newFilePath := "/path/to/new_book.pdf"
	expectedErr := errors.New("get book error")

	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:bookID}}).Return(nil, expectedErr).Once()

	err := uc.UpdateOneBookFile(bookID, newFilePath)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_UpdateOneBookFile_UpdateError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	bookID := uint(1)
	newFilePath := "/path/to/new_book.pdf"
	oldFilePath := "/path/to/old_book.pdf"
	expectedErr := errors.New("update error")

	// Corrected mock to use gorm.Model for ID
	mockRepo.On("GetOneBook", &book.Books{Model: gorm.Model{ID:bookID}}).Return(&book.Books{Model: gorm.Model{ID:bookID}, FilePath: oldFilePath}, nil).Once()
	mockRepo.On("UpdateNonZeroBookFields", bookID, &book.Books{FilePath: newFilePath}).Return(expectedErr).Once()

	err := uc.UpdateOneBookFile(bookID, newFilePath)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_SearchBooks_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo) // Simplified config

	page := uint(1)
	limit := 10
	minPrice := uint(50)
	maxPrice := uint(200)
	title := "Test"
	author := "Author"
	tagIDs := "1,2"


	req := &book.SearchBooksReq{
		PaginatieReq: models.PaginatieReq{Page: &page, Limit: &limit}, // Corrected: Use embedded struct
		MinPrice:     &minPrice,
		MaxPrice:     &maxPrice,
		Title:        title,  // Corrected: Dereference pointer
		AuthorName:   author, // Corrected: Dereference pointer
		TagIDs:       tagIDs,
	}
	roleID := uint(2) // Non-admin

	expectedBooksFromRepo := []book.Books{ // Entities from repo
		{Model: gorm.Model{ID: 1}, Title: "Test Book 1", Price: 100, CoverImagePath: "path/cover1.jpg", Tags: []book.Tags{{Model: gorm.Model{ID:1}, Name:"Tag1"}}},
	}
	expectedCount := int64(1)
	isAvailable := true // for non-admin

	parsedTagIDs := []*uint{}
	tID1, _ := uint(1), uint(2)
	parsedTagIDs = append(parsedTagIDs, &tID1, &tID1)


	mockRepo.On(
		"SearchBook",
		limit,
		"created_at ASC",
		0, // offset for page 1
		title, // Corrected: pass string directly
		&maxPrice,
		&minPrice,
		author,
		&isAvailable,
		mock.AnythingOfType("[]*uint"),
	).Return(expectedBooksFromRepo, expectedCount, nil).Once()

	res, err := uc.SearchBooks(req, roleID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Books, 1)
	assert.Equal(t, expectedCount, res.Pagination.Total)
	assert.Equal(t, "/book/cover/cover1.jpg", res.Books[0].CoverImageUrl) // Expect path
	mockRepo.AssertExpectations(t)
}

func TestBookUsecaseImpl_SearchBooks_DefaultPaginationAndAdmin(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo) // Simplified config

	// Nil Page, Limit, MinPrice
	req := &book.SearchBooksReq{}
	roleIDAdmin := uint(1) // Admin

	expectedBooksFromRepo := []book.Books{}
	expectedCount := int64(0)

	// For admin, isAvailable should be nil
	// Default limit 25, page 1 (offset 0), minPrice 0
	var defaultMinPrice uint = 0
	mockRepo.On(
		"SearchBook",
		25, // default limit
		"created_at ASC",
		0,  // default offset
		"", // title (empty string for string type)
		(*uint)(nil),   // maxPrice
		&defaultMinPrice, // default minPrice
		"", // authorName (empty string for string type)
		(*bool)(nil),
		mock.AnythingOfType("[]*uint"),
	).Return(expectedBooksFromRepo, expectedCount, nil).Once()

	res, err := uc.SearchBooks(req, roleIDAdmin)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 25, res.Pagination.Limit)
	assert.Equal(t, int64(1), res.Pagination.LastVisiblePage) // numOfPage = math.Round(0/25) = 0, then becomes 1
	mockRepo.AssertExpectations(t)
}


func TestBookUsecaseImpl_SearchBooks_TagParsingError(t *testing.T) {
	mockRepo := new(MockBookRepository) // Not used, but good practice
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)

	req := &book.SearchBooksReq{
		TagIDs: "1,abc,2", // "abc" will cause parsing error
	}
	roleID := uint(2)

	res, err := uc.SearchBooks(req, roleID)

	assert.Error(t, err) // Expect an error due to parsing "abc"
	assert.Nil(t, res)
	// No mockRepo expectations as it should fail before DB call
}


func TestBookUsecaseImpl_SearchBooks_DbError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	uc := NewBookUsecaseImpl(&config.Config{}, mockRepo)
	page := uint(1)
	limit := 10
	// Corrected: Use embedded struct for Page and Limit
	req := &book.SearchBooksReq{PaginatieReq: models.PaginatieReq{Page: &page, Limit: &limit}}
	roleID := uint(2)
	expectedErr := errors.New("search db error")
	isAvailable := true
	var minPriceDefault uint = 0


	mockRepo.On(
		"SearchBook",
		limit,
		"created_at ASC",
		0,
		"", // title (empty string for string type)
		(*uint)(nil),
		&minPriceDefault,
		"", // authorName (empty string for string type)
		&isAvailable,
		mock.AnythingOfType("[]*uint"),
	).Return(nil, int64(0), expectedErr).Once()

	res, err := uc.SearchBooks(req, roleID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}
