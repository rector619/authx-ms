package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SineChat/auth-ms/external/external_models"
	"github.com/SineChat/auth-ms/external/request"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
)

func PaymentRequestService(Logger *utility.Logger, db *mongodb.Database, req external_models.FlutterWavePaymentRequest, emailAddress string) (interface{}, int, error) {

	var response interface{}

	user := models.User{Email: emailAddress}
	err := user.GetUserByEmail(db)
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("invalid login details")
	}

	if req.Currency == "" {
		req.Currency = "NGN"
	}

	data := external_models.FlutterWavePaymentRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		RedirectURL: req.RedirectURL,
		Customer: struct {
			Email    string "json:\"email\""
			FullName string "json:\"name\""
		}{
			Email:    user.Email,
			FullName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		},
		TxRef: fmt.Sprint(time.Now().Unix()),
	}

	// external request
	flutterwave := request.ExternalRequest{
		Logger: Logger,
	}
	response, err = flutterwave.SendExternalRequest(request.FlutterWavePaymentRequest, data)
	if err != nil {
		return response, http.StatusBadRequest, err
	}
	// end

	// save transaction details to db
	payment := models.Payment{
		AccountID: user.ID,
		Amount:    data.Amount,
		Currency:  data.Currency,
		TxRef:     data.TxRef,
		IsPaid:    false,
	}

	err = payment.CreatePayment(db)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	return response, http.StatusOK, nil
}

// Send external request to flutterwave api to verify the transaction
func PaymentVerifyService(Logger *utility.Logger, db *mongodb.Database, txRef string) (interface{}, int, error) {

	flutterwave := request.ExternalRequest{
		Logger: Logger,
	}

	response, err := flutterwave.SendExternalRequest(request.FlutterWavePaymentVerify, txRef)
	if err != nil {
		return response, http.StatusBadRequest, err
	}

	return response, http.StatusOK, nil

}

func UpdatePaymentStatus(Logger *utility.Logger, db *mongodb.Database, resp external_models.FlutterWaveVerifyPaymentResponse) (int, error) {

	payment := models.Payment{
		TxRef: resp.Data.TxRef,
	}

	err := payment.GetByTxRef(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// check if the payment has been verified before
	if payment.IsVerified {
		return http.StatusOK, nil
	}

	// update these fields if status is successful
	if resp.Data.Status == "successful" {
		payment.IsPaid = true
		payment.PaymentType = resp.Data.PaymentType
	}

	payment.IsVerified = true

	err = payment.UpdateAllFields(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
