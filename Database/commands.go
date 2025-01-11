package Database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DataBase) Add(col string, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := db.Client.Collection(col).InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DataBase) Update(col string, id string, data interface{}) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := db.Client.Collection(col).UpdateOne(context.TODO(), bson.D{{"_id", objectID}}, bson.D{{"$set", data}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DataBase) Delete(col string, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = db.Client.Collection(col).DeleteOne(context.TODO(), bson.D{{"_id", objectID}})
	if err != nil {
		return err
	}
	return nil
}
