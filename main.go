package main

import (
	api "Notifications_API/API"
	"log"

	"github.com/IBM/sarama"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:Linux@123@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	producer, err := initKafkaProducer()
	if err != nil {
		log.Fatal("Error while creating kafka producer")
	}
	defer producer.Close()

	router := api.RegisterRoutes(db, producer)
	router.Run(":8080")

}

func initKafkaProducer() (sarama.SyncProducer, error) {
	brokerList := []string{"localhost:9092"}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
