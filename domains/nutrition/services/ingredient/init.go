package src

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kirvader/BodyController/internal/db"
)

type IngredientService struct {
	mongoClient *mongo.Client
}

func NewIngredientService(ctx context.Context) (*IngredientService, func(), error) {
	mongoClient, disconnectMongoClient, err := db.InitMongoDBClientFromENV(ctx)
	if err != nil {
		return nil, func() {}, err
	}
	if err = db.PingMongoDb(ctx, mongoClient); err != nil {
		return nil, func() {
			disconnectMongoClient()
		}, err
	}

	return &IngredientService{
		mongoClient: mongoClient,
	}, disconnectMongoClient, nil
}
