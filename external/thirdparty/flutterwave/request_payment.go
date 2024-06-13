package flutterwave

import (
	"fmt"

	"github.com/SineChat/auth-ms/external/external_models"
	"github.com/SineChat/auth-ms/internal/config"
)

func (r *RequestObj) FlutterWavePayment() (external_models.FlutterWavePaymentResponse, error) {

	var (
		outBoundResponse external_models.FlutterWavePaymentResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	// data, ok := idata.(models.PaymentRequest)
	data, ok := idata.(external_models.FlutterWavePaymentRequest)
	if !ok {
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.GetConfig().FlutterWave.SecretKey,
	}

	logger.Info("[FLUTTERWAVE]: request payment link", data)

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("[FLUTTERWAVE]: request payment link", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("[FLUTTERWAVE]: request payment link", outBoundResponse)

	return outBoundResponse, nil
}
