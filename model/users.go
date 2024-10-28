package model

type User struct {
	UserID      int    `gorm:"column:user_id;primaryKey"`
	Email       string `gorm:"column:email"`
	PhoneNumber string `gorm:"column:phone_number"`
}
