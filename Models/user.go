package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TelegramId int64              `json:"telegramId" bson:"telegramId"`
	Name       string             `json:"name" bson:"name"`
}
