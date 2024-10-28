package model

type Topic struct {
	UserID            int    `gorm:"column:user_id"`
	TopicName         string `gorm:"column:topic_name"`
	ChannelEmail      string `gorm:"column:channel_email"`
	ChannelPhone      string `gorm:"column:channel_phone"`
	ChannelPushNotify bool   `gorm:"column:channel_push_notify"`
}
