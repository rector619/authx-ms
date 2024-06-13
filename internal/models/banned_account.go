package models

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (BannedAccount) CollectionName() string {
	return "banned_accounts"
}

type BannedAccount struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID primitive.ObjectID `bson:"account_id" json:"account_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at" json:"-"`
	Deleted   bool               `bson:"deleted" json:"-"`
}

type UnBanAccountRequest struct {
	Email string `json:"email" validate:"required" mgvalidate:"exists=auth$users$email,email"`
}

func (b *BannedAccount) CreateBannedAccount(db *mongodb.Database) error {
	err := db.CreateOneRecord(&b)
	if err != nil {
		return fmt.Errorf("banned account creation failed: %v", err.Error())
	}
	return nil
}

func (b *BannedAccount) DeleteBannedAccountByAccountID(db *mongodb.Database) error {
	err := db.HardDeleteByFilter(&b, bson.M{"account_id": b.AccountID})
	if err != nil {
		return fmt.Errorf("banned account deletion failed: %v", err.Error())
	}
	return nil
}

func (b *BannedAccount) GetBannedAccountByAccountID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&b, bson.M{"account_id": b.AccountID})

	if err != nil {
		return err
	}
	return nil
}
