package shelfUsecase

import (
	"errors"
	"testing"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockShelfRepository is a mock implementation of ShelfRepository
type MockShelfRepository struct {
	mock.Mock
}

func (m *MockShelfRepository) GetShelves(shelfIDs []uint, userIDs []uint, bookIDs []uint) ([]shelf.Shelves, error) { // Corrected return type
	args := m.Called(shelfIDs, userIDs, bookIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]shelf.Shelves), args.Error(1) // Corrected return type
}

func (m *MockShelfRepository) GetOneItemInShelf(in *shelf.Shelves) (*shelf.Shelves, error) {
	args := m.Called(in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*shelf.Shelves), args.Error(1)
}

func (m *MockShelfRepository) InsertUserBooks(shelves []shelf.Shelves) error {
	// It's important that this mock can modify the passed `shelves` to simulate ID assignment by GORM.
	// However, the direct modification of input slice elements (like assigning IDs)
	// should be handled carefully, often by the .Run() method in specific tests
	// if the test needs to assert those IDs.
	// For simplicity here, we'll assume the .Run block in the test will handle ID assignments if necessary.
	args := m.Called(shelves)
	return args.Error(0)
}

func (m *MockShelfRepository) DeleteUserBookInIDs(cfg *config.Config, shelfIDs []uint) error {
	args := m.Called(cfg, shelfIDs)
	return args.Error(0)
}

// Mock gRPC client calls / Producing messages
func (m *MockShelfRepository) FindBookInIds(bookUrl string, req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error) {
	args := m.Called(bookUrl, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookPb.FindBooksInIdsRes), args.Error(1)
}

func (m *MockShelfRepository) AddBookRes(cfg *config.Config, res *shelf.AddBooksRes) { // Removed error return
	m.Called(cfg, res) // Removed args assignment and return
}

// --- Test Cases ---

func TestNewShelfUsecaseImpl(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	assert.NotNil(t, uc)
}

func TestShelfUsecaseImpl_AddBooks_Success(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{} // Add Grpc.OrderUrl if needed for AddBookRes

	req := &shelf.AddBooksReq{
		OrderID:       1,
		UserID:        10,
		TransactionID: 100,
		BookIDs:       []uint{201, 202},
	}

	// Mock InsertUserBooks
	// Use MatchedBy to check content and .Run to simulate ID assignment
	mockRepo.On("InsertUserBooks", mock.MatchedBy(func(s []shelf.Shelves) bool {
		return len(s) == 2 && s[0].UsersID == req.UserID && s[0].BooksID == req.BookIDs[0] &&
			s[1].UsersID == req.UserID && s[1].BooksID == req.BookIDs[1]
	})).Run(func(args mock.Arguments) {
		shelvesArg := args.Get(0).([]shelf.Shelves)
		// Simulate GORM assigning IDs
		for i := range shelvesArg {
			shelvesArg[i].ID = uint(i + 1) // e.g., ID 1 and 2
		}
	}).Return(nil).Once()

	// Mock AddBookRes (the callback/producer call)
	expectedAddBooksRes := &shelf.AddBooksRes{
		OrderID:       req.OrderID,
		UserID:        req.UserID,
		TransactionID: req.TransactionID,
		ShelfIDs:      []uint{1, 2}, // IDs assigned by InsertUserBooks mock
		Error:         "",
	}
	mockRepo.On("AddBookRes", cfg, expectedAddBooksRes).Return(nil).Once()

	uc.AddBooks(cfg, req)

	mockRepo.AssertExpectations(t)
}

func TestShelfUsecaseImpl_AddBooks_InsertError(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &shelf.AddBooksReq{
		OrderID:       1,
		UserID:        10,
		TransactionID: 100,
		BookIDs:       []uint{201},
	}

	insertErr := errors.New("insert error")
	mockRepo.On("InsertUserBooks", mock.AnythingOfType("[]shelf.Shelves")).Return(insertErr).Once()

	// Expect AddBookRes to be called with the error
	expectedAddBooksResWithError := &shelf.AddBooksRes{
		OrderID:       req.OrderID,
		UserID:        req.UserID,
		TransactionID: req.TransactionID,
		Error:         insertErr.Error(),
		// ShelfIDs would be empty or nil
	}
	// Use mock.MatchedBy for AddBooksRes if exact ShelfIDs are tricky due to empty slice vs nil
	mockRepo.On("AddBookRes", cfg, mock.MatchedBy(func(res *shelf.AddBooksRes) bool {
		return res.OrderID == expectedAddBooksResWithError.OrderID &&
			res.UserID == expectedAddBooksResWithError.UserID &&
			res.TransactionID == expectedAddBooksResWithError.TransactionID &&
			res.Error == expectedAddBooksResWithError.Error &&
			len(res.ShelfIDs) == 0 // Check that ShelfIDs is empty
	})).Return(nil).Once()


	uc.AddBooks(cfg, req)

	mockRepo.AssertExpectations(t)
}

func TestShelfUsecaseImpl_RollBacksAddBook_Success(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &shelf.RollbackAddBooks{
		ShelfIDs: []uint{1, 2, 3},
	}

	mockRepo.On("DeleteUserBookInIDs", cfg, req.ShelfIDs).Return(nil).Once()

	uc.RollBacksAddBook(cfg, req) // Assuming the method name is RollBacksAddBook

	mockRepo.AssertExpectations(t)
}

