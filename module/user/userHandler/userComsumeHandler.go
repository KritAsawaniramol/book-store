package userHandler

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
)

type userConsumerHandler struct {
	userUsecase userUsecase.UserUsecase
	cfg         *config.Config
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (u userConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session finished")
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (u userConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Topic(%s)| Key(%s) | Offset(%d) Message(%s) \n", msg.Topic, string(msg.Key), msg.Offset, string(msg.Value))
		switch string(msg.Key) {
		case "buy":
			req := &user.BuyBookReq{}
			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				log.Printf("error: ConsumeClaim: %s\n", err.Error())
				continue
			}
			u.userUsecase.BuyBook(u.cfg, req)
		case "rollbacktransaction":
			req := &user.RollbackUserTransactionReq{}
			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				continue
			}
			u.userUsecase.RollbackUserTransaction(req)
		default:
			log.Println("no consumer handler")
		}
		// MarkMessage marks a message as consumed.
		session.MarkMessage(msg, "")
	}
	return nil
}

// Setup implements sarama.ConsumerGroupHandler.
func (u userConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session started")
	return nil
}

func NewUserConsumerHandler(userUsecase userUsecase.UserUsecase, cfg *config.Config) sarama.ConsumerGroupHandler {
	return userConsumerHandler{userUsecase: userUsecase, cfg: cfg}
}
