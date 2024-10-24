package auth

import "time"

type (
	LoginReq struct {
		Username string `json:"username" validate:"required,max=64"`
		Password string `json:"password" validate:"required,max=64"`
	}

	LoginRes struct {
		ID         uint           `json:"id"`
		Username   string         `json:"username"`
		RoleID     uint           `json:"role_id"`
		Coin       int64          `json:"coin"`
		CreatedAt  time.Time      `json:"created_at"`
		UpdatedAt  time.Time      `json:"updated_at"`
		Credential *CredentialRes `json:"credential"`
	}

	CredentialRes struct {
		ID           uint      `json:"id"`
		UserID       uint      `json:"user_id"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	LogoutReq struct {
		CredentialID uint `json:"credential_id" validate:"required"`
	}

	RefreshTokenReq struct {
		CredentialID uint   `json:"credential_id" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"required,max=500"`
	}
)
