package server

import (
	"github.com/kritAsawaniramol/book-store/module/user/userHandler"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
)

func (g *ginServer) authService() {
	repo := userRepository.NewUserRepositoryImpl(g.db)
	usecase := userUsecase.NewUserUsecaseImpl(repo)
	httpHandler := userHandler.NewUserHttpHandler(usecase)
	queueHandler := userHandler.NewUserQueueHandler(usecase, g.cfg)
	grpcHandler := userHandler.NewUserGrpcHandler(usecase)

	_ = httpHandler 
	_ = queueHandler
	_ = grpcHandler 
}
