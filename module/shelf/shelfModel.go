package shelf

type (
	AddBooksReq struct {
		OrderID       uint   `json:"order_id" validate:"required"`
		UserID        uint   `json:"user_id"  validate:"required"`
		TransactionID uint   `json:"transaction_id" validate:"required"`
		BookIDs       []uint `json:"books" validate:"required,min=1"`
	}

	AddBooksRes struct {
		OrderID       uint   `json:"order_id" validate:"required"`
		UserID        uint   `json:"user_id"  validate:"required"`
		TransactionID uint   `json:"transaction_id" validate:"required"`
		ShelfIDs      []uint `json:"shelf_ids"`
		Error         string  `json:"error"`
	}

	RollbackAddBooks struct {
		ShelfIDs []uint `json:"shelf_ids" validate:"required,min=1"`
	}
)
