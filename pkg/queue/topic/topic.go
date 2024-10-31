package main

import (
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/kritAsawaniramol/book-store/config"
)

func main() {
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())
	// brokers := []string{"localhost:9092"}
	brokers := []string{cfg.Kafka.Url}
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	clusterAdmin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatal(err.Error())
	}

	topicDetail := &sarama.TopicDetail{
		// NumPartitions contains the number of partitions to create in the topic, or
		// -1 if we are either specifying a manual partition assignment or using the
		// default partitions.
		NumPartitions: 1,
		// ReplicationFactor contains the number of replicas to create for each
		// partition in the topic, or -1 if we are either specifying a manual
		// partition assignment or using the default replication factor.
		ReplicationFactor: 1,
	}

	topicName := []string{"order", "user", "shelf"}

	for _, name := range topicName {
		if err := clusterAdmin.CreateTopic(name, topicDetail, false); err != nil {
			log.Printf("Error creating topic: %v\n", err)
		} else {
			log.Printf("Topic %s created successfully\n", name)
		}
	}
}
