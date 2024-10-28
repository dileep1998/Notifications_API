package dataservice

import (
	"Notifications_API/model"

	"gorm.io/gorm"
)

func CheckUserExists(db *gorm.DB, userID int) (bool, error) {
	var user model.User
	result := db.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func CreateTopic(db *gorm.DB, topic *model.Topic) error {
	return db.Create(topic).Error
}

func RemoveTopic(db *gorm.DB, user_id int, topic_name string) error {
	if err := db.Where("user_id = ? AND topic_name = ?", user_id, topic_name).Delete(&model.Topic{}).Error; err != nil {
		return err
	}
	return nil
}

func GetUserTopics(db *gorm.DB, userID int) ([]model.Topic, error) {
	var topics []model.Topic
	if err := db.Where("user_id = ?", userID).Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func GetSubscribersForTopic(db *gorm.DB, topic string) ([]model.User, error) {
	var users []model.User

	err := db.Joins("JOIN topics ON topics.user_id = users.user_id").
		Where("topics.topic_name = ?", topic).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
