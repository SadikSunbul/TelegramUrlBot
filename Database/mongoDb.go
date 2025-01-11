package Database

import (
	"context"
	"github.com/SadikSunbul/TelegramUrlBot/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBase struct {
	Client *mongo.Database
}

func ConnectionDatabase() *DataBase {
	config := config.GetConfig()
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(config.MongoDbConnect))
	if err != nil {
		panic(err)
	}
	return &DataBase{Client: client.Database(config.DbName)}
}
