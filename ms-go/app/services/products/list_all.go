package products

import (
	"context"
	"ms-go/app/helpers"
	"ms-go/app/models"
	"ms-go/db"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func ListAll() ([]models.Product, error) {
	var products []models.Product

	data, err := db.Connection().Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	for data.Next(context.TODO()) {
		var product models.Product
		if err := data.Decode(&product); err != nil {
			return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
		}
		products = append(products, product)
	}

	return products, nil
}
