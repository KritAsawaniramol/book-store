package server

import (
	"github.com/kritAsawaniramol/book-store/module/auth/authHandler"
	"github.com/kritAsawaniramol/book-store/module/auth/authRepository"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
)

func (g *ginServer) authService() {
	repo := authRepository.NewAuthRepositoryImpl(g.db)
	usecase := authUsecase.NewAuthUsecaseImpl(repo)
	httpHandler := authHandler.NewAuthHttpHandlerImpl(usecase)

	_ = httpHandler
}
