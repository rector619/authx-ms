package auth

import (
	"fmt"
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FlagAccountService(Logger *utility.Logger, db *mongodb.Database, req models.FlagAccountRequest) (int, error) {

	// Get User
	user := models.User{
		Email: req.Email,
	}
	err := user.GetUserByEmail(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	flaggedAccount := models.FlaggedAccount{
		AccountID: user.ID,
	}

	// check if the flagged account already exist in the database.
	err = flaggedAccount.GetFlaggedAccountByAccountID(db)

	switch err {
	case nil: // account has been flagged before so -> nil

		// set the account as flagged if it has been unflagged before
		flaggedAccount.IsFlagged = true

		flaggedAccount.AddFlagReason(models.FlagReason{
			Reason: req.Reason,
		})

		flaggedAccount.NoOfTimesFlagged = len(flaggedAccount.Reasons)

		// if the account has been flagged 3 times, ban the account on the 4th time
		if flaggedAccount.NoOfTimesFlagged >= 4 {
			bannedAccount := models.BannedAccount{AccountID: user.ID}
			err = bannedAccount.CreateBannedAccount(db)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
		//

		// update record
		err = flaggedAccount.UpdateFlaggedAccount(db)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	default: // account hasnt been flagged before. create a new record

		flaggedAccount.IsFlagged = true // set the account as flagged

		flaggedAccount.AddFlagReason(models.FlagReason{
			Reason: req.Reason,
		})

		flaggedAccount.NoOfTimesFlagged = 1

		err = flaggedAccount.CreateFlaggedAccount(db)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Delete User Login Token from Db
	err = DeleteLoginTokensService(db, user.ID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error deleting login token: %v", err)
	}

	return http.StatusOK, nil
}

func DeleteLoginTokensService(db *mongodb.Database, accountID primitive.ObjectID) error {

	loginToken := models.LoginToken{
		AccountID: accountID,
	}

	err := loginToken.GetByAccountID(db)
	if err != nil {
		return err
	}

	err = loginToken.DeleteMultipleLoginTokensByAccountID(db)
	if err != nil {
		return err
	}

	return nil
}

func UnFlagAccountService(Logger *utility.Logger, db *mongodb.Database, req models.UnFlagAccountRequest) (int, error) {

	// Get User
	user := models.User{
		Email: req.Email,
	}

	err := user.GetUserByEmail(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	flaggedAccount := models.FlaggedAccount{
		AccountID: user.ID,
	}

	// get flagged account
	err = flaggedAccount.GetFlaggedAccountByAccountID(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	flaggedAccount.IsFlagged = false

	err = flaggedAccount.UpdateFlaggedAccount(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetFlaggedAccountsService(Logger *utility.Logger, db *mongodb.Database) ([]models.FlaggedAccount, int, error) {

	flag := models.FlaggedAccount{}

	result, err := flag.GetFlaggedAccounts(db)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}
