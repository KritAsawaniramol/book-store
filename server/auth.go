package server

import (
	"log"

	"github.com/kritAsawaniramol/book-store/module/auth/authHandler"
	"github.com/kritAsawaniramol/book-store/module/auth/authPb"
	"github.com/kritAsawaniramol/book-store/module/auth/authRepository"
	"github.com/kritAsawaniramol/book-store/module/auth/authUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
)

func (g *ginServer) authService() {
	userServiceGrpcClient, err := grpccon.NewGrpcClient(g.cfg.Grpc.UserUrl)
	if err != nil {
		log.Fatal(err)
	}
	repo := authRepository.NewAuthRepositoryImpl(g.db, userServiceGrpcClient.User())
	usecase := authUsecase.NewAuthUsecaseImpl(repo)
	httpHandler := authHandler.NewAuthHttpHandlerImpl(g.cfg, usecase)
	grpcHandler := authHandler.NewAuthGrpcHandlerImpl(usecase)

	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(g.cfg.Grpc.AuthUrl)
		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("gRPC server listening on %s\n", g.cfg.Grpc.AuthUrl)
		grpcServer.Serve(listener)
	}()

	auth := g.app.Group("/auth_v1")

	auth.GET("", g.healthCheck)
	auth.POST("/auth/login", httpHandler.Login)
	auth.POST("/auth/logout", httpHandler.Logout)
	auth.POST("/auth/refresh-token", httpHandler.RefreshToken)
}
