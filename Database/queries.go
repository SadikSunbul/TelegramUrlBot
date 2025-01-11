package Database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *DataBase) Get(col, id string) (interface{}, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := db.Client.Collection(col).FindOne(context.TODO(), bson.D{{"id", objectID}})

	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}

func (db *DataBase) GetBy(col string, data interface{}) (interface{}, error) {

	result := db.Client.Collection(col).FindOne(context.TODO(), data)

	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}

func (db *DataBase) GetList(col string, data interface{}) (interface{}, error) {

	result, err := db.Client.Collection(col).Find(context.TODO(), bson.D{{"$set", data}})
	if err != nil {
		return nil, err
	}
	var response interface{}
	err = result.Decode(response)
	if err != nil {
		return nil, err
	}
	return result, nil
}
