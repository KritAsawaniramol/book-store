package userUsecase

import (
	"errors"
	"log"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"github.com/kritAsawaniramol/book-store/util"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userUsecaseImpl struct {
	userRepository userRepository.UserRepository
}

// HandleStripeWebhook implements UserUsecase.
func (u *userUsecaseImpl) HandleStripeWebhook(sessionID string, sessionStatus string) error {
	topUpOrder, err := u.userRepository.GetOneTopUpOrder(&user.TopUpOrder{SessionID: sessionID})
	if err != nil {
		return err
	}

	if err := u.userRepository.UpdateOneTopUpOrderStatusBySessionID(sessionID, sessionStatus); err != nil {
		return err
	}

	if err := u.userRepository.CreateUserTransaction(&user.UserTransactions{
		UserID:       topUpOrder.UserID,
		Amount:       topUpOrder.Amount,
		Note:         "user top up",
		TopUpOrderID: &topUpOrder.ID,
	}); err != nil {
		u.rollbackUpdateOneTopUpOrderStatus(sessionID, topUpOrder.Status)
		return err
	}
	return nil
}

func (u *userUsecaseImpl) rollbackUpdateOneTopUpOrderStatus(sessionID, oldSessionStatus string) error {
	if err := u.userRepository.UpdateOneTopUpOrderStatusBySessionID(sessionID, oldSessionStatus); err != nil {
		return err
	}
	return nil
}

// GetOneTopUpOrderByID implements UserUsecase.
func (u *userUsecaseImpl) GetOneTopUpOrderByID(id uint) (*user.GetOneTopUpOrderRes, error) {
	condition := &user.TopUpOrder{}
	condition.ID = id
	result, err := u.userRepository.GetOneTopUpOrder(condition)
	if err != nil {
		return nil, err
	}
	return &user.GetOneTopUpOrderRes{
		ID:        result.ID,
		UserID:    result.UserID,
		Amount:    result.Amount,
		SessionID: result.SessionID,
		Status:    result.Status,
	}, nil
}

// TopUp implements UserUsecase.
func (u *userUsecaseImpl) TopUp(req *user.TopUpReq, cfg *config.Config) (string, error) {

	session, err := u.userRepository.CheckOutTopUp(cfg, req.Amount)
	if err != nil {
		return "", err
	}

	if err := u.userRepository.CreateOneTopUpOrder(&user.TopUpOrder{
		UserID:    req.UserID,
		Amount:    req.Amount,
		SessionID: session.ID,
		Status:    string(session.Status),
	}); err != nil {
		return "", err
	}

	return session.ID, nil
}

// SearchUserTransaction implements UserUsecase.
func (u *userUsecaseImpl) SearchUserTransaction(req *user.SearchUserTransactionReq) (*user.SearchUserTransactionRes, error) {
	condition := &user.UserTransactions{
		UserID: req.UsersID,
	}
	condition.ID = req.TransactionsID
	result, err := u.userRepository.GetUserTransactions(condition)
	if err != nil {
		return nil, err
	}

	util.PrintObjInJson(result)

	userMapUsername := map[uint]string{}
	uniqueUserID := []uint{}

	for _, r := range result {
		if _, ok := userMapUsername[r.UserID]; !ok {
			userMapUsername[r.UserID] = ""
			uniqueUserID = append(uniqueUserID, r.UserID)
		}
	}
	userDatum, _, err := u.userRepository.GetUserInIDs(uniqueUserID)
	if err != nil {
		return nil, err
	}

	for _, data := range userDatum {
		userMapUsername[data.ID] = data.Username
	}

	transactions := []user.UserTransactionsDatum{}
	for _, r := range result {
		transactions = append(transactions, user.UserTransactionsDatum{
			CreatedAt:     r.CreatedAt,
			TransactionID: r.ID,
			UserID:        r.UserID,
			Username:      userMapUsername[r.UserID],
			Amount:        r.Amount,
			UpdatedAt:     r.UpdatedAt,
			Note:          r.Note,
		})
	}
	return &user.SearchUserTransactionRes{Transactions: transactions}, nil
}

// BuyBook implements UserUsecase.
func (u *userUsecaseImpl) BuyBook(cfg *config.Config, req *user.BuyBookReq) {
	res := &user.BuyBookRes{
		OrderID: req.OrderID,
		UserID:  req.UserID,
		Total:   req.Total,
		BookIDs: req.BookIDs,
		Error:   "",
	}

	userBalance, err := u.GetUserBalance(req.UserID)
	if err != nil {
		res.Error = err.Error()
		u.userRepository.BuyBookRes(res, cfg)
		return
	}

	if userBalance.Balance < int64(req.Total) {
		res.Error = "error: not enough coin"
		u.userRepository.BuyBookRes(res, cfg)
		return
	}

	userTransaction := &user.UserTransactions{
		UserID: req.UserID,
		Amount: -int64(req.Total),
		Note:   "buy book",
	}
	if err := u.userRepository.CreateUserTransaction(userTransaction); err != nil {
		res.Error = err.Error()
		u.userRepository.BuyBookRes(res, cfg)
		return
	}
	res.TransactionID = userTransaction.ID
	u.userRepository.BuyBookRes(res, cfg)
}

// RollbackUserTransaction implements UserUsecase.
func (u *userUsecaseImpl) RollbackUserTransaction(req *user.RollbackUserTransactionReq) {
	u.userRepository.DeleteUserTransaction(req.TransactionID)
}

// CreateUserTransaction implements UserUsecase.
func (u *userUsecaseImpl) CreateUserTransaction(req *user.CreateUserTransactionReq, note string) (*user.CreateUserTransactionRes, error) {
	transactionEntity := &user.UserTransactions{
		UserID: req.UserID,
		Amount: req.Amount,
		Note:   note,
	}
	if err := u.userRepository.CreateUserTransaction(transactionEntity); err != nil {
		return nil, err
	}
	balance, err := u.GetUserBalance(req.UserID)
	if err != nil {
		return nil, err
	}
	return &user.CreateUserTransactionRes{
		TransactionID: transactionEntity.ID,
		Balance:       balance.Balance,
	}, nil
}

func (u *userUsecaseImpl) GetUserBalance(userID uint) (*user.UserBalanceRes, error) {
	transactions, err := u.userRepository.GetUserTransactions(&user.UserTransactions{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	var coin int64 = 0
	for _, t := range transactions {
		coin += t.Amount
	}
	return &user.UserBalanceRes{Balance: coin}, nil
}

// GetUserProfile implements UserUsecase.
func (u *userUsecaseImpl) GetUserProfile(userID uint) (*user.UserProfile, error) {
	condition := &user.User{}
	condition.ID = userID
	userData, err := u.userRepository.GetOneUser(condition)
	if err != nil {
		return nil, err
	}
	balance, err := u.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &user.UserProfile{
		ID:       userData.ID,
		Username: userData.Username,
		RoleID:   userData.RoleID,
		Coin:     balance.Balance,
	}, nil
}

// FindOneUserByID implements UserUsecase.
func (u *userUsecaseImpl) FindOneUserByID(userID uint) (*userPb.UserProfile, error) {
	condition := &user.User{}
	condition.ID = userID
	result, err := u.userRepository.GetOneUser(condition)
	if err != nil {
		return nil, err
	}
	return &userPb.UserProfile{
		Id:        uint64(result.ID),
		Username:  result.Username,
		RoleId:    uint32(result.RoleID),
		CreatedAt: timestamppb.New(result.CreatedAt),
		UpdatedAt: timestamppb.New(result.UpdatedAt),
	}, nil
}

// FindOneUserByUsernameAndPassword implements UserUsecase.
func (u *userUsecaseImpl) FindOneUserByUsernameAndPassword(username string, password string) (*userPb.UserProfile, error) {
	user, err := u.userRepository.GetOneUser(&user.User{Username: username})
	if err != nil {
		return nil, err
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("error: FindOneUserByUsernameAndPassword: %s\n", err.Error())
		return nil, errors.New("error: password is incorrect")
	}
	return &userPb.UserProfile{
		Id:       uint64(user.ID),
		Username: user.Username,
		RoleId:   uint32(user.RoleID),
		// Coin:      user.Coin,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// Register implements UserUsecase.
func (u *userUsecaseImpl) Register(registReq *user.UserRegisterReq) (uint, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error: Register: %s\n", err.Error())
		return 0, errors.New("error: failed to hash password")
	}

	userID, err := u.userRepository.CreateOneUser(&user.User{
		Username: registReq.Username,
		Password: string(hashedPassword),
		RoleID:   2,
		// Coin:     0,
	})
	if err != nil {
		return 0, err
	}
	return userID, err
}

func NewUserUsecaseImpl(userRepository userRepository.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepository: userRepository,
	}
}
