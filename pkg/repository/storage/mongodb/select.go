package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *Database) SelectAllFromDb(order string, receiver interface{}, query map[string]interface{}, result interface{}) error {
	query = AddDefaultGetParams(query)
	order, sortOrder := getSortDetails(order)
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return err
	}
	cur, err := coll.Find(context.Background(), query, options.Find().SetSort(bson.D{{order, sortOrder}}))
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	return cur.All(context.Background(), result)
}
func (db *Database) SelectAllFromDbWithLimit(order string, limit int, receiver interface{}, query map[string]interface{}, result interface{}) error {
	query = AddDefaultGetParams(query)
	order, sortOrder := getSortDetails(order)
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return err
	}
	cur, err := coll.Find(context.Background(), query, options.Find().SetSort(bson.D{{order, sortOrder}}).SetLimit(int64(limit)))
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	return cur.All(context.Background(), result)
}

func (db *Database) SelectOneFromDb(receiver interface{}, query map[string]interface{}) error {
	query = AddDefaultGetParams(query)
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return err
	}
	// filter := bson.M{"$and": []bson.M{{"_id": getDocumentID(receiver)}, {query: args}}}
	return coll.FindOne(context.Background(), query).Decode(receiver)
}

func (db *Database) SelectLatestFromDb(receiver interface{}, query map[string]interface{}) error {
	query = AddDefaultGetParams(query)
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return err
	}
	// filter := bson.M{"$and": []bson.M{{query: args}}}
	opt := options.FindOne().SetSort(bson.M{"_id": -1})
	return coll.FindOne(context.Background(), query, opt).Decode(receiver)
}

func (db *Database) SelectRandomFromDb(receiver interface{}, query map[string]interface{}) error {
	query = AddDefaultGetParams(query)
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return err
	}
	// filter := bson.M{"$and": []bson.M{{query: args}}}
	opt := options.FindOne().SetSort(bson.M{"rand()": 1})
	return coll.FindOne(context.Background(), query, opt).Decode(receiver)
}

func (db *Database) SelectFirstFromDb(receiver interface{}) error {
	query := AddDefaultGetParams(map[string]interface{}{})
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return err
	}
	return coll.FindOne(context.Background(), query).Decode(receiver)
}

func (db *Database) CheckExists(receiver interface{}, query map[string]interface{}) bool {
	query = AddDefaultGetParams(query)
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return false
	}
	// filter := bson.M{"$and": []bson.M{{"_id": getDocumentID(receiver)}, {query: args}}}
	count, err := coll.CountDocuments(context.Background(), query)
	return err == nil && count > 0
}

func (db *Database) CheckExistsInTable1(table string, query map[string]interface{}) bool {
	query = AddDefaultGetParams(query)
	coll := db.GetCollection(table)
	// filter := bson.M{query: args}
	count, err := coll.CountDocuments(context.Background(), query)
	return err == nil && count > 0
}

func (db *Database) CheckExistsInTable(table string, query map[string]interface{}) bool {
	query = AddDefaultGetParams(query)
	coll := db.GetCollection(table)
	// filter := bson.M{query: args}
	count, err := coll.CountDocuments(context.Background(), query)
	return err == nil && count > 0
}

func AddDefaultGetParams(query map[string]interface{}) map[string]interface{} {
	query["deleted"] = false
	return query
}

func getSortDetails(order string) (string, int) {
	if order == "" {
		order = "-_id"
	}

	sortOrder := 1
	if order[0] == '-' {
		sortOrder = -1
		order = order[1:]
	}

	return order, sortOrder
}