func TestShelfUsecaseImpl_RollBacksAddBook_DeleteError(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &shelf.RollbackAddBooks{
		ShelfIDs: []uint{1, 2, 3},
	}

	// Even if DeleteUserBookInIDs returns an error, RollBacksAddBook currently doesn't return it.
	// It's a fire-and-forget from the usecase perspective.
	mockRepo.On("DeleteUserBookInIDs", cfg, req.ShelfIDs).Return(errors.New("delete error")).Once()

	uc.RollBacksAddBook(cfg, req)

	mockRepo.AssertExpectations(t)
}


func TestShelfUsecaseImpl_GetMyShelves_Success_NoBookIDQuery(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}}
	userID := uint(1)
	now := time.Now()

	// Shelves returned by repository
	shelvesFromRepo := []shelf.Shelves{ // Corrected type
		{Model: gorm.Model{ID: 1, CreatedAt: now}, UsersID: userID, BooksID: 101},
		{Model: gorm.Model{ID: 2, CreatedAt: now}, UsersID: userID, BooksID: 102},
	}
	mockRepo.On("GetShelves", ([]uint)(nil), []uint{userID}, ([]uint)(nil)).Return(shelvesFromRepo, nil).Once()

	// Books returned by Book gRPC service
	booksFromService := &bookPb.FindBooksInIdsRes{
		Book: []*bookPb.Book{
			{Id: 101, Title: "Book 101", AuthorName: "Author A", CoverImagePath: "/cover/101.jpg"},
			{Id: 102, Title: "Book 102", AuthorName: "Author B", CoverImagePath: "/cover/102.jpg"},
		},
	}
	expectedBookIDsReq := []uint64{101, 102} // Order might vary, use mock.MatchedBy
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, mock.MatchedBy(func(req *bookPb.FindBooksInIdsReq) bool {
		if len(req.Ids) != len(expectedBookIDsReq) { return false }
		m := make(map[uint64]bool)
		for _, id := range req.Ids { m[id] = true }
		for _, id := range expectedBookIDsReq { if !m[id] { return false } }
		return true
	})).Return(booksFromService, nil).Once()


	res, err := uc.GetMyShelves(cfg, userID, 0) // 0 for no specific bookIDQuery

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Shelves, 2)

	// Check shelf 1 (order depends on how uniqueBooks map iteration happens, so could be 101 or 102 first)
	// For simplicity, let's assume order or find by BookID
	var shelf101, shelf102 *shelf.ShelfRes
	for i := range res.Shelves {
		if res.Shelves[i].Book.ID == 101 {
			shelf101 = &res.Shelves[i]
		} else if res.Shelves[i].Book.ID == 102 {
			shelf102 = &res.Shelves[i]
		}
	}
	assert.NotNil(t, shelf101)
	assert.NotNil(t, shelf102)

	assert.Equal(t, uint(1), shelf101.ID) // Shelf ID
	assert.Equal(t, userID, shelf101.UserID)
	assert.Equal(t, uint(101), shelf101.Book.ID) // Book ID
	assert.Equal(t, "Book 101", shelf101.Book.Title)
	assert.Equal(t, "/cover/101.jpg", shelf101.Book.CoverImageUrl)


	assert.Equal(t, uint(2), shelf102.ID) // Shelf ID
	assert.Equal(t, "Book 102", shelf102.Book.Title)

	mockRepo.AssertExpectations(t)
}

func TestShelfUsecaseImpl_GetMyShelves_Success_WithBookIDQuery(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}}
	userID := uint(1)
	bookIDQuery := uint(101)

	shelvesFromRepo := []shelf.Shelves{ // Corrected type
		{Model: gorm.Model{ID: 1}, UsersID: userID, BooksID: bookIDQuery},
	}
	// Note: bookIDsDQuery in usecase becomes []uint{bookIDQuery}
	mockRepo.On("GetShelves", ([]uint)(nil), []uint{userID}, []uint{bookIDQuery}).Return(shelvesFromRepo, nil).Once()

	booksFromService := &bookPb.FindBooksInIdsRes{
		Book: []*bookPb.Book{{Id: 101, Title: "Book 101"}},
	}
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, &bookPb.FindBooksInIdsReq{Ids: []uint64{101}}).Return(booksFromService, nil).Once()

	res, err := uc.GetMyShelves(cfg, userID, bookIDQuery)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Shelves, 1)
	assert.Equal(t, bookIDQuery, res.Shelves[0].Book.ID)
	mockRepo.AssertExpectations(t)
}


func TestShelfUsecaseImpl_GetMyShelves_GetShelvesError(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	expectedErr := errors.New("get shelves error")

	mockRepo.On("GetShelves", ([]uint)(nil), []uint{1}, ([]uint)(nil)).Return(nil, expectedErr).Once()

	res, err := uc.GetMyShelves(cfg, 1, 0)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestShelfUsecaseImpl_GetMyShelves_FindBookInIdsError(t *testing.T) {
	mockRepo := new(MockShelfRepository)
	uc := NewShelfUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}}
	userID := uint(1)

	shelvesFromRepo := []shelf.Shelves{ // Corrected type
		{UsersID: userID, BooksID: 101},
	}
	mockRepo.On("GetShelves", ([]uint)(nil), []uint{userID}, ([]uint)(nil)).Return(shelvesFromRepo, nil).Once()

	expectedErr := errors.New("find book error")
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, mock.Anything).Return(nil, expectedErr).Once()

	res, err := uc.GetMyShelves(cfg, userID, 0)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}
