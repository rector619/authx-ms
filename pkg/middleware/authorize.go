package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SineChat/auth-ms/internal/config"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AuthType       AuthorizationType = "auth"
	AdminType      AuthorizationType = "admin"
	ApiPublicType  AuthorizationType = "api_public"
	ApiPrivateType AuthorizationType = "api_private"
	AppType        AuthorizationType = "app"
)

type (
	AuthorizationType  string
	AuthorizationTypes []AuthorizationType
)

func Authorize(db *mongodb.Database, authTypes ...AuthorizationType) gin.HandlerFunc {

	return func(c *gin.Context) {
		if len(authTypes) > 0 {

			msg := ""
			for _, v := range authTypes {
				ms, status := v.ValidateAuthorizationRequest(c, db)
				if status {
					return
				}
				msg = ms
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.UnauthorisedResponse(http.StatusUnauthorized, fmt.Sprint(http.StatusUnauthorized), "Unauthorized", msg))
		}
	}
}

func (at AuthorizationType) in(authTypes AuthorizationTypes) bool {
	for _, v := range authTypes {
		if v == at {
			return true
		}
	}
	return false
}

func (at AuthorizationType) ValidateAuthorizationRequest(c *gin.Context, db *mongodb.Database) (string, bool) {

	switch at {
	case ApiPublicType:
		return at.ValidateApiPublicType(c, db)
	case ApiPrivateType:
		return at.ValidateApiPrivateType(c, db)
	case AppType:
		return at.ValidateAppType(c)
	case AuthType:
		return at.ValidateAuthType(c, db, false)
	case AdminType:
		return at.ValidateAuthType(c, db, true)

	}

	return "authorized", true
}

func (at AuthorizationType) ValidateAuthType(c *gin.Context, db *mongodb.Database, isAdmin bool) (string, bool) {

	var invalidToken = "Your request was made with invalid credentials."
	authorizationToken := GetHeader(c, "Authorization")
	if authorizationToken == "" {
		return "token not provided", false
	}

	bearerTokenArr := strings.Split(authorizationToken, " ")
	if len(bearerTokenArr) != 2 {
		return invalidToken, false
	}

	bearerToken := bearerTokenArr[1]

	if bearerToken == "" {
		return invalidToken, false
	}

	token, err := TokenValid(bearerToken)
	if err != nil {
		return invalidToken, false
	}

	claims := token.Claims.(jwt.MapClaims)

	activeUserAccountIDStr, ok := claims["account_id"].(string)
	if !ok {
		return invalidToken, false
	}
	activeUserAccountID, err := primitive.ObjectIDFromHex(activeUserAccountIDStr)
	if err != nil {
		return err.Error(), false
	}

	tokenIDStr, ok := claims["token_id"].(string)
	if !ok {
		return invalidToken, false
	}

	tokenID, err := primitive.ObjectIDFromHex(tokenIDStr)
	if err != nil {
		return err.Error(), false
	}

	loginToken := models.LoginToken{ID: tokenID}
	err = loginToken.GetByID(db)
	if err != nil {
		return invalidToken, false
	}

	authoriseStatus, ok := claims["authorised"].(bool) //check if token is authorised for middleware
	if !ok && !authoriseStatus {
		return invalidToken, false
	}

	if time.Now().After(loginToken.AtExpiresTime) {
		return "expired token", false
	}

	user := models.User{ID: activeUserAccountID}
	err = user.GetUserByID(db)
	if err != nil {
		return "user does not exist: " + err.Error(), false
	}

	if isAdmin && !user.IsAdmin {
		return "access denied", false
	}

	models.MyIdentity = &user
	models.IdentityLoginToken = &loginToken

	return "authorized", true
}

func (at AuthorizationType) ValidateAppType(c *gin.Context) (string, bool) {
	config := config.GetConfig().App
	appKey := GetHeader(c, "app-key")
	if appKey == "" {
		return "missing app key", false
	}

	if appKey != config.Key {
		return "invalid app key", false
	}

	return "authorized", true
}

func (at AuthorizationType) ValidateApiPublicType(c *gin.Context, db *mongodb.Database) (string, bool) {
	_, msg, status := at.CheckAccessTokensWithPublicKey(c, db)
	return msg, status
}
func (at AuthorizationType) ValidateApiPrivateType(c *gin.Context, db *mongodb.Database) (string, bool) {
	_, msg, status := at.CheckAccessTokensWithPrivateKey(c, db)
	return msg, status
}

func (at AuthorizationType) CheckAccessTokensWithPublicKey(c *gin.Context, db *mongodb.Database) (models.AccessToken, string, bool) {
	publicKey := GetHeader(c, "public-key")

	if publicKey == "" {
		return models.AccessToken{}, "missing api key", false
	}

	token := models.AccessToken{PublicKey: publicKey, IsLive: true}
	err := token.LiveTokensWithPublicKey(db)
	if err != nil {
		return token, "invalid key", false
	}
	return token, "authorized", true
}

func (at AuthorizationType) CheckAccessTokensWithPrivateKey(c *gin.Context, db *mongodb.Database) (models.AccessToken, string, bool) {
	privateKey := GetHeader(c, "private-key")

	if privateKey == "" {
		return models.AccessToken{}, "missing api key", false
	}

	token := models.AccessToken{PrivateKey: privateKey, IsLive: true}
	err := token.LiveTokensWithPrivateKey(db)
	if err != nil {
		return token, "invalid key", false
	}
	return token, "authorized", true
}

func GetHeader(c *gin.Context, key string) string {
	header := ""
	if c.GetHeader(key) != "" {
		header = c.GetHeader(key)
	} else if c.GetHeader(strings.ToLower(key)) != "" {
		header = c.GetHeader(strings.ToLower(key))
	} else if c.GetHeader(strings.ToUpper(key)) != "" {
		header = c.GetHeader(strings.ToUpper(key))
	} else if c.GetHeader(strings.Title(key)) != "" {
		header = c.GetHeader(strings.Title(key))
	}
	return header
}
