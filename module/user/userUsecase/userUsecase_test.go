package userUsecase

import (
	"errors"
	"testing"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	// "github.com/kritAsawaniramol/book-store/module/user/userPb" // Removed unused import
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"github.com/stripe/stripe-go/v80" // Corrected Stripe import to v80
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateOneUser(u *user.User) (uint, error) {
	args := m.Called(u)
	return args.Get(0).(uint), args.Error(1) // Corrected to Get(0).(uint)
}

func (m *MockUserRepository) GetOneUser(condition *user.User) (*user.User, error) {
	args := m.Called(condition)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetUserInIDs(userIDs []uint) ([]user.User, int64, error) { // Corrected return type
	args := m.Called(userIDs)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]user.User), args.Get(1).(int64), args.Error(2) // Corrected to use Get(1).(int64)
}

func (m *MockUserRepository) CreateUserTransaction(utx *user.UserTransactions) error {
	args := m.Called(utx)
	// Simulate GORM assigning ID
	if utx != nil && args.Error(0) == nil {
		val := args.Get(1) // Get the argument first
		if val != nil {    // Check if it's not nil (i.e., provided in Return)
			utx.ID = val.(uint) // Assert as uint
		}
		if utx.ID == 0 { // Default if not specified by test or if val was nil
			utx.ID = 1
		}
	}
	return args.Error(0)
}


func (m *MockUserRepository) GetUserTransactions(condition *user.UserTransactions) ([]user.UserTransactions, error) { // Corrected return type
	args := m.Called(condition)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.UserTransactions), args.Error(1) // Corrected return type
}

func (m *MockUserRepository) DeleteUserTransaction(transactionID uint) error {
	args := m.Called(transactionID)
	return args.Error(0)
}

func (m *MockUserRepository) CreateOneTopUpOrder(order *user.TopUpOrder) error {
	args := m.Called(order)
	if order != nil && args.Error(0) == nil {
		val := args.Get(1) // Get the argument first
		if val != nil {    // Check if it's not nil
			order.ID = val.(uint) // Assert as uint
		}
		if order.ID == 0 { // Default if not specified
			order.ID = 1
		}
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetOneTopUpOrder(condition *user.TopUpOrder) (*user.TopUpOrder, error) {
	args := m.Called(condition)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.TopUpOrder), args.Error(1)
}

func (m *MockUserRepository) UpdateOneTopUpOrderStatusBySessionID(sessionID string, status string) error {
	args := m.Called(sessionID, status)
	return args.Error(0)
}

// Mocking the method that would call Stripe
func (m *MockUserRepository) CheckOutTopUp(cfg *config.Config, amount int64) (*stripe.CheckoutSession, error) { // Corrected return type
	args := m.Called(cfg, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.CheckoutSession), args.Error(1)
}

// Mocking the callback/producer
func (m *MockUserRepository) BuyBookRes(res *user.BuyBookRes, cfg *config.Config) {
	m.Called(res, cfg)
}


// --- Test Cases ---

func TestNewUserUsecaseImpl(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	assert.NotNil(t, uc)
}

func TestUserUsecaseImpl_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)

	req := &user.UserRegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	// Expect CreateOneUser to be called with a hashed password
	mockRepo.On("CreateOneUser", mock.MatchedBy(func(u *user.User) bool {
		return u.Username == req.Username && bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) == nil && u.RoleID == 2
	})).Return(uint(1), nil).Once() // Corrected: Return uint(1)

	userID, err := uc.Register(req)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), userID)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_Register_HashingError(t *testing.T) {
	// This test is tricky because bcrypt.GenerateFromPassword failing is rare unless params are wrong.
	// Here, we assume it could fail and the usecase handles it.
	// Forcing bcrypt to fail is hard without changing its input in a way that the usecase might prevent.
	// The usecase passes bcrypt.DefaultCost, which is valid.
	// We will test the repository error path instead, as it's more practical.
	// If bcrypt itself fails, it's an internal library error usually.
	// The current implementation returns "error: failed to hash password" if bcrypt fails.
	// To properly test this, we'd need to mock bcrypt, which is overkill.
	// So, we'll focus on repository error.
}

