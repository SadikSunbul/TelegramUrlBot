package Database

import (
	"context"

	"github.com/SadikSunbul/TelegramUrlBot/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ICollection, koleksiyon işlemleri için interface
type ICollection interface {
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}

// IDatabase, veritabanı işlemleri için interface
type IDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
	RunCommand(ctx context.Context, runCommand interface{}) *mongo.SingleResult
}

type DataBase struct {
	Client IDatabase
}

// DatabaseWrapper, gerçek MongoDB veritabanını sarmalar
type DatabaseWrapper struct {
	*mongo.Database
}

func (d *DatabaseWrapper) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return d.Database.Collection(name, opts...)
}

func (d *DatabaseWrapper) RunCommand(ctx context.Context, runCommand interface{}) *mongo.SingleResult {
	return d.Database.RunCommand(ctx, runCommand)
}

func ConnectionDatabase() *DataBase {
	config := config.GetConfig()
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(config.MongoDbConnect))
	if err != nil {
		panic(err)
	}
	return &DataBase{Client: &DatabaseWrapper{client.Database(config.DbName)}}
}
