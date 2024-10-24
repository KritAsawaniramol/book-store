package orderHandler

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/order/orderUsecase"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
)

type orderConsumerHandler struct {
	orderUsecase orderUsecase.OrderUsecase
	cfg          *config.Config
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (o *orderConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Topic(%s)| Key(%s) | Offset(%d) Message(%s) \n", msg.Topic, string(msg.Key), msg.Offset, string(msg.Value))
		switch string(msg.Key) {
		case "usertransaction":
			res := &user.BuyBookRes{}
			if err := queue.DecodeMessage(res, msg.Value); err != nil {
				log.Printf("error: ConsumeClaim: %s\n", err.Error())
				continue
			}
			o.orderUsecase.HandleBuyBooksRes(o.cfg, res)
		case "addbook":
			res := &shelf.AddBooksRes{}
			if err := queue.DecodeMessage(res, msg.Value); err != nil {
				log.Printf("error: ConsumeClaim: %s\n", err.Error())
				continue
			}
			o.orderUsecase.HandleAddBookRes(o.cfg, res)
		default:
			log.Println("no consumer handler")
		}
		// MarkMessage marks a message as consumed.
		session.MarkMessage(msg, "")
	}
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (o *orderConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// Setup implements sarama.ConsumerGroupHandler.
func (o *orderConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func NewOrderConsumerHandler(orderUsecase orderUsecase.OrderUsecase, cfg *config.Config) sarama.ConsumerGroupHandler {
	return &orderConsumerHandler{orderUsecase: orderUsecase, cfg: cfg}
}
