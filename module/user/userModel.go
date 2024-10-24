package user

import "time"

type (
	UserRegisterReq struct {
		Username string `json:"username" validate:"required,max=64,min=1"`
		Password string `json:"password" validate:"required,max=64,min=6"`
	}

	UserProfile struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		RoleID   uint   `json:"role_id"`
		Coin     int64  `json:"coin"`
	}

	CreateUserTransactionReq struct {
		UserID uint  `json:"user_id"`
		Amount int64 `json:"amount"  validate:"required"`
	}

	BuyBookReq struct {
		OrderID uint   `json:"order_id" validate:"required"`
		UserID  uint   `json:"user_id"  validate:"required"`
		Total   uint   `json:"total" validate:"required"`
		BookIDs []uint `json:"books" validate:"required,min=1"`
	}

	RollbackUserTransactionReq struct {
		TransactionID uint `json:"transaction_id"`
	}

	BuyBookRes struct {
		OrderID       uint   `json:"order_id" validate:"required"`
		UserID        uint   `json:"user_id"  validate:"required"`
		Total         uint   `json:"total" validate:"required"`
		BookIDs       []uint `json:"books" validate:"required,min=1"`
		TransactionID uint   `json:"transaction_id"`
		Error         string `json:"error"`
	}

	CreateUserTransactionRes struct {
		TransactionID uint  `json:"transaction_id"`
		Balance       int64 `json:"balance"`
		// Error         string `json:"error"`
	}

	UserBalanceRes struct {
		Balance int64 `json:"balance"`
	}

	SearchUserTransactionReq struct {
		UsersID        uint
		TransactionsID uint
	}

	SearchUserTransactionRes struct {
		Transactions []UserTransactionsDatum `json:"transactions"`
	}

	UserTransactionsDatum struct {
		CreatedAt     time.Time `json:"created_at"`
		TransactionID uint      `json:"transaction_id"`
		UserID        uint      `json:"user_id"`
		Username      string    `json:"username"`
		Amount        int64     `json:"amount"`
		UpdatedAt     time.Time `json:"updated_at"`
		Note          string    `json:"note"`
	}

	TopUpReq struct {
		UserID uint  `json:"user_id"`
		Amount int64 `json:"amount"  validate:"required,min=10"`
	}

	GetOneTopUpOrderRes struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		Amount    int64  `json:"amount"`
		SessionID string `json:"session_id"`
		Status    string `json:"status"`
	}
)
