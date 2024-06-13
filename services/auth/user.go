package auth

import (
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
)

func UpdateUserService(Logger *utility.Logger, db *mongodb.Database, user models.User, req models.UpdateUserRequest) (int, error) {

	// if the user email is changed, set verified to false
	if req.Email != "" {
		user.IsVerified = false
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Country != "" {
		user.Country = req.Country
	}

	err := user.UpdateAllfields(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
