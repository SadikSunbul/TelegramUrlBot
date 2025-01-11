package Database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *DataBase) Add(col string, data interface{}) (interface{}, error) {
	result, err := db.Client.Collection(col).InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (db *DataBase) Update(col string, id string, data interface{}) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	result, err := db.Client.Collection(col).UpdateOne(context.TODO(), bson.D{{"id", objectID}}, bson.D{{"$set", data}})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (db *DataBase) Delete(col string, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = db.Client.Collection(col).DeleteOne(context.TODO(), bson.D{{"id", objectID}})
	if err != nil {
		return err
	}
	return nil
}
