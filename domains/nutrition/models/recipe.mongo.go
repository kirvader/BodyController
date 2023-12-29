package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/durationpb"
)

type RecipeMongoDB struct {
	Id                             primitive.ObjectID          `bson:"_id,omitempty"`
	Title                          string                      `bson:"title"`
	TasteDescription               string                      `bson:"taste_description"`
	CookingStepsDescription        string                      `bson:"cooking_steps_description"`
	Author                         string                      `bson:"author"`
	RequiredIngredientsProportions []WeightedIngredientMongoDB `bson:"required_ingredients_proportions"`
	CookingTime                    int64                       `bson:"cooking_time"`
	Macros100G                     MacrosMongoDB               `bson:"macros_100g"`
}

func (proto *Recipe) ConvertToMongoDocument() (*RecipeMongoDB, error) {
	mongoMacros, _ := proto.Macros100G.ConvertToMongoDocument()

	mongo := &RecipeMongoDB{
		Title:                          proto.Title,
		TasteDescription:               proto.TasteDescription,
		CookingStepsDescription:        proto.CookingStepsDescription,
		Author:                         proto.Author,
		RequiredIngredientsProportions: make([]WeightedIngredientMongoDB, 0, len(proto.RequiredIngredientsProportions)),
		CookingTime:                    proto.CookingTime.Seconds,
		Macros100G:                     *mongoMacros,
	}

	for _, protoWeightedIngredient := range proto.RequiredIngredientsProportions {
		mongoWeightedIngredient, _ := protoWeightedIngredient.ConvertToMongoDocument()
		mongo.RequiredIngredientsProportions = append(mongo.RequiredIngredientsProportions, *mongoWeightedIngredient)
	}

	if len(proto.GetHexId()) != 0 {
		objectId, err := primitive.ObjectIDFromHex(proto.GetHexId())
		if err != nil {
			return nil, fmt.Errorf("Recipe.ConvertToMongoDocument returned error: %v", err)
		}
		mongo.Id = objectId
	}

	return mongo, nil
}

func (mongo *RecipeMongoDB) ConvertToProtoMessage() (*Recipe, error) {
	protoMacros, _ := mongo.Macros100G.ConvertToProtoMessage()
	proto := &Recipe{
		HexId:                          mongo.Id.Hex(),
		Title:                          mongo.Title,
		TasteDescription:               mongo.TasteDescription,
		CookingStepsDescription:        mongo.CookingStepsDescription,
		Author:                         mongo.Author,
		RequiredIngredientsProportions: make([]*WeightedIngredient, 0, len(mongo.RequiredIngredientsProportions)),
		CookingTime:                    &durationpb.Duration{Seconds: mongo.CookingTime},
		Macros100G:                     protoMacros,
	}

	for _, mongoWeightedIngredient := range mongo.RequiredIngredientsProportions {
		protoWeightedIngredient, _ := mongoWeightedIngredient.ConvertToProtoMessage()
		proto.RequiredIngredientsProportions = append(proto.RequiredIngredientsProportions, protoWeightedIngredient)
	}

	return proto, nil
}