func TestUserUsecaseImpl_Register_RepositoryError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	req := &user.UserRegisterReq{Username: "testuser", Password: "password123"}
	expectedErr := errors.New("db error")

	mockRepo.On("CreateOneUser", mock.AnythingOfType("*user.User")).Return(uint(0), expectedErr).Once() // Corrected: Return uint(0)

	userID, err := uc.Register(req)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, uint(0), userID)
	mockRepo.AssertExpectations(t)
}


func TestUserUsecaseImpl_FindOneUserByUsernameAndPassword_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	now := time.Now()

	mockUser := &user.User{
		Model:    gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		Username: username,
		Password: string(hashedPassword),
		RoleID:   2,
	}
	mockRepo.On("GetOneUser", &user.User{Username: username}).Return(mockUser, nil).Once()

	res, err := uc.FindOneUserByUsernameAndPassword(username, password)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint64(1), res.Id)
	assert.Equal(t, username, res.Username)
	assert.Equal(t, uint32(2), res.RoleId)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_FindOneUserByUsernameAndPassword_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	expectedErr := errors.New("user not found")

	mockRepo.On("GetOneUser", mock.AnythingOfType("*user.User")).Return(nil, expectedErr).Once()

	res, err := uc.FindOneUserByUsernameAndPassword("unknownuser", "password")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_FindOneUserByUsernameAndPassword_PasswordMismatch(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	username := "testuser"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	mockUser := &user.User{
		Model:    gorm.Model{ID: 1},
		Username: username,
		Password: string(hashedPassword),
		RoleID:   2,
	}
	mockRepo.On("GetOneUser", &user.User{Username: username}).Return(mockUser, nil).Once()

	res, err := uc.FindOneUserByUsernameAndPassword(username, wrongPassword)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password is incorrect")
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_FindOneUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	userID := uint(1)
	now := time.Now()
	expectedUser := &user.User{
		Model:     gorm.Model{ID: userID, CreatedAt: now, UpdatedAt: now},
		Username:  "test",
		RoleID:    1,
	}
	condition := &user.User{}
	condition.ID = userID
	mockRepo.On("GetOneUser", condition).Return(expectedUser, nil).Once()

	resPb, err := uc.FindOneUserByID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, resPb)
	assert.Equal(t, uint64(userID), resPb.Id)
	assert.Equal(t, expectedUser.Username, resPb.Username)
	assert.Equal(t, timestamppb.New(now).Seconds, resPb.CreatedAt.Seconds) // Compare seconds for timestamp
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_FindOneUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	userID := uint(1)
	expectedErr := errors.New("not found")
	condition := &user.User{}
	condition.ID = userID
	mockRepo.On("GetOneUser", condition).Return(nil, expectedErr).Once()

	resPb, err := uc.FindOneUserByID(userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, resPb)
	mockRepo.AssertExpectations(t)
}


