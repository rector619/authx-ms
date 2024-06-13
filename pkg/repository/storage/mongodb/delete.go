package mongodb

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *Database) SoftDelete(model interface{}) error {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": getDocumentID(model)}
	update := bson.M{"$set": bson.M{"deleted": true, "deleted_at": time.Now()}}

	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) SoftDeleteByFilter(model interface{}, filter map[string]interface{}) error {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"deleted": true, "deleted_at": time.Now()}}
	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) HardDelete(model interface{}) error {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": getDocumentID(model)}
	_, err = coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) HardDeleteByFilter(model interface{}, filter map[string]interface{}) error {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return err
	}

	_, err = coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) HardDeleteManyByFilter(model interface{}, filter map[string]interface{}) error {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return err
	}

	_, err = coll.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func getDocumentID(model interface{}) primitive.ObjectID {
	var (
		modelValue = reflect.ValueOf(model)
		modelType  = reflect.TypeOf(model)
		value      reflect.Value
	)

	if modelType.Kind() == reflect.Ptr && modelType.Elem().Kind() == reflect.Pointer {
		value = modelValue.Elem().Elem().FieldByName("ID")
	} else if modelType.Kind() == reflect.Ptr && modelType.Elem().Kind() == reflect.Struct {
		value = modelValue.Elem().FieldByName("ID")
	} else {
		value = modelValue.Elem().FieldByName("ID")
	}

	if !value.IsValid() {
		return primitive.NilObjectID
	}
	return value.Interface().(primitive.ObjectID)
}
