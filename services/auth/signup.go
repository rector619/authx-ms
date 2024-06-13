package auth

import (
	"net/http"
	"strings"

	"github.com/SineChat/auth-ms/external/external_models"
	"github.com/SineChat/auth-ms/external/request"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
)

func SignupService(Logger *utility.Logger, db *mongodb.Database, req models.SignupRequest) (int, error) {
	var (
		firstname    = strings.Title(strings.ToLower(req.Firstname))
		lastname     = strings.Title(strings.ToLower(req.Lastname))
		emailAddress = strings.ToLower(req.Email)
		password     = req.Password
	)

	password, err := utility.Hash(req.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user := models.User{
		FirstName: firstname,
		LastName:  lastname,
		Email:     emailAddress,
		Password:  password,
	}

	err = user.CreateUser(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// send welcome email
	data := external_models.SendWelcomeMail{
		Email: user.Email,
	}
	notification := request.ExternalRequest{
		Logger: Logger,
	}

	_, err = notification.SendExternalRequest(request.SendWelcomeMail, data)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
