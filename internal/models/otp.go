package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (OTP) CollectionName() string {
	return "otps"
}

type OTP struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID   primitive.ObjectID `bson:"account_id,omitempty" json:"account_id"`
	Token       int                `bson:"token" json:"token"`               // 6 digit
	TokenType   string             `bson:"token_type" json:"token_type"`     // i.e. "reset_password"
	TokenLength int                `bson:"token_length" json:"token_length"` // i.e. length of token
	Duration    int32              `bson:"duration" json:"duration"`         // i.e. duration in seconds
	IsActive    bool               `bson:"is_active" json:"-"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt   time.Time          `bson:"deleted_at" json:"-"`
	Deleted     bool               `bson:"deleted" json:"-"`
}

func (otp *OTP) CreateOTP(db *mongodb.Database) error {
	if otp.TokenType == "" {
		return fmt.Errorf("token type not provided")
	}

	if otp.Duration == 0 {
		otp.Duration = 60
	}

	if otp.TokenLength == 0 {
		otp.TokenLength = 6
	}

	otp.Token = utility.RandomNumbers(otp.TokenLength)

	err := db.CreateOneRecord(&otp)
	if err != nil {
		return fmt.Errorf("token creation failed: %v", err.Error())
	}

	// TTL index: automatically remove documents from collection after a certain duration
	var hr_24 int32 = 3600 * 24
	sessionCollection := db.DB.Collection(otp.CollectionName())
	index := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1},
		Options: options.Index().SetExpireAfterSeconds(hr_24), // created OTP's Will be removed after 24 hours
	}

	_, err = sessionCollection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return err
	}

	return nil
}

func (otp *OTP) CheckIfExpired(db *mongodb.Database, token int) (bool, error) {

	err := db.SelectOneFromDb(&otp, bson.M{"token": token})

	// if the error is `mongo: no documents in result, return 'invalid token'`
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = fmt.Errorf("invalid/expired token")
	}

	if err != nil {
		return false, fmt.Errorf("%v", err.Error())
	}

	expiresAt := otp.CreatedAt.Add(time.Second * time.Duration(otp.Duration))
	if expiresAt.Before(time.Now()) {
		return false, fmt.Errorf("token has expired")
	}

	return true, nil
}

func (otp *OTP) DeleteOTP(db *mongodb.Database) error {
	err := db.HardDelete(&otp)
	if err != nil {
		return fmt.Errorf("otp token deletion failed: %v", err.Error())
	}
	return nil
}

func (otp *OTP) GetAccountIDByToken(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&otp, bson.M{"token": otp.Token})
	if err != nil {
		return err
	}
	return nil
}
