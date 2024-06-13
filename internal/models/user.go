package models

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (User) CollectionName() string {
	return "users"
}

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName  string             `bson:"firstname" json:"firstname"`
	LastName   string             `bson:"lastname" json:"lastname"`
	Email      string             `bson:"email" index:"email,unique" json:"email"`
	Phone      string             `bson:"phone" index:"phone" json:"phone"`
	Password   string             `bson:"password" json:"password"`
	Country    string             `bson:"country" json:"country"`
	IsVerified bool               `bson:"is_verified" json:"is_verified"`
	IsAdmin    bool               `bson:"is_admin" json:"is_admin"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt  time.Time          `bson:"deleted_at" json:"-"`
	Deleted    bool               `bson:"deleted" json:"-"`
}

var (
	MyIdentity         *User
	IdentityLoginToken *LoginToken
)

type SignupRequest struct {
	Email     string `json:"email" validate:"required" mgvalidate:"notexists=auth$users$email,email"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName string `bson:"firstname" json:"firstname"`
	LastName  string `bson:"lastname" json:"lastname"`
	Email     string `bson:"email" json:"email" mgvalidate:"notexists=auth$users$email,email"`
	Phone     string `bson:"phone" json:"phone" mgvalidate:"notexists=auth$users$phone,phone"`
	Country   string `bson:"country" json:"country"`
}

type GetUserModel struct {
	ID           string `json:"id" pgvalidate:"exists=auth$users$id"`
	EmailAddress string `json:"email_address" pgvalidate:"exists=auth$users$email"`
	// PhoneNumber  string `json:"phone_number"`
}

func (u *User) CreateUser(db *mongodb.Database) error {
	err := db.CreateOneRecord(&u)
	if err != nil {
		return fmt.Errorf("user creation failed: %v", err.Error())
	}
	return nil
}

func (u *User) GetUserByID(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&u, bson.M{"_id": u.ID})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserByEmail(db *mongodb.Database) error {
	err := db.SelectOneFromDb(&u, bson.M{"email": u.Email})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateAllfields(db *mongodb.Database) error {

	err := db.SaveAllFields(&u)
	if err != nil {
		return fmt.Errorf("user update failed: %v", err.Error())
	}
	return nil
}
