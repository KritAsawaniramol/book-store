package authUsecase

import "github.com/kritAsawaniramol/book-store/module/auth/authRepository"

type authUsecaseImpl struct {
	authRepository authRepository.AuthRepository
}

func NewAuthUsecaseImpl(authRepository authRepository.AuthRepository) AuthUsecase  {
	return &authUsecaseImpl{authRepository: authRepository}
}