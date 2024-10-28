package api

import (
	"Notifications_API/queue"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationRequest struct {
	Topic   string `json:"topic" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func RegisterRoutes(db *gorm.DB, producer sarama.SyncProducer) *gin.Engine {
	router := gin.Default()

	topicHandler := NewTopicHandler(db)
	router.POST("/subscribe_to_topic", topicHandler.AddTopic)

	router.POST("/notifications/send", topicHandler.SendNotifications)

	router.POST("/unsubscribe_to_topics", topicHandler.unsubscribeTopic)

	router.GET("/subsciptions/:user_id", topicHandler.GetUserTopics)

	router.POST("/notifications/send_to_topic", func(c *gin.Context) {
		var req NotificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := queue.SendMessage(req.Topic, req.Message, producer); err != nil {
			log.Printf("Failed to send message to Kafka: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
			return
		}

		log.Printf("Message sent to topic '%s': %s", req.Topic, req.Message)
		c.JSON(http.StatusOK, gin.H{"status": "Message sent"})

		brokers := []string{"localhost:9092"}

		consumer, partitionConsumer := initKafkaConsumer(brokers, req.Topic)
		defer consumer.Close()
		defer partitionConsumer.Close()
		queue.ConsumeMessages(db, partitionConsumer)

	})

	return router
}

func initKafkaConsumer(brokers []string, topic string) (sarama.Consumer, sarama.PartitionConsumer) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}

	return consumer, partitionConsumer
}
