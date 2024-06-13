package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MigrateModel(ctx context.Context, model interface{}) error {
	// Get the collection name from the model.
	name := CollectionName(model)

	// Get the collection from the context.
	coll := GetCollection(ctx, name)
	if coll == nil {
		return fmt.Errorf("collection '%s' not found in context", name)
	}

	// Create the indexes and options from the model.
	indexes := make([]mongo.IndexModel, 0)
	collation := &options.Collation{}
	structType := reflect.TypeOf(model).Elem()
	fields := make([]string, 0)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		bsonTag, ok := field.Tag.Lookup("bson")
		if !ok {
			jsonTag, ok := field.Tag.Lookup("json")
			if !ok {
				bsonTag = field.Name
			} else {
				bsonTag = strings.Split(jsonTag, ",")[0]
			}
		} else if bsonTag == "-" {
			continue
		}
		fields = append(fields, bsonTag)
		indexTag, ok := field.Tag.Lookup("index")
		if ok {
			indexModel := mongo.IndexModel{
				Keys:    bson.M{bsonTag: 1},
				Options: options.Index(),
			}
			if indexTag != "" {
				indexModel.Options.SetName(indexTag)
			}
			indexes = append(indexes, indexModel)
		}
	}
	schema := bson.M{}
	for _, field := range fields {
		schema[field] = 1
	}

	// Create the collection if it doesn't exist.
	if err := coll.Database().CreateCollection(ctx, name, &options.CreateCollectionOptions{
		Collation: collation,
	}); err != nil {
		return err
	}

	// Create the indexes if there are any.
	if len(indexes) > 0 {
		_, err := coll.Indexes().CreateMany(ctx, indexes)
		if err != nil {
			return err
		}
	}

	// Update the schema.
	if _, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: schema, Options: options.Index().SetName(fmt.Sprintf("%s_schema", name))}); err != nil {
		return err
	}

	return nil
}
