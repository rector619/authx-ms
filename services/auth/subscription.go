package auth

import (
	"net/http"
	"strings"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSubscriptionService(Logger *utility.Logger, db *mongodb.Database, req models.CreateSubscriptionRequest) (int, error) {

	var (
		subscriptionType = strings.ToUpper(req.SubscriptionType)
		description      = req.Description
		currencyCode     = strings.ToUpper(req.CurrencyCode)
		price            = req.Price
	)

	subscription := models.Subscription{
		SubscriptionType: subscriptionType,
		Description:      description,
		CurrencyCode:     currencyCode,
		Price:            price,
	}

	err := subscription.CreateSubscription(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func GetSubscriptionService(Logger *utility.Logger, db *mongodb.Database) ([]models.Subscription, int, error) {
	sub := models.Subscription{}
	result, err := sub.GetSubscriptions(db)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func GetActiveSubscriptionService(Logger *utility.Logger, db *mongodb.Database) ([]models.Subscription, int, error) {
	sub := models.Subscription{}
	result, err := sub.GetActiveSubscriptions(db)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func UpdateSubscriptionService(Logger *utility.Logger, db *mongodb.Database, subID string, req models.UpdateSubscriptionRequest) (int, error) {
	id, err := primitive.ObjectIDFromHex(subID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	sub := models.Subscription{ID: id}
	err = sub.GetByID(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if req.SubscriptionType != "" {
		sub.SubscriptionType = strings.ToUpper(req.SubscriptionType)
	}
	if req.Description != "" {
		sub.Description = req.Description
	}
	if req.CurrencyCode != "" {
		sub.CurrencyCode = strings.ToUpper(req.CurrencyCode)
	}
	if req.Price > 0 {
		sub.Price = req.Price
	}
	if req.IsActive != nil {
		sub.IsActive = *req.IsActive
	}

	err = sub.UpdateSubscription(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func DeleteSubscriptionService(Logger *utility.Logger, db *mongodb.Database, subID string) (int, error) {
	id, err := primitive.ObjectIDFromHex(subID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	subscription := models.Subscription{ID: id}
	err = subscription.DeleteSubscription(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
