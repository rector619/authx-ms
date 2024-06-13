package auth

import (
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
)

func FlutterWaveWebhookService(Logger *utility.Logger, db *mongodb.Database, req models.FlutterWaveWebhookRequest, requestBody []byte) (int, error) {

	webhook := models.WebhookLog{
		TransactionID:   req.ID,
		TransactionRef:  req.TxRef,
		ResponsePayload: req,
	}

	err := webhook.CreateWebhookLog(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func FlutterWaveWebhookVerifyPaymentService(Logger *utility.Logger, db *mongodb.Database, req models.FlutterWaveWebhookRequest) (int, error) {

	payment := models.Payment{
		TxRef: req.TxRef,
	}

	err := payment.GetByTxRef(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// check if payment has been verified already.
	if payment.IsVerified {
		return http.StatusOK, nil
	}

	if req.Status == "successful" {
		payment.IsPaid = true
	}

	payment.IsVerified = true
	payment.PaymentType = req.EventType

	err = payment.UpdateAllFields(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
