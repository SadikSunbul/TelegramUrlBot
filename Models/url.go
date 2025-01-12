package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Url struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserTelegramId int64              `json:"userTelegramId" bson:"userTelegramId"`
	OriginalUrl    string             `json:"originalUrl" bson:"originalUrl"`
	ShortUrl       string             `json:"ShortUrl" bson:"shortUrl"`
	EndDate        primitive.DateTime `json:"endDate,omitempty" bson:"endDate,omitempty"` // null olabilir
}

type UrlInfo struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UrlId     primitive.ObjectID `json:"urlId" bson:"urlId"`
	ClickTime primitive.DateTime `json:"clickTime,omitempty" bson:"clickTime,omitempty"`
	Country   string             `json:"country,omitempty" bson:"country,omitempty"`
}
