package api

import (
	"Notifications_API/dataservice"
	"Notifications_API/model"

	"gorm.io/gorm"
)

func CheckUserExists(db *gorm.DB, userID int) (bool, error) {
	return dataservice.CheckUserExists(db, userID)
}

func CreateTopic(db *gorm.DB, topic *model.Topic) error {
	return dataservice.CreateTopic(db, topic)
}

func RemoveTopic(db *gorm.DB, user_id int, topic_name string) error {
	return dataservice.RemoveTopic(db, user_id, topic_name)
}

func GetUserTopics(db *gorm.DB, user_id int) ([]model.Topic, error) {
	return dataservice.GetUserTopics(db, user_id)
}
