package products

import (
	"context"
	"fmt"
	"ms-go/app/helpers"
	"ms-go/app/models"
	"ms-go/app/producers"
	"ms-go/db"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(data models.Product, isAPI bool) (*models.Product, error) {
	if data.ID == 0 {
		var max models.Product
		opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

		if err := db.Connection().FindOne(context.TODO(), bson.D{}, opts).Decode(&max); err != nil {
			return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
		}

		data.ID = max.ID + 1
	}

	if err := data.Validate(); err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusUnprocessableEntity}
	}

	data.CreatedAt = time.Now()
	data.UpdatedAt = data.CreatedAt

	if _, err := db.Connection().InsertOne(context.TODO(), data); err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	topic := "go-to-rails"
	if err := producers.ProduceMessage(topic, data); err != nil {
		fmt.Printf("Error producing message for topic %s: %v\n", topic, err)
	}

	fmt.Println("Message successfully sent to topic", topic)

	return &data, nil
}
