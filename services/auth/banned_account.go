package auth

import (
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
)

func UnBanAccountService(Logger *utility.Logger, db *mongodb.Database, req models.UnBanAccountRequest) (int, error) {

	// Get User
	user := models.User{
		Email: req.Email,
	}

	err := user.GetUserByEmail(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	bannedAccount := models.BannedAccount{
		AccountID: user.ID,
	}

	// check if the specified account_id exists in the banned account collection
	err = bannedAccount.GetBannedAccountByAccountID(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// FLAG
	// when a user is unbanned, delete the flagged records
	flaggedAccount := models.FlaggedAccount{
		AccountID: user.ID,
	}

	// get flagged account
	err = flaggedAccount.GetFlaggedAccountByAccountID(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// delete flagged account record
	err = flaggedAccount.DeleteFlaggedAccountByAccountID(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// delete banned account record
	err = bannedAccount.DeleteBannedAccountByAccountID(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
