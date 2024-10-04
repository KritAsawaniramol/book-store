package server

import (
	"context"
	"log"

	"github.com/kritAsawaniramol/book-store/module/user/userHandler"
	userPb "github.com/kritAsawaniramol/book-store/module/user/userPb"
	"github.com/kritAsawaniramol/book-store/module/user/userRepository"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/grpccon"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
)

func (g *ginServer) userService() {
	repo := userRepository.NewUserRepositoryImpl(g.db)
	usecase := userUsecase.NewUserUsecaseImpl(repo)
	httpHandler := userHandler.NewUserHttpHandler(usecase)
	grpcHandler := userHandler.NewUserGrpcHandler(usecase)
	queueConn, err := queue.ConnectConsumer([]string{g.cfg.Kafka.Url}, g.cfg.Kafka.GroupID)
	if err != nil {
		log.Fatal(err)
	}

	consumerHandler := userHandler.NewUserConsumerHandler(usecase, g.cfg)
	go func() {
		for {
			if err := queueConn.Consume(context.Background(), []string{"user"}, consumerHandler); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// go func() {
	// 	err := <-queueConn.Errors()
	// 	log.Printf("error: queueConn: %s\n", err.Error())
	// }()

	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(g.cfg.Grpc.UserUrl)
		userPb.RegisterUserGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("gRPC server listening on %s\n", g.cfg.Grpc.UserUrl)
		grpcServer.Serve(listener)
	}()

	g.app.POST("/register", httpHandler.Register)
	g.app.POST("/user/top-up",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(
			map[uint]bool{
				1: true,
			},
		),
		httpHandler.AddUserMoney)
	// g.app.GET("/user/coin/:id", )
	// g.app.GET("/user/", )
}
