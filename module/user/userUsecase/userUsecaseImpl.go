package userUsecase

import (
	"errors"
	"log"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userUsecaseImpl struct {
	userRepository userRepository.UserRepository
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
		log.Println("case #1")
		res.Error = err.Error()
		u.userRepository.BuyBookRes(res, cfg)
		return
	}

	if userBalance.Balance < int64(req.Total) {
		log.Printf("case #2 %v < %v", userBalance.Balance, int64(req.Total))
		res.Error = "error: not enough coin"
		u.userRepository.BuyBookRes(res, cfg)
		return
	}

	userTransaction := &user.UserTransactions{
		UserID: req.UserID,
		Amount: -int64(req.Total),
	}
	if err := u.userRepository.CreateUserTransaction(userTransaction); err != nil {
		log.Println("case #3")
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
func (u *userUsecaseImpl) CreateUserTransaction(req *user.CreateUserTransactionReq) (*user.CreateUserTransactionRes, error) {
	transactionEntity := &user.UserTransactions{
		UserID: req.UserID,
		Amount: req.Amount,
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

// FindOneUserByID implements UserUsecase.
func (u *userUsecaseImpl) FindOneUserByID(userID uint) (*userPb.UserProfile, error) {
	condition := &user.User{}
	condition.ID = userID
	user, err := u.userRepository.GetOneUser(condition)
	if err != nil {
		return nil, err
	}
	return &userPb.UserProfile{
		Id:        uint64(user.ID),
		Username:  user.Username,
		RoleId:    uint32(user.RoleID),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
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
