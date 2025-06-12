package orderUsecase

import (
	"errors"
	"testing"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/book/bookPb"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(req *order.Orders) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrders(orderIDs []uint, userIDs []uint, status []string, isPreload bool) ([]order.Orders, error) { // Return []order.Orders
	args := m.Called(orderIDs, userIDs, status, isPreload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]order.Orders), args.Error(1) // Return []order.Orders
}

func (m *MockOrderRepository) FindOrdersWithBookDetail(filter *order.Orders, status []string) ([]order.Orders, error) { // Return []order.Orders
	args := m.Called(filter, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]order.Orders), args.Error(1) // Return []order.Orders
}

func (m *MockOrderRepository) UpdateOrderByID(orderID uint, updatedFields *order.Orders) error {
	args := m.Called(orderID, updatedFields)
	return args.Error(0)
}

// Mock gRPC client calls
func (m *MockOrderRepository) FindBookInIds(bookUrl string, req *bookPb.FindBooksInIdsReq) (*bookPb.FindBooksInIdsRes, error) {
	args := m.Called(bookUrl, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookPb.FindBooksInIdsRes), args.Error(1)
}

func (m *MockOrderRepository) DecreaseUserMoney(cfg *config.Config, req *user.BuyBookReq) error {
	args := m.Called(cfg, req)
	return args.Error(0)
}

func (m *MockOrderRepository) RollbackUserTransaction(cfg *config.Config, req *user.RollbackUserTransactionReq) error {
	args := m.Called(cfg, req)
	return args.Error(0)
}

func (m *MockOrderRepository) AddBookToShelf(cfg *config.Config, req *shelf.AddBooksReq) error {
	args := m.Called(cfg, req)
	return args.Error(0)
}

func (m *MockOrderRepository) RollbackAddBooks(cfg *config.Config, req *shelf.RollbackAddBooks) error {
	args := m.Called(cfg, req)
	return args.Error(0)
}

// --- Test Cases ---

func TestNewOrderUsecaseImpl(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	assert.NotNil(t, uc)
}

func TestOrderUsecaseImpl_BuyBooks_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}} // Corrected GrpcConfig

	req := &order.BuyBooksReq{
		UserID:  1,
		BookIDs: []uint{101, 102},
	}

	// Mock FindOrdersWithBookDetail (no existing conflicting orders)
	mockRepo.On("FindOrdersWithBookDetail",
		&order.Orders{UserID: req.UserID},
		[]string{order.Completed, order.Pendding},
	).Return([]order.Orders{}, nil).Once() // Return []order.Orders

	// Mock FindBookInIds
	mockBooksResponse := &bookPb.FindBooksInIdsRes{
		Book: []*bookPb.Book{
			{Id: 101, Title: "Book 101", Price: 100},
			{Id: 102, Title: "Book 102", Price: 150},
		},
	}
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, &bookPb.FindBooksInIdsReq{Ids: []uint64{101, 102}}).Return(mockBooksResponse, nil).Once()

	// Mock CreateOrder
	mockRepo.On("CreateOrder", mock.MatchedBy(func(o *order.Orders) bool {
		return o.UserID == req.UserID && o.Status == order.Pendding && o.Total == 250 && len(o.OrdersBooks) == 2
	})).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*order.Orders)
		arg.ID = 99 // Simulate DB assigning an ID
	}).Return(nil).Once()

	// Mock DecreaseUserMoney
	mockRepo.On("DecreaseUserMoney", cfg, mock.MatchedBy(func(br *user.BuyBookReq) bool {
		return br.OrderID == 99 && br.UserID == req.UserID && br.Total == 250
	})).Return(nil).Once()

	_, err := uc.BuyBooks(cfg, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_BuyBooks_BookAlreadyInOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}

	req := &order.BuyBooksReq{
		UserID:  1,
		BookIDs: []uint{101},
	}

	// Mock FindOrdersWithBookDetail (existing order with the book)
	existingOrders := []order.Orders{ // Return []order.Orders
		{
			UserID: req.UserID,
			Status: order.Pendding,
			OrdersBooks: []order.OrdersBooks{
				{BookID: 101, Price: 100},
			},
		},
	}
	mockRepo.On("FindOrdersWithBookDetail",
		&order.Orders{UserID: req.UserID},
		[]string{order.Completed, order.Pendding},
	).Return(existingOrders, nil).Once()

	_, err := uc.BuyBooks(cfg, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book id 101 already in another order")
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_BuyBooks_FindBookInIdsError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}} // Corrected GrpcConfig

	req := &order.BuyBooksReq{
		UserID:  1,
		BookIDs: []uint{101},
	}

	mockRepo.On("FindOrdersWithBookDetail", mock.Anything, mock.Anything).Return([]order.Orders{}, nil).Once() // Return []order.Orders
	expectedErr := errors.New("find book error")
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, mock.Anything).Return(nil, expectedErr).Once()

	_, err := uc.BuyBooks(cfg, req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_BuyBooks_CreateOrderError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}} // Corrected GrpcConfig
	req := &order.BuyBooksReq{UserID: 1, BookIDs: []uint{101}}

	mockRepo.On("FindOrdersWithBookDetail", mock.Anything, mock.Anything).Return([]order.Orders{}, nil).Once() // Return []order.Orders
	mockBooksResponse := &bookPb.FindBooksInIdsRes{Book: []*bookPb.Book{{Id: 101, Price: 100}}}
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, mock.Anything).Return(mockBooksResponse, nil).Once()

	expectedErr := errors.New("create order error")
	mockRepo.On("CreateOrder", mock.Anything).Return(expectedErr).Once()

	_, err := uc.BuyBooks(cfg, req)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_BuyBooks_DecreaseMoneyError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}} // Corrected GrpcConfig
	req := &order.BuyBooksReq{UserID: 1, BookIDs: []uint{101}}

	mockRepo.On("FindOrdersWithBookDetail", mock.Anything, mock.Anything).Return([]order.Orders{}, nil).Once() // Return []order.Orders
	mockBooksResponse := &bookPb.FindBooksInIdsRes{Book: []*bookPb.Book{{Id: 101, Price: 100}}}
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, mock.Anything).Return(mockBooksResponse, nil).Once()

	mockRepo.On("CreateOrder", mock.Anything).Run(func(args mock.Arguments) {
		args.Get(0).(*order.Orders).ID = 99
	}).Return(nil).Once()

	decreaseErr := errors.New("decrease money error")
	mockRepo.On("DecreaseUserMoney", cfg, mock.Anything).Return(decreaseErr).Once()
	// Expect UpdateOrderByID to be called to set status to Failed
	mockRepo.On("UpdateOrderByID", uint(99), mock.MatchedBy(func(o *order.Orders) bool {
		return o.Status == order.Failed && o.Note == decreaseErr.Error()
	})).Return(nil).Once()


	_, err := uc.BuyBooks(cfg, req)
	assert.Error(t, err)
	assert.Equal(t, decreaseErr, err)
	mockRepo.AssertExpectations(t)
}


