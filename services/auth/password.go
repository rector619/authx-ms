package auth

import (
	"fmt"
	"net/http"

	"github.com/SineChat/auth-ms/external/external_models"
	"github.com/SineChat/auth-ms/external/request"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
)

func PasswordResetService(Logger *utility.Logger, db *mongodb.Database, req models.PasswordReset) (int, error) {

	var (
		otp  models.OTP
		user models.User
	)

	// get otp
	otp = models.OTP{Token: req.Token}
	err := otp.GetAccountIDByToken(db)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid/expired token")
	}

	// get user
	user = models.User{ID: otp.AccountID}
	err = user.GetUserByID(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// check if token has expired
	if val, err := otp.CheckIfExpired(db, req.Token); val != true {
		return http.StatusBadRequest, err
	}

	// check if the token type is password reset
	if otp.TokenType != "password_reset" {
		return http.StatusBadRequest, fmt.Errorf("invalid token")
	}

	// hash password
	hashedPassword, err := utility.Hash(req.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user.Password = hashedPassword

	err = user.UpdateAllfields(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// delete otp immediately after use
	err = otp.DeleteOTP(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func PasswordResetRequestService(Logger *utility.Logger, db *mongodb.Database, req models.PasswordResetRequest) (int, error) {

	// get user
	user := models.User{Email: req.Email}
	err := user.GetUserByEmail(db)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid email")
	}

	otp, err := generatePasswordResetToken(db, user)
	if err != nil {
		return http.StatusBadRequest, err
	}

	redirect_url := utility.ParseUrl(req.RedirectURL, fmt.Sprint(otp.Token))

	// SEND EMAIL: password reset mail
	data := external_models.SendResetPasswordMail{
		Email:       user.Email,
		RedirectURL: redirect_url,
	}

	notification := request.ExternalRequest{
		Logger: Logger,
	}

	_, err = notification.SendExternalRequest(request.SendResetPasswordMail, data)
	if err != nil {
		return http.StatusBadRequest, err
	}
	// End

	return http.StatusOK, nil
}

func generatePasswordResetToken(db *mongodb.Database, user models.User) (models.OTP, error) {

	var otp models.OTP

	otp = models.OTP{AccountID: user.ID, TokenType: "password_reset", TokenLength: 10, Duration: 120}

	err := otp.CreateOTP(db)
	if err != nil {
		return otp, err
	}

	return otp, nil
}
