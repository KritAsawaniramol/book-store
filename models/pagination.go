package models

type (
	PaginatieReq struct {
		Page  *uint `form:"page,omitempty"`
		Limit *int  `form:"limit,omitempty" validate:"omitempty,lte=25"`
	}

	PaginatieRes struct {
		Limit           int    `json:"limit"`
		LastVisiblePage int64  `json:"last_visible_page"`
		HasNextPage     bool   `json:"has_next_page"`
		Total           int64 `json:"total"`
	}
)
