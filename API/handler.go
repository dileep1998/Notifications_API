package api

import (
	"Notifications_API/dataservice"
	"Notifications_API/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TopicHandler struct {
	db *gorm.DB
}

func NewTopicHandler(db *gorm.DB) *TopicHandler {
	return &TopicHandler{db: db}
}

type ChannelRequest struct {
	Email             string `json:"email" binding:"omitempty,email"`
	Sms               string `json:"sms" binding:"omitempty,sms"`
	PushNotifications bool   `json:"push_notifications"`
}

type SubscriptionRequest struct {
	Topic    string         `json:"topic" binding:"required"`
	Channels ChannelRequest `json:"channels" binding:"required"`
}

type UserSubscriptionRequest struct {
	UserID        int                   `json:"user_id" binding:"required"`
	Subscriptions []SubscriptionRequest `json:"subscriptions" binding:"required"`
}

type UserUnSubscribeRequest struct {
	UserID int      `json:"user_id" gorm:"primaryKey"`
	Topics []string `json:"topics" gorm:"type:text[]"`
}

func (h *TopicHandler) unsubscribeTopic(c *gin.Context) {
	var request UserUnSubscribeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userExists, err := CheckUserExists(h.db, request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	for _, unsubscription := range request.Topics {
		user_id := request.UserID
		topic_name := unsubscription

		if err := RemoveTopic(h.db, user_id, topic_name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove topic"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully removed topics"})

}

func (h *TopicHandler) AddTopic(c *gin.Context) {
	var req UserSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userExists, err := CheckUserExists(h.db, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	for _, subscription := range req.Subscriptions {
		topic := model.Topic{
			UserID:            req.UserID,
			TopicName:         subscription.Topic,
			ChannelEmail:      subscription.Channels.Email,
			ChannelPhone:      subscription.Channels.Sms,
			ChannelPushNotify: subscription.Channels.PushNotifications,
		}

		if err := CreateTopic(h.db, &topic); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create topic"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscriptions added successfully"})
}

func (h *TopicHandler) GetUserTopics(c *gin.Context) {
	userID := c.Param("user_id")
	user_ID, err := strconv.Atoi(userID)

	topics, err := GetUserTopics(h.db, user_ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "topics": topics})
}

func (h *TopicHandler) SendNotifications(c *gin.Context) {
	var request struct {
		Topic string `json:"topic" binding:"required"`
		Event struct {
			EventID   string `json:"event_id"`
			Timestamp string `json:"timestamp"`
			Details   struct {
				UserID   string `json:"user_id"`
				Email    string `json:"email"`
				Username string `json:"username"`
			} `json:"details"`
		} `json:"event"`
		Message struct {
			Title string `json:"title" binding:"required"`
			Body  string `json:"body" binding:"required"`
		} `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	users, err := dataservice.GetSubscribersForTopic(h.db, request.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscribers"})
		return
	}

	for _, user := range users {
		h.sendNotificationToUser(user, request.Message.Title, request.Message.Body)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications sent successfully"})
}

func (h *TopicHandler) sendNotificationToUser(user model.User, title, body string) {
	println("Sending notification to:", user.Email)
	println("Title:", title)
	println("Body:", body)
}
