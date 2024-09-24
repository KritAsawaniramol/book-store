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
)
