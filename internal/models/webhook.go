package models

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (WebhookLog) CollectionName() string {
	return "webhook_log"
}

type WebhookLog struct {
	ID              primitive.ObjectID        `bson:"_id,omitempty" json:"id"`
	TransactionID   int                       `bson:"transaction_id" json:"transaction_id"`
	TransactionRef  string                    `bson:"transaction_ref" json:"transaction_ref"`
	ResponsePayload FlutterWaveWebhookRequest `bson:"response_payload" json:"response_payload"`
	CreatedAt       time.Time                 `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time                 `bson:"updated_at" json:"updated_at"`
	DeletedAt       time.Time                 `bson:"deleted_at" json:"-"`
	Deleted         bool                      `bson:"deleted" json:"-"`
}

func (w *WebhookLog) CreateWebhookLog(db *mongodb.Database) error {
	err := db.CreateOneRecord(&w)
	if err != nil {
		return fmt.Errorf("webhook creation failed: %v", err.Error())
	}
	return nil
}

func (w *WebhookLog) GetByTransactionID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&w, bson.M{"transaction_id": w.TransactionID})
	if err != nil {
		return err
	}
	return nil
}
