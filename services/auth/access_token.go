package auth

import (
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAccessTokenService(Logger *utility.Logger, db *mongodb.Database, userID primitive.ObjectID) (models.AccessToken, int, error) {

	token := models.AccessToken{AccountID: userID}

	err := token.GetByAccountID(db)
	// if the err is nil, it shows the access token with the AccountID exists
	if err == nil {
		return token, http.StatusOK, nil
	}

	// if it doesnt exist, create one
	err = token.CreateAccessToken(db)
	if err != nil {
		return token, http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil
}
