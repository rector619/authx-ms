package auth_model

import (
	"fmt"
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
)

func ValidateOnDbService(req models.ValidateOnDBReq, db *mongodb.Database) (bool, int, error) {

	req.Query = mongodb.ValidateMapQuery(req.Query)
	fmt.Println(req.Query)
	if req.Type == "notexists" {
		return !db.CheckExistsInTable(req.Table, req.Query), http.StatusOK, nil

	} else if req.Type == "exists" {
		return db.CheckExistsInTable(req.Table, req.Query), http.StatusOK, nil

	} else {
		return false, http.StatusOK, nil
	}
}
