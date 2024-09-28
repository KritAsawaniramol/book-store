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
	repo := authRepository.NewAuthRepositoryImpl(g.db)
	usecase := authUsecase.NewAuthUsecaseImpl(repo)
	httpHandler := authHandler.NewAuthHttpHandlerImpl(g.cfg, usecase)
	grpcHandler := authHandler.NewAuthGrpcHandlerImpl(usecase)

	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(g.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("gRPC server listening on %s\n", g.cfg.Grpc.AuthUrl)
		grpcServer.Serve(listener)
	}()

	g.app.POST("/auth/login", httpHandler.Login)
	g.app.POST("/auth/logout", httpHandler.Logout)
	g.app.POST("/auth/refresh-token", httpHandler.RefreshToken)
}