func TestOrderUsecaseImpl_HandleBuyBooksRes_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	res := &user.BuyBookRes{
		Error:         "", // No error from user service
		OrderID:       1,
		UserID:        10,
		TransactionID: 123, // Corrected to uint
		BookIDs:       []uint{101, 102},
	}

	// Expect AddBookToShelf to be called
	mockRepo.On("AddBookToShelf", cfg, &shelf.AddBooksReq{
		OrderID:       res.OrderID,
		UserID:        res.UserID,
		TransactionID: res.TransactionID,
		BookIDs:       res.BookIDs,
	}).Return(nil).Once()

	uc.HandleBuyBooksRes(cfg, res)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_HandleBuyBooksRes_UserServiceError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	res := &user.BuyBookRes{
		Error:   "user service error",
		OrderID: 1,
	}

	// Expect UpdateOrderByID to be called to set status to Failed
	mockRepo.On("UpdateOrderByID", res.OrderID, mock.MatchedBy(func(o *order.Orders) bool {
		return o.Status == order.Failed && o.Note == res.Error
	})).Return(nil).Once()

	uc.HandleBuyBooksRes(cfg, res)
	mockRepo.AssertExpectations(t)
	// No call to AddBookToShelf or RollbackUserTransaction should occur
}

func TestOrderUsecaseImpl_HandleBuyBooksRes_AddBookToShelfError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	res := &user.BuyBookRes{
		Error:         "", // No error from user service
		OrderID:       1,
		UserID:        10,
		TransactionID: 123, // Corrected to uint
		BookIDs:       []uint{101},
	}

	addShelfErr := errors.New("add to shelf error")
	mockRepo.On("AddBookToShelf", cfg, mock.Anything).Return(addShelfErr).Once()

	// Expect RollbackUserTransaction to be called
	mockRepo.On("RollbackUserTransaction", cfg, &user.RollbackUserTransactionReq{
		TransactionID: res.TransactionID,
	}).Return(nil).Once()

	uc.HandleBuyBooksRes(cfg, res)
	mockRepo.AssertExpectations(t)
}


