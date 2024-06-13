package flutterwave

import (
	"fmt"

	"github.com/SineChat/auth-ms/external/external_models"
	"github.com/SineChat/auth-ms/internal/config"
)

func (r *RequestObj) FlutterWaveVerifyPaymentByTxRef() (external_models.FlutterWaveVerifyPaymentResponse, error) {

	var (
		outBoundResponse external_models.FlutterWaveVerifyPaymentResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(string)
	if !ok {
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Authorization": "Bearer " + config.GetConfig().FlutterWave.SecretKey,
	}

	logger.Info("[FLUTTERWAVE]: verify payment", data)

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("[FLUTTERWAVE]: verify payment", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("[FLUTTERWAVE]: verify payment", outBoundResponse)

	return outBoundResponse, nil
}
