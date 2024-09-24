package server

import (
	"github.com/kritAsawaniramol/book-store/module/user/userHandler"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
)

func (g *ginServer) userService() {
	repo := userRepository.NewUserRepositoryImpl(g.db)
	usecase := userUsecase.NewUserUsecaseImpl(repo)
	httpHandler := userHandler.NewUserHttpHandler(usecase)
	queueHandler := userHandler.NewUserQueueHandler(usecase, g.cfg)
	grpcHandler := userHandler.NewUserQueueHandler(usecase, g.cfg)

	_ = grpcHandler
	_ = queueHandler
	
	g.app.POST("/register", httpHandler.Register)
}
