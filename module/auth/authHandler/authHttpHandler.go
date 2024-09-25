package authHandler

import "github.com/kritAsawaniramol/book-store/module/auth/authUsecase"

type (
	AuthHttpHandler interface{}

	authHttpHandlerImpl struct {
		authUsecase authUsecase.AuthUsecase
	}
)

func NewAuthHttpHandlerImpl(authUsecase authUsecase.AuthUsecase) AuthHttpHandler {
	return &authHttpHandlerImpl{authUsecase: authUsecase}
}
