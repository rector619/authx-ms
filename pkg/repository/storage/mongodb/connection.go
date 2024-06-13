package mongodb

import (
	"fmt"
	"strings"
	"time"

	"context"

	"github.com/SineChat/auth-ms/internal/config"
	"github.com/SineChat/auth-ms/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	DB     *mongo.Database
	Logger *utility.Logger
}

var DB Database

// Connection gets connection of mysqlDB database
func Connection() *Database {
	return &DB
}

func ConnectToDB(logger *utility.Logger, conectionString string) (*Database, error) {
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions := options.Client().ApplyURI(conectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("error connecting to db %v", err.Error()))
		panic(err)
	}
	utility.LogAndPrint(logger, "connected to db")
	db := Database{DB: client.Database(config.Config.Databases.DBName), Logger: logger}
	DB = db
	return &db, nil
}

func (db *Database) GetCollection(collectionName string) *mongo.Collection {
	return db.DB.Collection(collectionName)
}

func (db *Database) CreateUniqueIndex(logger *utility.Logger, collName, field string, order int) error {
	collection := db.GetCollection(collName)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{field: order},
		Options: options.Index().SetUnique(true),
	}

	timeOutFactor := 3
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutFactor)*time.Second)

	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create unique index on field %s in %s", field, collName))
		return fmt.Errorf("failed to create unique index on field %s in %s", field, collName)
	}

	return nil
}

func (db *Database) CollectionExists(ctx context.Context, name string) bool {
	filter := bson.M{"name": name}
	collections, err := db.DB.ListCollectionNames(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, coll := range collections {
		if strings.EqualFold(coll, name) {
			return true
		}
	}
	return false
}

func (db *Database) GetCollectionNameForModel(model interface{}) (string, error) {
	ctx := context.Background()
	name := CollectionName(model)
	if !db.CollectionExists(ctx, name) {
		return "", fmt.Errorf("collection for model %v does not exist, add to migrations and apply", name)
	}
	return name, nil
}

func (db *Database) GetCollectionForModel(model interface{}) (*mongo.Collection, error) {
	name, err := db.GetCollectionNameForModel(model)
	if err != nil {
		return nil, err
	}
	return db.GetCollection(name), nil
}
