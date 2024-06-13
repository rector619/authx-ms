package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (Subscription) CollectionName() string {
	return "subscriptions"
}

type Subscription struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SubscriptionType string             `bson:"subscription_type" json:"subscription_type"`
	Description      string             `bson:"description" json:"description"`
	Price            float64            `bson:"price" json:"price"`
	CurrencyCode     string             `bson:"currency_code" json:"currency_code"`
	IsActive         bool               `bson:"is_active" json:"-"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt        time.Time          `bson:"deleted_at" json:"-"`
	Deleted          bool               `bson:"deleted" json:"-"`
}

type CreateSubscriptionRequest struct {
	SubscriptionType string  `json:"subscription_type" validate:"required"`
	Description      string  `json:"description" validate:"required"`
	CurrencyCode     string  `json:"currency_code" validate:"required"`
	Price            float64 `json:"price" validate:"required"`
}
type UpdateSubscriptionRequest struct {
	SubscriptionType string  `json:"subscription_type"`
	Description      string  `json:"description"`
	CurrencyCode     string  `json:"currency_code"`
	Price            float64 `json:"price"`
	IsActive         *bool   `json:"is_active"`
}

func (s *Subscription) GetByID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&s, bson.M{"_id": s.ID})

	if err != nil {
		return err
	}
	return nil
}

func (s *Subscription) CreateSubscription(db *mongodb.Database) error {

	// a query to check if a susbcription with the same price and type exist.
	query := bson.M{"price": s.Price, "subscription_type": strings.ToUpper(s.SubscriptionType)}

	exist := db.CheckExists(&s, query)
	if exist {
		return fmt.Errorf("A subscription with the same TYPE and PRICE already exists")
	}

	err := db.CreateOneRecord(&s)
	if err != nil {
		return fmt.Errorf("subscription creation failed: %v", err.Error())
	}

	return nil
}

func (s *Subscription) GetSubscriptions(db *mongodb.Database) ([]Subscription, error) {

	var result []Subscription

	query := bson.M{}

	err := db.SelectAllFromDb("updated_at", &s, query, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *Subscription) GetActiveSubscriptions(db *mongodb.Database) ([]Subscription, error) {

	var result []Subscription

	query := bson.M{"is_active": true}

	err := db.SelectAllFromDb("updated_at", &s, query, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *Subscription) UpdateSubscription(db *mongodb.Database) error {
	s.SubscriptionType = strings.ToUpper(s.SubscriptionType)
	s.CurrencyCode = strings.ToUpper(s.CurrencyCode)
	err := db.SaveAllFields(&s)
	if err != nil {
		return fmt.Errorf("subscription update failed: %v", err.Error())
	}
	return nil
}

func (s *Subscription) DeleteSubscription(db *mongodb.Database) error {

	query := bson.M{"_id": s.ID}

	// check if the subscription with the specified id exists
	exist := db.CheckExists(&s, query)
	if !exist {
		return fmt.Errorf("invalid Subscription ID")
	}

	err := db.HardDeleteByFilter(&s, query)
	if err != nil {
		return err
	}

	return nil
}
