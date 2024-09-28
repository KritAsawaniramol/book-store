package book

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
		ID   uint   `json:"id" `
		Name string `json:"name"`
	}
)
