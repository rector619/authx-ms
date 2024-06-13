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

func EmailVerificationService(Logger *utility.Logger, db *mongodb.Database, req models.EmailVerification) (int, error) {

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

	// check if the token type is for email verification
	if otp.TokenType != "email_verification" {
		return http.StatusBadRequest, fmt.Errorf("invalid token")
	}

	// set the user to verified
	user.IsVerified = true

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

func EmailVerificationRequestService(Logger *utility.Logger, db *mongodb.Database, req models.EmailVerificationRequest) (int, error) {

	// get user
	user := models.User{Email: req.Email}
	err := user.GetUserByEmail(db)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid email")
	}

	// if the user has already been verified, return this
	if user.IsVerified {
		return http.StatusOK, fmt.Errorf("your account has already been verified")
	}

	// generate verification token
	otp, err := generateEmailVerificationToken(db, user)
	if err != nil {
		return http.StatusBadRequest, err
	}

	redirect_url := utility.ParseUrl(req.RedirectURL, fmt.Sprint(otp.Token))

	fmt.Println(redirect_url)
	// SEND EMAILS: verification mail
	data := external_models.SendVerificationMail{
		Email:       user.Email,
		RedirectURL: redirect_url,
	}
	notification := request.ExternalRequest{
		Logger: Logger,
	}
	_, err = notification.SendExternalRequest(request.SendVerificationMail, data)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func generateEmailVerificationToken(db *mongodb.Database, user models.User) (models.OTP, error) {

	var otp models.OTP

	var hour_24 int32 = 24 * 3600 // 24 hours

	otp = models.OTP{AccountID: user.ID, TokenType: "email_verification", TokenLength: 15, Duration: hour_24}

	err := otp.CreateOTP(db)
	if err != nil {
		return otp, err
	}

	return otp, nil
}
