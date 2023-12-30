package src

import (
	"context"
	"fmt"

	"github.com/kirvader/BodyController/domains/nutrition/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateNutritionLifestyleTemplateRequest struct {
	HexId   string
	NewData *models.NutritionLifestyleTemplate
}

type UpdateNutritionLifestyleTemplateResponse struct {
	HexId string
}

func (svc *NutritionLifestyleTemplateService) UpdateIngredient(ctx context.Context, req *UpdateNutritionLifestyleTemplateRequest) (*UpdateNutritionLifestyleTemplateResponse, error) {
	coll := svc.mongoClient.Database("BodyController").Collection("NutritionLifestyleTemplates")

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

	return &UpdateNutritionLifestyleTemplateResponse{
		HexId: req.HexId,
	}, nil
}
