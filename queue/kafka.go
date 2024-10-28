package queue

import (
	"Notifications_API/dataservice"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

func SendMessage(topic string, message string, producer sarama.SyncProducer) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := producer.SendMessage(msg)
	return err
}

func ConsumeMessages(db *gorm.DB, partitionConsumer sarama.PartitionConsumer) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	fmt.Println("Listening for messages...")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			topic := msg.Topic
			message := string(msg.Value)
			fmt.Printf("Received message on topic '%s': %s\n", topic, message)

			users, err := dataservice.GetSubscribersForTopic(db, topic)
			if err != nil {
				log.Printf("Failed to fetch subscribers for topic %s: %v", topic, err)
				continue
			}

			for _, user := range users {
				fmt.Printf("Sending message to %s (%s): %s\n", user.Email, message)
			}

		case <-signals:
			fmt.Println("Interrupt received, shutting down...")
			return
		}
	}
}
