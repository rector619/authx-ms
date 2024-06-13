package authorization

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SineChat/auth-ms/internal/config"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/middleware"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateAuthorizationService(req models.ValidateAuthorizationReq, db *mongodb.Database) (interface{}, string, bool, int, error) {
	switch req.Type {
	case string(middleware.ApiPublicType):
		data, msg, status := validateApiPublicType(db, req.PublicKey)
		return data, msg, status, http.StatusOK, nil
	case string(middleware.ApiPrivateType):
		data, msg, status := validateApiPrivateType(db, req.PrivateKey)
		return data, msg, status, http.StatusOK, nil
	case string(middleware.AppType):
		data, msg, status := validateAppType(db, req.AppKey)
		return data, msg, status, http.StatusOK, nil
	case string(middleware.AuthType):
		data, msg, status := validateAuthType(db, req.AuthorizationToken)
		return data, msg, status, http.StatusOK, nil
	default:
		return nil, "not implemented", false, http.StatusBadRequest, fmt.Errorf("not implemented")
	}
}

func validateAuthType(db *mongodb.Database, bearerToken string) (*models.User, string, bool) {

	var invalidToken = "Your request was made with invalid credentials."

	if bearerToken == "" {
		return nil, invalidToken, false
	}

	token, err := middleware.TokenValid(bearerToken)
	if err != nil {
		return nil, invalidToken, false
	}

	claims := token.Claims.(jwt.MapClaims)

	activeUserAccountIDStr, ok := claims["account_id"].(string)
	if !ok {
		return nil, invalidToken, false
	}
	activeUserAccountID, err := primitive.ObjectIDFromHex(activeUserAccountIDStr)
	if err != nil {
		return nil, err.Error(), false
	}

	tokenIDStr, ok := claims["token_id"].(string)
	if !ok {
		return nil, invalidToken, false
	}

	tokenID, err := primitive.ObjectIDFromHex(tokenIDStr)
	if err != nil {
		return nil, err.Error(), false
	}

	loginToken := models.LoginToken{ID: tokenID}
	err = loginToken.GetByID(db)
	if err != nil {
		return nil, invalidToken, false
	}

	authoriseStatus, ok := claims["authorised"].(bool) //check if token is authorised for middleware
	if !ok && !authoriseStatus {
		return nil, invalidToken, false
	}

	if time.Now().After(loginToken.AtExpiresTime) {
		return nil, "expired token", false
	}

	user := models.User{ID: activeUserAccountID}
	err = user.GetUserByID(db)
	if err != nil {
		return nil, "user does not exist: " + err.Error(), false
	}

	return &user, "authorized", true
}

func validateAppType(db *mongodb.Database, appKey string) (interface{}, string, bool) {
	if appKey == "" {
		return nil, "missing app key", false
	}

	if appKey != config.GetConfig().App.Key {
		return nil, "invalid app key", false
	}

	return nil, "authorized", true
}

func validateApiPublicType(db *mongodb.Database, publicKey string) (*models.AccessToken, string, bool) {
	if publicKey == "" {
		return &models.AccessToken{}, "missing api key", false
	}

	token := models.AccessToken{PublicKey: publicKey, IsLive: true}
	err := token.LiveTokensWithPublicKey(db)
	if err != nil {
		return &token, "invalid key", false
	}
	return &token, "authorized", true
}

func validateApiPrivateType(db *mongodb.Database, privateKey string) (*models.AccessToken, string, bool) {
	if privateKey == "" {
		return &models.AccessToken{}, "missing api key", false
	}

	token := models.AccessToken{PrivateKey: privateKey, IsLive: true}
	err := token.LiveTokensWithPrivateKey(db)
	if err != nil {
		return &token, "invalid key", false
	}
	return &token, "authorized", true
}
