package models

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (Payment) CollectionName() string {
	return "payments"
}

type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID   primitive.ObjectID `bson:"account_id" json:"account_id"`
	Amount      string             `bson:"amount" json:"amount"`
	Currency    string             `bson:"currency" json:"currency"`
	TxRef       string             `bson:"tx_ref" json:"tx_ref"`
	PaymentType string             `bson:"payment_type" json:"payment_type"`
	IsVerified  bool               `bson:"is_verified" json:"is_verified"`
	IsPaid      bool               `bson:"is_paid" json:"is_paid"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt   time.Time          `bson:"deleted_at" json:"-"`
	Deleted     bool               `bson:"deleted" json:"-"`
}

func (p *Payment) CreatePayment(db *mongodb.Database) error {
	err := db.CreateOneRecord(&p)
	if err != nil {
		return fmt.Errorf("payment record creation failed: %v", err.Error())
	}
	return nil
}

func (p *Payment) GetByTxRef(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&p, bson.M{"tx_ref": p.TxRef})
	if err != nil {
		return err
	}
	return nil
}

func (p *Payment) UpdateAllFields(db *mongodb.Database) error {
	err := db.SaveAllFields(&p)
	if err != nil {
		return fmt.Errorf("payment update failed: %v", err.Error())
	}
	return nil
}