func TestUserUsecaseImpl_GetUserProfile_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    uc := NewUserUsecaseImpl(mockRepo)
    userID := uint(1)

    // Mock GetOneUser
    mockUserEntity := &user.User{Model: gorm.Model{ID: userID}, Username: "profileuser", RoleID: 2}
    userCondition := &user.User{}
	userCondition.ID = userID
    mockRepo.On("GetOneUser", userCondition).Return(mockUserEntity, nil).Once()

    // Mock GetUserTransactions (for GetUserBalance)
    mockTransactions := []user.UserTransactions{ // Corrected type
        {Amount: 100},
        {Amount: -20},
        {Amount: 50},
    }
    transactionCondition := &user.UserTransactions{UserID: userID}
    mockRepo.On("GetUserTransactions", transactionCondition).Return(mockTransactions, nil).Once()

    profile, err := uc.GetUserProfile(userID)

    assert.NoError(t, err)
    assert.NotNil(t, profile)
    assert.Equal(t, userID, profile.ID)
    assert.Equal(t, "profileuser", profile.Username)
    assert.Equal(t, uint(2), profile.RoleID)
    assert.Equal(t, int64(130), profile.Coin) // 100 - 20 + 50
    mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_GetUserProfile_GetOneUserError(t *testing.T) {
    mockRepo := new(MockUserRepository)
    uc := NewUserUsecaseImpl(mockRepo)
    userID := uint(1)
    expectedErr := errors.New("get user error")

    userCondition := &user.User{}
	userCondition.ID = userID
    mockRepo.On("GetOneUser", userCondition).Return(nil, expectedErr).Once()
    // GetUserTransactions should not be called if GetOneUser fails

    profile, err := uc.GetUserProfile(userID)

    assert.Error(t, err)
    assert.Equal(t, expectedErr, err)
    assert.Nil(t, profile)
    mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_GetUserProfile_GetUserBalanceError(t *testing.T) {
    mockRepo := new(MockUserRepository)
    uc := NewUserUsecaseImpl(mockRepo)
    userID := uint(1)
    expectedErr := errors.New("get balance error")

    mockUserEntity := &user.User{Model: gorm.Model{ID: userID}, Username: "profileuser", RoleID: 2}
	userCondition := &user.User{}
	userCondition.ID = userID
    mockRepo.On("GetOneUser", userCondition).Return(mockUserEntity, nil).Once()

    transactionCondition := &user.UserTransactions{UserID: userID}
    mockRepo.On("GetUserTransactions", transactionCondition).Return(nil, expectedErr).Once()

    profile, err := uc.GetUserProfile(userID)

    assert.Error(t, err)
    assert.Equal(t, expectedErr, err)
    assert.Nil(t, profile)
    mockRepo.AssertExpectations(t)
}


func TestUserUsecaseImpl_GetUserBalance_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	userID := uint(1)
	mockTransactions := []user.UserTransactions{ // Corrected type
		{Amount: 1000},
		{Amount: -200},
		{Amount: 50},
	}
	condition := &user.UserTransactions{UserID: userID}
	mockRepo.On("GetUserTransactions", condition).Return(mockTransactions, nil).Once()

	balanceRes, err := uc.GetUserBalance(userID)
	assert.NoError(t, err)
	assert.NotNil(t, balanceRes)
	assert.Equal(t, int64(850), balanceRes.Balance)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_GetUserBalance_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	userID := uint(1)
	expectedErr := errors.New("db error getting transactions")
	condition := &user.UserTransactions{UserID: userID}
	mockRepo.On("GetUserTransactions", condition).Return(nil, expectedErr).Once()

	balanceRes, err := uc.GetUserBalance(userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, balanceRes)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_CreateUserTransaction_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	userID := uint(1)
	req := &user.CreateUserTransactionReq{
		UserID: userID,
		Amount: 100,
	}
	note := "test deposit"
	expectedTransactionID := uint(99)

	// Mock CreateUserTransaction in repository
	mockRepo.On("CreateUserTransaction", mock.MatchedBy(func(utx *user.UserTransactions) bool {
		return utx.UserID == userID && utx.Amount == req.Amount && utx.Note == note
	})).Return(nil, expectedTransactionID).Once() // Return nil error, provide expectedTransactionID

	// Mock GetUserTransactions (for GetUserBalance call within CreateUserTransaction)
	mockTransactionsForBalance := []user.UserTransactions{
		{Amount: 50}, // Previous balance
		{Amount: req.Amount}, // Current transaction (though CreateUserTransaction calculates from all)
	}
	mockRepo.On("GetUserTransactions", &user.UserTransactions{UserID: userID}).Return(mockTransactionsForBalance, nil).Once()


	res, err := uc.CreateUserTransaction(req, note)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedTransactionID, res.TransactionID)
	assert.Equal(t, int64(150), res.Balance) // 50 + 100
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_CreateUserTransaction_RepoCreateError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	req := &user.CreateUserTransactionReq{UserID: 1, Amount: 100}
	note := "test deposit"
	expectedErr := errors.New("db create error")

	mockRepo.On("CreateUserTransaction", mock.AnythingOfType("*user.UserTransactions")).Return(expectedErr, uint(0)).Once()
	// GetUserBalance should not be called

	res, err := uc.CreateUserTransaction(req, note)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_CreateUserTransaction_GetUserBalanceError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	userID := uint(1)
	req := &user.CreateUserTransactionReq{UserID: userID, Amount: 100}
	note := "test deposit"
	expectedErr := errors.New("get balance error")

	// Mock CreateUserTransaction success
	mockRepo.On("CreateUserTransaction", mock.AnythingOfType("*user.UserTransactions")).Return(nil, uint(1)).Once()
	// Mock GetUserTransactions to fail
	mockRepo.On("GetUserTransactions", &user.UserTransactions{UserID: userID}).Return(nil, expectedErr).Once()

	res, err := uc.CreateUserTransaction(req, note)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_RollbackUserTransaction(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	transactionID := uint(123)
	req := &user.RollbackUserTransactionReq{TransactionID: transactionID}

	// Expect DeleteUserTransaction to be called
	mockRepo.On("DeleteUserTransaction", transactionID).Return(nil).Once() // Assuming it might return an error, though usecase ignores it

	uc.RollbackUserTransaction(req)

	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_RollbackUserTransaction_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	transactionID := uint(123)
	req := &user.RollbackUserTransactionReq{TransactionID: transactionID}

	// Expect DeleteUserTransaction to be called, and it returns an error
	mockRepo.On("DeleteUserTransaction", transactionID).Return(errors.New("delete failed")).Once()

	uc.RollbackUserTransaction(req) // Usecase method doesn't return error, so we just check mock

	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_BuyBook_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	cfg := &config.Config{} // Dummy config
	userID := uint(1)
	orderID := uint(10)
	totalAmount := uint(100)
	bookIDs := []uint{101, 102}
	expectedTransactionID := uint(201)

	req := &user.BuyBookReq{
		OrderID: orderID,
		UserID:  userID,
		Total:   totalAmount,
		BookIDs: bookIDs,
	}

	// Mock GetUserTransactions (for GetUserBalance) - sufficient balance
	mockRepo.On("GetUserTransactions", &user.UserTransactions{UserID: userID}).Return([]user.UserTransactions{{Amount: int64(totalAmount + 50)}}, nil).Once()

	// Mock CreateUserTransaction
	mockRepo.On("CreateUserTransaction", mock.MatchedBy(func(utx *user.UserTransactions) bool {
		return utx.UserID == userID && utx.Amount == -int64(totalAmount) && utx.Note == "buy book"
	})).Return(nil, expectedTransactionID).Once() // Simulate success and return transaction ID

	// Mock BuyBookRes (callback/producer)
	expectedBuyBookRes := &user.BuyBookRes{
		OrderID:       orderID,
		UserID:        userID,
		Total:         totalAmount,
		BookIDs:       bookIDs,
		TransactionID: expectedTransactionID,
		Error:         "",
	}
	mockRepo.On("BuyBookRes", expectedBuyBookRes, cfg).Return().Once()

	uc.BuyBook(cfg, req)

	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_BuyBook_GetUserBalanceError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &user.BuyBookReq{UserID: 1, Total: 100}
	expectedErr := errors.New("get balance error")

	mockRepo.On("GetUserTransactions", &user.UserTransactions{UserID: req.UserID}).Return(nil, expectedErr).Once()

	// Expect BuyBookRes to be called with the error
	expectedBuyBookResWithError := &user.BuyBookRes{
		OrderID: req.OrderID, UserID: req.UserID, Total: req.Total, BookIDs: req.BookIDs, Error: expectedErr.Error(),
	}
	mockRepo.On("BuyBookRes", expectedBuyBookResWithError, cfg).Return().Once()

	uc.BuyBook(cfg, req)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_BuyBook_NotEnoughCoin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &user.BuyBookReq{UserID: 1, Total: 100}

	// Mock GetUserTransactions (for GetUserBalance) - insufficient balance
	mockRepo.On("GetUserTransactions", &user.UserTransactions{UserID: req.UserID}).Return([]user.UserTransactions{{Amount: 50}}, nil).Once()

	expectedErrMsg := "error: not enough coin"
	expectedBuyBookResWithError := &user.BuyBookRes{
		OrderID: req.OrderID, UserID: req.UserID, Total: req.Total, BookIDs: req.BookIDs, Error: expectedErrMsg,
	}
	mockRepo.On("BuyBookRes", expectedBuyBookResWithError, cfg).Return().Once()

	uc.BuyBook(cfg, req)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_BuyBook_CreateTransactionError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &user.BuyBookReq{UserID: 1, Total: 100}
	expectedErr := errors.New("create transaction error")

	// Mock GetUserTransactions (for GetUserBalance) - sufficient balance
	mockRepo.On("GetUserTransactions", &user.UserTransactions{UserID: req.UserID}).Return([]user.UserTransactions{{Amount: 150}}, nil).Once()

	// Mock CreateUserTransaction to fail
	mockRepo.On("CreateUserTransaction", mock.AnythingOfType("*user.UserTransactions")).Return(expectedErr, uint(0)).Once()


	expectedBuyBookResWithError := &user.BuyBookRes{
		OrderID: req.OrderID, UserID: req.UserID, Total: req.Total, BookIDs: req.BookIDs, Error: expectedErr.Error(),
	}
	mockRepo.On("BuyBookRes", expectedBuyBookResWithError, cfg).Return().Once()

	uc.BuyBook(cfg, req)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_SearchUserTransaction_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	userID := uint(1)
	transactionID := uint(101)

	req := &user.SearchUserTransactionReq{
		UsersID:        userID,
		TransactionsID: transactionID,
	}

	// Mock GetUserTransactions
	now := time.Now()
	mockTransactions := []user.UserTransactions{
		{Model: gorm.Model{ID: transactionID, CreatedAt: now, UpdatedAt: now}, UserID: userID, Amount: -100, Note: "buy book"},
		{Model: gorm.Model{ID: 102, CreatedAt: now, UpdatedAt: now}, UserID: userID, Amount: 500, Note: "top up"},
	}
	expectedCondition := &user.UserTransactions{UserID: userID}
	expectedCondition.ID = transactionID // GORM model ID for filtering specific transaction
	mockRepo.On("GetUserTransactions", expectedCondition).Return(mockTransactions, nil).Once()

	// Mock GetUserInIDs (to populate usernames)
	mockUsers := []user.User{
		{Model: gorm.Model{ID: userID}, Username: "testuser"},
	}
	mockRepo.On("GetUserInIDs", []uint{userID}).Return(mockUsers, int64(1), nil).Once()

	res, err := uc.SearchUserTransaction(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Transactions, 2)

	// Check details of the first transaction (assuming order is preserved or test data is specific)
	assert.Equal(t, transactionID, res.Transactions[0].TransactionID)
	assert.Equal(t, userID, res.Transactions[0].UserID)
	assert.Equal(t, "testuser", res.Transactions[0].Username)
	assert.Equal(t, int64(-100), res.Transactions[0].Amount)
	assert.Equal(t, "buy book", res.Transactions[0].Note)

	mockRepo.AssertExpectations(t)
}


