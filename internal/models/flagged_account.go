package models

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (FlaggedAccount) CollectionName() string {
	return "flagged_accounts"
}

type FlaggedAccount struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID        primitive.ObjectID `bson:"account_id" json:"account_id"`
	NoOfTimesFlagged int                `bson:"no_of_times_flagged" json:"no_of_times_flagged"`
	Reasons          []FlagReason       `bson:"reasons" json:"reasons"`
	IsFlagged        bool               `bson:"is_flagged" json:"is_flagged"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt        time.Time          `bson:"deleted_at" json:"-"`
	Deleted          bool               `bson:"deleted" json:"-"`
}

type FlagReason struct {
	Reason    string    `bson:"reason" json:"reason"` // reason for being flagged
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type FlagAccountRequest struct {
	Email  string `json:"email" validate:"required" mgvalidate:"exists=auth$users$email,email"`
	Reason string `json:"reason" validate:"required"`
}

type UnFlagAccountRequest struct {
	Email string `json:"email" validate:"required"`
}

func (f *FlaggedAccount) CreateFlaggedAccount(db *mongodb.Database) error {
	err := db.CreateOneRecord(&f)
	if err != nil {
		return fmt.Errorf("flagged user creation failed: %v", err.Error())
	}
	return nil
}

func (f *FlaggedAccount) UpdateFlaggedAccount(db *mongodb.Database) error {
	err := db.SaveAllFields(&f)
	if err != nil {
		return fmt.Errorf("flagged account update failed: %v", err.Error())
	}
	return nil
}

func (f *FlaggedAccount) GetFlaggedAccountByAccountID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&f, bson.M{"account_id": f.AccountID})
	if err != nil {
		return err
	}
	return nil
}

func (f *FlaggedAccount) DeleteFlaggedAccountByAccountID(db *mongodb.Database) error {
	err := db.HardDeleteByFilter(&f, bson.M{"account_id": f.AccountID})
	if err != nil {
		return fmt.Errorf("flagged account deletion failed: %v", err.Error())
	}
	return nil
}

func (f *FlaggedAccount) GetFlaggedAccounts(db *mongodb.Database) ([]FlaggedAccount, error) {

	var result []FlaggedAccount

	query := bson.M{"is_flagged": true}

	err := db.SelectAllFromDb("updated_at", &f, query, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// This appends the new data into FlaggedAccountReason
func (f *FlaggedAccount) AddFlagReason(reason FlagReason) {
	reason.CreatedAt = time.Now()
	f.Reasons = append(f.Reasons, reason)
}
