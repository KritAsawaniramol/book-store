package queue

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
)

func ConnectProducer(brokerUrls []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	// config.Producer.RequiredAcks
	// This setting determines how many acknowledgments the producer needs to
	// receive from Kafka brokers before considering a message as successfully sent.
	// It impacts the reliability of message delivery.

	// sarama.NoResponse (0):
	// The producer doesn't wait for any acknowledgment from the broker.
	// This is the fastest option but the least reliable. If the message
	// is lost before being written to a broker, the producer won't know about it.

	// sarama.WaitForLocal (1):
	// The producer waits for acknowledgment from the leader broker only. This is
	// faster than WaitForAll, but less reliable because if the leader broker crashes
	// before replicating the message, the message could be lost.

	// sarama.WaitForAll (-1):
	// The producer waits for acknowledgments from all in-sync replicas.
	// This is the most reliable option since the message is replicated to
	// multiple brokers before being considered successful.
	config.Producer.RequiredAcks = sarama.WaitForAll


	//max number of retry when failed to connect to Producer
	config.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(brokerUrls, config)
	if err != nil {
		log.Printf("error: ConnectProducer: Failed to connect to producer: %s", err.Error())
		return nil, errors.New("error: failed to connect to producer")
	}
	return producer, nil
}

func ConnectConsumer(brokerUrls []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Fetch.Max = 3

	consumer, err := sarama.NewConsumerGroup(brokerUrls, groupID, config)
	if err != nil {
		log.Printf("error: ConnectConsumer: Failed to connect to consumer: %s", err.Error())
		return nil, errors.New("error: failed to connect to consumer")
	}
	return consumer, nil
}

func DecodeMessage(obj any, value []byte) error {
	if err := json.Unmarshal(value, obj); err != nil {
		log.Printf("error: DecodeMessage: %s\n", err.Error())
		return errors.New("error: decode message failed")
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		log.Printf("error: DecodeMessage: %s\n", err.Error())
		return errors.New("error: validate message failed")
	}
	return nil
}

func PushMessageWithKeyToQueue(brokerUrls []string, topic, key string, message []byte) error {
	producer, err := ConnectProducer(brokerUrls)
	if err != nil {
		log.Printf("error: PushMessageWithKeyToQueue: %s\n", err.Error())
		return errors.New("error: connect to producer failed")
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
		Key:   sarama.StringEncoder(key),
	}

	log.Printf("msg: %v\n", msg)
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("error: PushMessageWithKeyToQueue: %s\n", err.Error())
		return errors.New("error: send message failed")
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
