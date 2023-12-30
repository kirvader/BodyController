package src

import (
	"context"
	"fmt"

	"github.com/kirvader/BodyController/domains/nutrition/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdatePersonalNutritionLifestyleRequest struct {
	HexId   string
	NewData *models.PersonalNutritionLifestyle
}

type UpdatePersonalNutritionLifestyleResponse struct {
	HexId string
}

func (svc *PersonalNutritionLifestyleService) UpdateIngredient(ctx context.Context, req *UpdatePersonalNutritionLifestyleRequest) (*UpdatePersonalNutritionLifestyleResponse, error) {
	coll := svc.mongoClient.Database("BodyController").Collection("PersonalNutritionLifestyles")

	objectId, err := primitive.ObjectIDFromHex(req.HexId)
	if err != nil {
		return nil, err
	}

	newDocumentData, err := req.NewData.ConvertToMongoDocument()
	if err != nil {
		return nil, err
	}

	resp, err := coll.UpdateByID(ctx, objectId, bson.M{"$set": newDocumentData})
	if err != nil {
		return nil, fmt.Errorf("update error: %w", err)
	}
	if resp.ModifiedCount == 0 {
		return nil, fmt.Errorf("no record updated: id='%s'", req.HexId)
	}

	return &UpdatePersonalNutritionLifestyleResponse{
		HexId: req.HexId,
	}, nil
}
