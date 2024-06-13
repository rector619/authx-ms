package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/middleware"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
)

func LoginService(Logger *utility.Logger, db *mongodb.Database, req models.LoginRequest) (map[string]interface{}, int, error) {
	var (
		emailAddress = strings.ToLower(req.Email)
		response     = map[string]interface{}{}
	)

	user := models.User{Email: emailAddress}
	err := user.GetUserByEmail(db)
	if err != nil {
		fmt.Println(err.Error())
		return response, http.StatusBadRequest, fmt.Errorf("invalid login details")
	}

	bannedAccount := models.BannedAccount{AccountID: user.ID}
	err = bannedAccount.GetBannedAccountByAccountID(db)
	if err == nil {
		return response, http.StatusBadRequest, fmt.Errorf("this account has been banned")
	}

	flaggedAccount := models.FlaggedAccount{AccountID: user.ID}
	err = flaggedAccount.GetFlaggedAccountByAccountID(db)
	if err == nil {
		if flaggedAccount.IsFlagged == true {
			return response, http.StatusBadRequest, fmt.Errorf("this account has been flagged")
		}
	}

	if !utility.CompareHash(req.Password, user.Password) {
		return response, http.StatusBadRequest, fmt.Errorf("invalid login details")
	}

	if !user.IsVerified {
		return response, http.StatusBadRequest, fmt.Errorf("your account needs to be verified before logging in")
	}

	token, err := middleware.CreateToken(db, user, false)
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("error creating token: " + err.Error())
	}

	response["user"] = user
	response["token"] = token.AccessToken
	return response, http.StatusOK, nil
}

func ChangeAuthUserPasswordService(Logger *utility.Logger, db *mongodb.Database, req models.ChangeAuthUserPassword, emailAddress string) (int, error) {

	user := models.User{Email: emailAddress}
	err := user.GetUserByEmail(db)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid login details")
	}

	if !utility.CompareHash(req.OldPassword, user.Password) {
		return http.StatusBadRequest, fmt.Errorf("old password does not match")
	}

	// hash password
	hashedPassword, err := utility.Hash(req.NewPassword)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user.Password = hashedPassword
	err = user.UpdateAllfields(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func AdminLoginService(Logger *utility.Logger, db *mongodb.Database, req models.LoginRequest) (map[string]interface{}, int, error) {
	var (
		emailAddress = strings.ToLower(req.Email)
		response     = map[string]interface{}{}
	)

	user := models.User{Email: emailAddress}
	err := user.GetUserByEmail(db)
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("invalid login details")
	}

	if !utility.CompareHash(req.Password, user.Password) {
		return response, http.StatusBadRequest, fmt.Errorf("invalid login details")
	}

	if !user.IsAdmin {
		return response, http.StatusBadRequest, fmt.Errorf("access denied; not an admin.")
	}

	if !user.IsVerified {
		return response, http.StatusBadRequest, fmt.Errorf("your account needs to be verified before logging in")
	}

	token, err := middleware.CreateToken(db, user, false)
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("error creating token: " + err.Error())
	}

	response["user"] = user
	response["token"] = token.AccessToken
	return response, http.StatusOK, nil
}
