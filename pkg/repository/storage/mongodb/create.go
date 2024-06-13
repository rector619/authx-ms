package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *Database) CreateOneRecord(model interface{}) error {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return err
	}

	result, err := coll.InsertOne(context.Background(), model)
	if err != nil {
		return err
	}

	if result.InsertedID == nil {
		return fmt.Errorf("record creation for %v failed", coll.Name())
	}

	filter := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
	update := bson.M{"$set": bson.M{"updated_at": time.Now(), "created_at": time.Now(), "deleted": false}}
	re, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if re.ModifiedCount > 0 {
		err := db.SelectOneFromDb(model, bson.M{"_id": result.InsertedID.(primitive.ObjectID)})
		if err != nil {
			return err
		}
	}
	return nil
}
