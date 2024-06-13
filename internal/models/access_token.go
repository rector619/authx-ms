package models

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/internal/config"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (AccessToken) CollectionName() string {
	return "access_tokens"
}

type AccessToken struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID  primitive.ObjectID `bson:"account_id" json:"account_id"`
	PublicKey  string             `bson:"public_key" json:"public_key"`
	PrivateKey string             `bson:"private_key" json:"private_key"`
	IsLive     bool               `bson:"is_live" json:"is_live"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt  time.Time          `bson:"deleted_at" json:"-"`
	Deleted    bool               `bson:"deleted" json:"-"`
}

func (a *AccessToken) GetByAccountID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&a, bson.M{"account_id": a.AccountID})

	if err != nil {
		return err
	}
	return nil
}

func (a *AccessToken) GetLatestByAccountIDAndIsLive(db *mongodb.Database) error {
	err := db.SelectLatestFromDb(&a, bson.M{"account_id": a.AccountID, "is_live": a.IsLive})
	if err != nil {
		return err
	}
	return nil
}

func (a *AccessToken) CreateAccessToken(db *mongodb.Database) error {
	app := config.GetConfig().App
	var basePrimitive primitive.ObjectID
	if a.AccountID == basePrimitive {
		return fmt.Errorf("account id not provided to create access token")
	}
	a.IsLive = true
	a.PrivateKey = "s_" + app.Name + "_" + utility.RandomString(50)
	a.PublicKey = "s_" + app.Name + "_" + utility.RandomString(50)

	err := db.CreateOneRecord(&a)
	if err != nil {
		return fmt.Errorf("access token creation failed: %v", err.Error())
	}
	return nil
}

func (a *AccessToken) LiveTokensWithPublicKey(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&a, bson.M{"public_key": a.PublicKey, "is_live": a.IsLive})

	if err != nil {
		return err
	}
	return nil
}
func (a *AccessToken) LiveTokensWithPrivateKey(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&a, bson.M{"private_key": a.PrivateKey, "is_live": a.IsLive})

	if err != nil {
		return err
	}
	return nil
}

func (a *AccessToken) LiveTokensWithKey(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&a, bson.M{
		"$or": bson.A{
			bson.M{"private_key": a.PrivateKey},
			bson.M{"public_key": a.PublicKey},
		},
		"is_live": a.IsLive,
	})

	if err != nil {
		return err
	}
	return nil
}