func TestOrderUsecaseImpl_HandleAddBookRes_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	res := &shelf.AddBooksRes{
		Error:   "", // No error from shelf service
		OrderID: 1,
		// ShelfIDs not asserted here as UpdateOrderByID is the primary check for success path
	}

	// Expect UpdateOrderByID to be called to set status to Completed
	mockRepo.On("UpdateOrderByID", res.OrderID, &order.Orders{Status: order.Completed}).Return(nil).Once()

	uc.HandleAddBookRes(cfg, res)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_HandleAddBookRes_ShelfServiceError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	res := &shelf.AddBooksRes{
		Error:         "shelf service error",
		OrderID:       1,
		TransactionID: 123, // Corrected to uint
	}

	// Expect UpdateOrderByID to be called to set status to Failed
	mockRepo.On("UpdateOrderByID", res.OrderID, &order.Orders{Status: order.Failed, Note: res.Error}).Return(nil).Once()
	// Expect RollbackUserTransaction to be called
	mockRepo.On("RollbackUserTransaction", cfg, &user.RollbackUserTransactionReq{
		TransactionID: res.TransactionID,
	}).Return(nil).Once()


	uc.HandleAddBookRes(cfg, res)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_HandleAddBookRes_UpdateOrderErrorAfterSuccess(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	res := &shelf.AddBooksRes{
		Error:         "", // No error from shelf service
		OrderID:       1,
		TransactionID: 123, // Corrected to uint
		ShelfIDs:      []uint{501, 502},
	}

	updateErr := errors.New("update order error")
	mockRepo.On("UpdateOrderByID", res.OrderID, &order.Orders{Status: order.Completed}).Return(updateErr).Once()

	// Expect RollbackUserTransaction and RollbackAddBooks
	mockRepo.On("RollbackUserTransaction", cfg, &user.RollbackUserTransactionReq{
		TransactionID: res.TransactionID,
	}).Return(nil).Once()
	mockRepo.On("RollbackAddBooks", cfg, &shelf.RollbackAddBooks{
		ShelfIDs: res.ShelfIDs,
	}).Return(nil).Once()

	uc.HandleAddBookRes(cfg, res)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_SearchOneUserOrderByBookID_Found(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	req := &order.SearchOneMyOrderReq{UserID: 1, BookID: 101}

	ordersFromRepo := []order.Orders{ // Return []order.Orders
		{
			Model:  gorm.Model{ID: 99},
			UserID: 1,
			Status: order.Completed,
			Total:  150,
			OrdersBooks: []order.OrdersBooks{
				{BookID: 100, Price: 50},
				{BookID: 101, Price: 100}, // Target book
			},
		},
		{ // Another order, should not interfere if first one matches
			Model:  gorm.Model{ID: 98},
			UserID: 1,
			Status: order.Pendding,
			Total:  200,
			OrdersBooks: []order.OrdersBooks{
				{BookID: 102, Price: 200},
			},
		},
	}

	mockRepo.On("FindOrdersWithBookDetail",
		&order.Orders{UserID: req.UserID},
		[]string{order.Completed, order.Pendding},
	).Return(ordersFromRepo, nil).Once()

	res, err := uc.SearchOneUserOrderByBookID(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(99), res.OrderID)
	assert.Equal(t, req.UserID, res.UserID)
	assert.Equal(t, order.Completed, res.Status)
	assert.Len(t, res.Books, 2)
	assert.Equal(t, uint(101), res.Books[1].BookID) // Check if book 101 is there
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_SearchOneUserOrderByBookID_NotFound(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	req := &order.SearchOneMyOrderReq{UserID: 1, BookID: 999} // Book not in orders

	ordersFromRepo := []order.Orders{ // Return []order.Orders
		{
			UserID: 1, Status: order.Completed, OrdersBooks: []order.OrdersBooks{{BookID: 100}},
		},
	}
	mockRepo.On("FindOrdersWithBookDetail", mock.Anything, mock.Anything).Return(ordersFromRepo, nil).Once()

	res, err := uc.SearchOneUserOrderByBookID(req)
	assert.NoError(t, err)
	assert.Nil(t, res) // Expect nil result when not found
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_SearchOneUserOrderByBookID_Error(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	req := &order.SearchOneMyOrderReq{UserID: 1, BookID: 101}

	expectedErr := errors.New("db error")
	mockRepo.On("FindOrdersWithBookDetail", mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

	res, err := uc.SearchOneUserOrderByBookID(req)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

// Removed TestOrderUsecaseImpl_CanAddBooksToCart as CanAddBooksToCart is not in the OrderUsecase interface

func TestOrderUsecaseImpl_GetMyOrders_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}} // Corrected GrpcConfig
	userID := uint(1)
	now := time.Now()

	ordersFromRepo := []order.Orders{ // Return []order.Orders
		{
			Model:  gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
			UserID: userID,
			Total:  250,
			Status: order.Completed,
			OrdersBooks: []order.OrdersBooks{
				{BookID: 101, Price: 100},
				{BookID: 102, Price: 150},
			},
		},
		{
			Model:  gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
			UserID: userID,
			Total:  50,
			Status: order.Pendding,
			OrdersBooks: []order.OrdersBooks{
				{BookID: 103, Price: 50},
				{BookID: 101, Price: 0}, // Book 101 appears again, but its info should be fetched once
			},
		},
	}
	mockRepo.On("GetOrders", []uint{}, []uint{userID}, []string{}, true).Return(ordersFromRepo, nil).Once()

	mockBooksResponse := &bookPb.FindBooksInIdsRes{
		Book: []*bookPb.Book{
			{Id: 101, Title: "Book 101 Title"},
			{Id: 102, Title: "Book 102 Title"},
			{Id: 103, Title: "Book 103 Title"},
		},
	}
	// Check that bookIDs are unique and correctly gathered
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, mock.MatchedBy(func(req *bookPb.FindBooksInIdsReq) bool {
		// Simple check for length and presence of expected IDs. Order doesn't matter for the request.
		if len(req.Ids) != 3 { return false }
		found101, found102, found103 := false, false, false
		for _, id := range req.Ids {
			if id == 101 { found101 = true }
			if id == 102 { found102 = true }
			if id == 103 { found103 = true }
		}
		return found101 && found102 && found103
	})).Return(mockBooksResponse, nil).Once()

	res, err := uc.GetMyOrders(cfg, userID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Orders, 2)

	// Check order 1
	assert.Equal(t, uint(1), res.Orders[0].OrderID)
	assert.Len(t, res.Orders[0].Books, 2)
	assert.Equal(t, uint(101), res.Orders[0].Books[0].BookID)
	assert.Equal(t, "Book 101 Title", res.Orders[0].Books[0].Title)
	assert.Equal(t, uint(100), res.Orders[0].Books[0].Price) // Price from OrdersBooks

	// Check order 2
	assert.Equal(t, uint(2), res.Orders[1].OrderID)
	assert.Len(t, res.Orders[1].Books, 2) // Book 101 and 103
	// The order of books within res.Orders[1].Books might vary.
	// Find Book 101 in Order 2
	var book101InOrder2 *order.GetMyOrdersResBookDatum
	for _, b := range res.Orders[1].Books {
		if b.BookID == 101 {
			book101InOrder2 = b
			break
		}
	}
	assert.NotNil(t, book101InOrder2)
	assert.Equal(t, "Book 101 Title", book101InOrder2.Title)
	assert.Equal(t, uint(0), book101InOrder2.Price) // Price from OrdersBooks for this specific order item

	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_GetMyOrders_GetOrdersError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	userID := uint(1)
	expectedErr := errors.New("get orders error")

	mockRepo.On("GetOrders", []uint{}, []uint{userID}, []string{}, true).Return(nil, expectedErr).Once()

	res, err := uc.GetMyOrders(cfg, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestOrderUsecaseImpl_GetMyOrders_FindBookInIdsError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	uc := NewOrderUsecaseImpl(mockRepo)
	cfg := &config.Config{Grpc: config.Grpc{BookUrl: "dummy_book_url"}} // Corrected GrpcConfig
	userID := uint(1)

	ordersFromRepo := []order.Orders{ // Return []order.Orders
		{UserID: userID, OrdersBooks: []order.OrdersBooks{{BookID: 101}}},
	}
	mockRepo.On("GetOrders", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ordersFromRepo, nil).Once()

	expectedErr := errors.New("find book error")
	mockRepo.On("FindBookInIds", cfg.Grpc.BookUrl, mock.Anything).Return(nil, expectedErr).Once()

	res, err := uc.GetMyOrders(cfg, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

// Helper to create a basic config for tests that need it for gRPC URLs etc.
// Currently, only BookUrl is used from Grpc config by the usecase.
func newTestConfig() *config.Config {
	return &config.Config{
		Grpc: config.Grpc{ // Corrected GrpcConfig
			BookUrl:  "dummy_book_url",
			UserUrl:  "dummy_user_url",
			ShelfUrl: "dummy_shelf_url",
		},
	}
}
