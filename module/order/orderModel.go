package order

import "time"

const (
	Completed = "completed"
	Pendding  = "pendding"
	Failed    = "failed"
)

type (
	BuyBooksReq struct {
		UserID  uint
		BookIDs []uint `json:"books" validate:"required,min=1,dive,required"`
	}

	BookReqDatum struct {
		BookID uint `json:"book_id" validate:"required"`
	}

	BuyBooksRes struct {
		OrderID    uint            `json:"order_id"`
		UserID     uint            `json:"user_id"`
		Books      []*BookResDatum `json:"books" validate:"required,min=1,dive,required"`
		TotalPrice uint            `json:"total_price"`
		Status     string          `json:"status"`
	}

	BookResDatum struct {
		BookID uint `json:"book_id"`
		Price  uint `json:"price"`
	}

	SearchOneMyOrderReq struct {
		UserID uint
		BookID uint
	}

	SearchOneMyOrderRes struct {
		OrderID    uint            `json:"order_id"`
		UserID     uint            `json:"user_id"`
		Books      []*BookResDatum `json:"books"`
		TotalPrice uint            `json:"total_price"`
		Status     string          `json:"status"`
	}

	GetMyOrdersRes struct {
		Orders []GetMyOrdersResOrder `json:"orders"`
	}

	GetMyOrdersResOrder struct {
		OrderID    uint                       `json:"order_id"`
		UserID     uint                       `json:"user_id"`
		Books      []*GetMyOrdersResBookDatum `json:"books"`
		TotalPrice uint                       `json:"total_price"`
		Note       string                     `json:"note"`
		Status     string                     `json:"status"`
		CreatedAt  time.Time                  `json:"created_at"`
		UpdatedAt  time.Time                  `json:"updated_at"`
	}
	GetMyOrdersResBookDatum struct {
		BookID uint   `json:"book_id"`
		Title  string `json:"title"`
		Price  uint   `json:"price"`
	}
)
