package models

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (LoginToken) CollectionName() string {
	return "login_tokens"
}

type LoginToken struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID     primitive.ObjectID `bson:"account_id" json:"account_id"`
	AccessUuid    string             `bson:"access_uuid"  json:"access_uuid"`
	AccessToken   string             `bson:"access_token"  json:"access_token"`
	AtExpiresTime time.Time          `bson:"at_expires_time"  json:"at_expires_time"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt     time.Time          `bson:"deleted_at" json:"-"`
	Deleted       bool               `bson:"deleted" json:"-"`
}

func (a *LoginToken) GetByID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&a, bson.M{"_id": a.ID})

	if err != nil {
		return err
	}
	return nil
}

func (a *LoginToken) GetByAccountID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&a, bson.M{"account_id": a.AccountID})
	if err != nil {
		return err
	}
	return nil
}

func (l *LoginToken) CreateLoginToken(db *mongodb.Database) error {
	err := db.CreateOneRecord(&l)
	if err != nil {
		return fmt.Errorf("login token creation failed: %v", err.Error())
	}
	return nil
}

func (l *LoginToken) UpdateLoginToken(db *mongodb.Database) error {
	err := db.SaveAllFields(&l)
	if err != nil {
		return fmt.Errorf("login token update failed: %v", err.Error())
	}
	return nil
}

func (l *LoginToken) DeleteLoginToken(db *mongodb.Database) error {
	err := db.HardDelete(&l)
	if err != nil {
		return fmt.Errorf("login token deletion failed: %v", err.Error())
	}
	return nil
}

func (l *LoginToken) DeleteLoginTokensByAccountID(db *mongodb.Database) error {
	err := db.HardDeleteByFilter(&l, bson.M{"account_id": l.AccountID})
	if err != nil {
		return fmt.Errorf("login tokens deletion failed: %v", err.Error())
	}
	return nil
}

func (l *LoginToken) DeleteMultipleLoginTokensByAccountID(db *mongodb.Database) error {
	err := db.HardDeleteManyByFilter(&l, bson.M{"account_id": l.AccountID})
	if err != nil {
		return fmt.Errorf("login tokens deletion failed: %v", err.Error())
	}
	return nil
}
