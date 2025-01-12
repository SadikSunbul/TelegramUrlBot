package Database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DataBase) Get(col, id string) (*mongo.SingleResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := db.Client.Collection(col).FindOne(context.TODO(), bson.D{{"_id", objectID}})

	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}

func (db *DataBase) GetBy(col string, data interface{}) (*mongo.SingleResult, error) {

	result := db.Client.Collection(col).FindOne(context.TODO(), data)

	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}

func (db *DataBase) GetList(col string, data interface{}) (*mongo.Cursor, error) {

	result, err := db.Client.Collection(col).Find(context.TODO(), data)
	if err != nil {
		return nil, err
	}

	return result, nil
}
