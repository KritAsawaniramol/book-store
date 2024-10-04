package book

import (
	"time"

	"github.com/kritAsawaniramol/book-store/models"
)

type (
	CreateBookReq struct {
		Title          string     `json:"title" validate:"required,max=200"`
		Price          uint       `json:"price" validate:"required"`
		FilePath       string     `json:"file_path"`
		CoverImagePath string     `json:"cover_image_path"`
		AuthorName     string     `json:"author_name"`
		Tags           []BookTags `json:"tags"`
	}

	BookTags struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	SearchBooksReq struct {
		Title      string `form:"title" validate:"max=200"`
		MaxPrice   *uint  `form:"max_price"`
		MinPrice   *uint  `form:"min_price"`
		AuthorName string `form:"author_name"`
		TagIDs     string `form:"tag_ids"`
		models.PaginatieReq
	}

	BookRes struct {
		Title         string     `json:"title"`
		Price         uint       `json:"price"`
		CoverImageUrl string     `json:"cover_image_Url"`
		AuthorName    string     `json:"author_name"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
		Tags          []BookTags `json:"tags"`
	}

	SearchBooksRes struct {
		Books      []BookRes           `json:"books"`
		Pagination models.PaginatieRes `json:"pagination"`
	}
)
