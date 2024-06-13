package auth_model

import (
	"fmt"
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserService(db *mongodb.Database, req models.GetUserModel) (*models.User, int, error) {
	user := models.User{}
	if req.ID != "" {
		id, err := primitive.ObjectIDFromHex(req.ID)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		user.ID = id
		err = user.GetUserByID(db)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return &user, http.StatusOK, nil
	} else if req.EmailAddress != "" {
		user.Email = req.EmailAddress
		err := user.GetUserByEmail(db)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return &user, http.StatusOK, nil
	} else {
		return nil, http.StatusBadRequest, fmt.Errorf("no request values provided")
	}
}
