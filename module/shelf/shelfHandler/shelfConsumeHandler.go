package shelfHandler

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/shelf/shelfUsecase"
	"github.com/kritAsawaniramol/book-store/pkg/queue"
)

type shelfConsumerHandler struct {
	shelfUsecase shelfUsecase.ShelfUsecase
	cfg          *config.Config
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (s shelfConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Topic(%s)| Key(%s) | Offset(%d) Message(%s) \n", msg.Topic, string(msg.Key), msg.Offset, string(msg.Value))
		switch string(msg.Key) {
		case "addbook":
			req := &shelf.AddBooksReq{}
			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				continue
			}
			s.shelfUsecase.AddBooks(s.cfg, req)
		case "rollbackaddbook":
		default:
			log.Println("no consumer handler")
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (s shelfConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// Setup implements sarama.ConsumerGroupHandler.
func (s shelfConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func NewShelfConsumerHandler(shelfUsecase shelfUsecase.ShelfUsecase, cfg *config.Config) sarama.ConsumerGroupHandler {
	return shelfConsumerHandler{shelfUsecase: shelfUsecase, cfg: cfg}
}
