package middleware

import (
	"fmt"
	"time"

	"github.com/SineChat/auth-ms/internal/config"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

type TokenDetailsDTO struct {
	AccessUuid    string `json:"access_uuid"`
	AccessToken   string `json:"access_token"`
	AtExpiresTime time.Time
}

// CreateToken method
func CreateToken(db *mongodb.Database, user models.User, universalAccess bool) (*models.LoginToken, error) {

	config := config.GetConfig()

	var err error
	td := &models.LoginToken{
		AccountID: user.ID,
	}

	td.AtExpiresTime = time.Now().Add(time.Hour * time.Duration(config.Server.AccessTokenExpirationDuration))
	AccessUuid, _ := uuid.NewV4()
	td.AccessUuid = AccessUuid.String()

	err = td.CreateLoginToken(db)
	if err != nil {
		return nil, err
	}

	//create access token
	atClaims := jwt.MapClaims{}
	atClaims["email"] = user.Email
	atClaims["account_id"] = user.ID.Hex()
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["authorised"] = true
	atClaims["universal_access"] = universalAccess
	atClaims["token_id"] = td.ID.Hex()
	atClaims["exp"] = td.AtExpiresTime.Unix()

	// atClaims["exp"] = time.Now().AddDate(0, 0, 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = token.SignedString([]byte(config.Server.Secret))
	if err != nil {
		return nil, err
	}

	err = td.UpdateLoginToken(db)
	if err != nil {
		return nil, err
	}

	//generatedPass, err := GenerateSecureKey(16)
	if err != nil {
		return nil, err
	}
	//td.TransmissionKey = generatedPass

	return td, nil
}

// TokenValid method
func TokenValid(bearerToken string) (*jwt.Token, error) {
	token, err := verifyToken(bearerToken)
	if err != nil {
		if token != nil {
			return token, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Unauthorized")
	}
	return token, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	config := config.GetConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Server.Secret), nil
	})
	if err != nil {
		return token, fmt.Errorf("Unauthorized")
	}
	return token, nil
}
