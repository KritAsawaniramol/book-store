package userHandler

import (
	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/user/userUsecase"
)

type UserQueueHandler interface {
}

type userQueueHandlerImpl struct {
	userUsecase userUsecase.UserUsecase
	cfg         *config.Config
}

// func (u *userQueueHandlerImpl) UserConsumer() {
// 	consumer, err := queue.ConnectConsumer([]string{u.cfg.Kafka.Url}, u.cfg.Kafka.GroupID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	consumer.Consume()
// }