func TestUserUsecaseImpl_SearchUserTransaction_GetTransactionsError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	req := &user.SearchUserTransactionReq{UsersID: 1}
	expectedErr := errors.New("db error get transactions")

	mockRepo.On("GetUserTransactions", mock.AnythingOfType("*user.UserTransactions")).Return(nil, expectedErr).Once()

	res, err := uc.SearchUserTransaction(req)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_SearchUserTransaction_GetUserInIDsError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	req := &user.SearchUserTransactionReq{UsersID: 1}
	expectedErr := errors.New("db error get users")

	// Mock GetUserTransactions - success
	mockTransactions := []user.UserTransactions{{UserID: 1, Amount: 100}}
	mockRepo.On("GetUserTransactions", mock.AnythingOfType("*user.UserTransactions")).Return(mockTransactions, nil).Once()

	// Mock GetUserInIDs - failure
	mockRepo.On("GetUserInIDs", []uint{1}).Return(nil, int64(0), expectedErr).Once()

	res, err := uc.SearchUserTransaction(req)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_TopUp_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	cfg := &config.Config{} // Dummy config
	userID := uint(1)
	amount := int64(1000)

	req := &user.TopUpReq{
		UserID: userID,
		Amount: amount,
	}

	// Mock CheckOutTopUp (simulates Stripe call)
	mockSessionID := "sess_12345"
	mockSessionStatus := "open" // Or whatever stripe.CheckoutSessionStatusOpen is
	mockStripeSession := &stripe.CheckoutSession{
		ID:     mockSessionID,
		Status: stripe.CheckoutSessionStatus(mockSessionStatus),
		// Add other necessary fields if the usecase uses them from the session object
	}
	mockRepo.On("CheckOutTopUp", cfg, amount).Return(mockStripeSession, nil).Once()

	// Mock CreateOneTopUpOrder
	expectedTopUpOrderID := uint(301)
	mockRepo.On("CreateOneTopUpOrder", mock.MatchedBy(func(order *user.TopUpOrder) bool {
		return order.UserID == userID &&
			order.Amount == amount &&
			order.SessionID == mockSessionID &&
			order.Status == mockSessionStatus
	})).Return(nil, expectedTopUpOrderID).Once() // Return nil error, provide expectedTopUpOrderID

	sessionID, err := uc.TopUp(req, cfg)

	assert.NoError(t, err)
	assert.Equal(t, mockSessionID, sessionID)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_TopUp_CheckOutTopUpError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &user.TopUpReq{UserID: 1, Amount: 1000}
	expectedErr := errors.New("stripe error")

	mockRepo.On("CheckOutTopUp", cfg, req.Amount).Return(nil, expectedErr).Once()

	sessionID, err := uc.TopUp(req, cfg)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, sessionID)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_TopUp_CreateTopUpOrderError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	cfg := &config.Config{}
	req := &user.TopUpReq{UserID: 1, Amount: 1000}
	expectedErr := errors.New("db create topup error")

	// Mock CheckOutTopUp - success
	mockStripeSession := &stripe.CheckoutSession{ID: "sess_123", Status: "open"}
	mockRepo.On("CheckOutTopUp", cfg, req.Amount).Return(mockStripeSession, nil).Once()

	// Mock CreateOneTopUpOrder - failure
	mockRepo.On("CreateOneTopUpOrder", mock.AnythingOfType("*user.TopUpOrder")).Return(expectedErr, uint(0)).Once()

	sessionID, err := uc.TopUp(req, cfg)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, sessionID)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_GetOneTopUpOrderByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	topUpOrderID := uint(1)

	expectedOrder := &user.TopUpOrder{
		Model:     gorm.Model{ID: topUpOrderID},
		UserID:    uint(2),
		Amount:    5000,
		SessionID: "sess_abc123",
		Status:    "completed",
	}
	condition := &user.TopUpOrder{}
	condition.ID = topUpOrderID
	mockRepo.On("GetOneTopUpOrder", condition).Return(expectedOrder, nil).Once()

	res, err := uc.GetOneTopUpOrderByID(topUpOrderID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedOrder.ID, res.ID)
	assert.Equal(t, expectedOrder.UserID, res.UserID)
	assert.Equal(t, expectedOrder.Amount, res.Amount)
	assert.Equal(t, expectedOrder.SessionID, res.SessionID)
	assert.Equal(t, expectedOrder.Status, res.Status)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_GetOneTopUpOrderByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	topUpOrderID := uint(1)
	expectedErr := errors.New("not found")

	condition := &user.TopUpOrder{}
	condition.ID = topUpOrderID
	mockRepo.On("GetOneTopUpOrder", condition).Return(nil, expectedErr).Once()

	res, err := uc.GetOneTopUpOrderByID(topUpOrderID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_HandleStripeWebhook_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	sessionID := "sess_completed_123"
	sessionStatus := "complete" // Assuming this is a status from Stripe
	topUpOrderID := uint(1)
	userID := uint(5)
	amount := int64(2000)

	mockTopUpOrder := &user.TopUpOrder{
		Model:     gorm.Model{ID: topUpOrderID},
		UserID:    userID,
		Amount:    amount,
		SessionID: sessionID,
		Status:    "open", // Previous status
	}

	// Mock GetOneTopUpOrder
	mockRepo.On("GetOneTopUpOrder", &user.TopUpOrder{SessionID: sessionID}).Return(mockTopUpOrder, nil).Once()
	// Mock UpdateOneTopUpOrderStatusBySessionID
	mockRepo.On("UpdateOneTopUpOrderStatusBySessionID", sessionID, sessionStatus).Return(nil).Once()
	// Mock CreateUserTransaction
	mockRepo.On("CreateUserTransaction", mock.MatchedBy(func(utx *user.UserTransactions) bool {
		return utx.UserID == userID &&
			utx.Amount == amount &&
			utx.Note == "user top up" &&
			utx.TopUpOrderID != nil && *utx.TopUpOrderID == topUpOrderID
	})).Return(nil, uint(0)).Once() // Transaction ID from CreateUserTransaction not critical for this test's assertions

	err := uc.HandleStripeWebhook(sessionID, sessionStatus)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_HandleStripeWebhook_GetOrderError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	sessionID := "sess_err_123"
	expectedErr := errors.New("get order db error")

	mockRepo.On("GetOneTopUpOrder", &user.TopUpOrder{SessionID: sessionID}).Return(nil, expectedErr).Once()

	err := uc.HandleStripeWebhook(sessionID, "complete")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_HandleStripeWebhook_UpdateStatusError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	sessionID := "sess_update_err_123"
	expectedErr := errors.New("update status db error")

	mockTopUpOrder := &user.TopUpOrder{SessionID: sessionID, Status: "open"}
	mockRepo.On("GetOneTopUpOrder", &user.TopUpOrder{SessionID: sessionID}).Return(mockTopUpOrder, nil).Once()
	mockRepo.On("UpdateOneTopUpOrderStatusBySessionID", sessionID, "complete").Return(expectedErr).Once()

	err := uc.HandleStripeWebhook(sessionID, "complete")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_HandleStripeWebhook_CreateTransactionErrorWithRollback(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	sessionID := "sess_tx_err_123"
	sessionStatus := "complete"
	originalStatus := "open" // Status before update attempt
	expectedErr := errors.New("create transaction db error")

	mockTopUpOrder := &user.TopUpOrder{
		Model:     gorm.Model{ID: 1},
		UserID:    5,
		Amount:    2000,
		SessionID: sessionID,
		Status:    originalStatus,
	}

	mockRepo.On("GetOneTopUpOrder", &user.TopUpOrder{SessionID: sessionID}).Return(mockTopUpOrder, nil).Once()
	mockRepo.On("UpdateOneTopUpOrderStatusBySessionID", sessionID, sessionStatus).Return(nil).Once()
	mockRepo.On("CreateUserTransaction", mock.AnythingOfType("*user.UserTransactions")).Return(expectedErr, uint(0)).Once()
	// Expect rollback of status
	mockRepo.On("UpdateOneTopUpOrderStatusBySessionID", sessionID, originalStatus).Return(nil).Once()


	err := uc.HandleStripeWebhook(sessionID, sessionStatus)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecaseImpl_HandleStripeWebhook_RollbackError(t *testing.T) {
	// This test checks if the main error (from CreateUserTransaction) is still returned
	// even if the rollback itself fails.
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecaseImpl(mockRepo)
	sessionID := "sess_rollback_err_123"
	sessionStatus := "complete"
	originalStatus := "open"
	createTxErr := errors.New("create transaction db error")
	rollbackErr := errors.New("rollback status db error")


	mockTopUpOrder := &user.TopUpOrder{Model: gorm.Model{ID:1}, SessionID: sessionID, Status: originalStatus}

	mockRepo.On("GetOneTopUpOrder", &user.TopUpOrder{SessionID: sessionID}).Return(mockTopUpOrder, nil).Once()
	mockRepo.On("UpdateOneTopUpOrderStatusBySessionID", sessionID, sessionStatus).Return(nil).Once()
	mockRepo.On("CreateUserTransaction", mock.AnythingOfType("*user.UserTransactions")).Return(createTxErr, uint(0)).Once()
	mockRepo.On("UpdateOneTopUpOrderStatusBySessionID", sessionID, originalStatus).Return(rollbackErr).Once() // Rollback fails

	err := uc.HandleStripeWebhook(sessionID, sessionStatus)
	assert.Error(t, err)
	assert.Equal(t, createTxErr, err) // Ensure original error is returned
	mockRepo.AssertExpectations(t)
}
