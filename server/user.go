package server

import (
	"log"

	"github.com/kritAsawaniramol/book-store/module/user/userHandler"
	userPb "github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
)

func (g *ginServer) userService() {
	repo := userRepository.NewUserRepositoryImpl(g.db)
	usecase := userUsecase.NewUserUsecaseImpl(repo)
	httpHandler := userHandler.NewUserHttpHandler(usecase)
	grpcHandler := userHandler.NewUserGrpcHandler(usecase)
	queueHandler := userHandler.NewUserQueueHandler(usecase, g.cfg)

	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(g.cfg.Grpc.UserUrl)

		userPb.RegisterUserGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("gRPC server listening on %s\n", g.cfg.Grpc.UserUrl)
		grpcServer.Serve(listener)
	}()

	_ = queueHandler

	g.app.POST("/register", httpHandler.Register)
}
