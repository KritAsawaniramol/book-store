package user

type (
	UserRegisterReq struct {
		Username string `json:"username" validate:"required,max=64"`
		Password string `json:"password" validate:"required,max=64"`
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
		Error         string  `json:"error"`
	}

	CreateUserTransactionRes struct {
		TransactionID uint  `json:"transaction_id"`
		Balance       int64 `json:"balance"`
		// Error         string `json:"error"`
	}

	UserBalanceRes struct {
		TransactionID uint  `json:"transaction_id"`
		Balance       int64 `json:"balance"`
	}

	
)
