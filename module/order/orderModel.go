package order

type (
	BuyBooksReq struct {
		UserID uint
		Books  []*BookReqDatum `json:"books" validate:"required,min=1,dive,required"`
	}

	BookReqDatum struct {
		BookID uint `json:"book_id" validate:"required"`
	}

	BuyBooksRes struct {
		OrderID    uint            `json:"order_id"`
		UserID     uint            `json:"user_id"`
		Books      []*BookResDatum `json:"books" validate:"required,min=1,dive,required"`
		TotalPrice uint            `json:"total_price"`
		Status     string          `json:"stauts"`
	}

	BookResDatum struct {
		BookID uint `json:"book_id"`
		Price  uint `json:"price"`
	}
)
