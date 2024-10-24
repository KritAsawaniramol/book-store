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
	"github.com/stripe/stripe-go/v80"
)

func (g *ginServer) userService() {
	stripe.Key = g.cfg.Stripe.SecretKey
	repo := userRepository.NewUserRepositoryImpl(g.db)
	usecase := userUsecase.NewUserUsecaseImpl(repo)
	httpHandler := userHandler.NewUserHttpHandler(usecase, g.cfg)
	grpcHandler := userHandler.NewUserGrpcHandler(usecase)
	queueConn, err := queue.ConnectConsumer([]string{g.cfg.Kafka.Url}, g.cfg.Kafka.ApiKey, g.cfg.Kafka.Secret, g.cfg.Kafka.GroupID)
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

	go func() {
		grpcServer, listener := grpccon.NewGrpcServer(g.cfg.Grpc.UserUrl)
		userPb.RegisterUserGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("gRPC server listening on %s\n", g.cfg.Grpc.UserUrl)
		grpcServer.Serve(listener)
	}()

	adminOnly := map[uint]bool{1: true}

	user := g.app.Group("/user_v1")
	user.GET("", g.healthCheck)
	user.POST("/user/register", httpHandler.Register)
	user.POST("/user/top-up",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(map[uint]bool{2: true}),
		httpHandler.TopUp,
	)
	user.GET("/user/top-up/:id",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(map[uint]bool{2: true}),
		httpHandler.GetOneTopUpOrder,
	)

	user.POST("/user/transaction",
		g.middleware.JwtAuthorization(),
		g.middleware.RbacAuthorization(adminOnly),
		httpHandler.AddUserTransaction,
	)

	user.GET("/user/balance", g.middleware.JwtAuthorization(),g.middleware.RbacAuthorization(map[uint]bool{2: true}), httpHandler.GetUserBalance)

	user.POST("/webhook", httpHandler.StripeWebhook)

	user.GET("/user/transaction",
		g.middleware.JwtAuthorization(), g.middleware.RbacAuthorization(adminOnly), httpHandler.SearchUserTransaction)

	user.GET("/user/profile", g.middleware.JwtAuthorization(), httpHandler.GetUserProfile)
}
