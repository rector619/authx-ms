package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (db *Database) SaveAllFields(model interface{}) error {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": getDocumentID(model)}
	update := bson.M{"$set": model}

	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	update = bson.M{"$set": bson.M{"updated_at": time.Now()}}
	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 {
		err := db.SelectOneFromDb(model, bson.M{"_id": getDocumentID(model)})
		if err != nil {
			return err
		}
	}

	return nil
}
