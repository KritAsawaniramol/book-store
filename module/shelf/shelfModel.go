package shelf

import "time"

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
		Error         string `json:"error"`
	}

	RollbackAddBooks struct {
		ShelfIDs []uint `json:"shelf_ids" validate:"required,min=1"`
	}

	IsBookInShelfReq struct {
		BookID uint
		UserID uint
	}

	IsBookInShelfRes struct {
		ID     uint `json:"id"`
		BookID uint `json:"book_id"`
		UserID uint `json:"user_id"`
	}

	GetMyShelf struct {
		Shelves []ShelfRes `json:"shelves"`
	}

	ShelfRes struct {
		ID     uint    `json:"id"`
		UserID uint    `json:"user_id"`
		Book   BookRes `json:"book"`
	}

	BookRes struct {
		ID            uint      `json:"book_id"`
		Title         string    `json:"title"`
		CoverImageUrl string    `json:"cover_image_url"`
		AuthorName    string    `json:"author_name"`
		CreatedAt     time.Time `json:"created_at"`
	}
)
