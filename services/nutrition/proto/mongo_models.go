package proto

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var ErrNilInstance = errors.New("nil instance")
var ErrInvalidId = errors.New("invalid id")

type MacrosMongo struct {
	Proteins float32 `bson:"proteins"`
	Carbs    float32 `bson:"carbs"`
	Fats     float32 `bson:"fats"`
	Calories float32 `bson:"calories"`
}

func (instance *Macros) Mongo() (*MacrosMongo, error) {
	return &MacrosMongo{
		Proteins: float32(instance.GetProteins()),
		Carbs:    float32(instance.GetCarbs()),
		Fats:     float32(instance.GetFats()),
		Calories: float32(instance.GetCalories()),
	}, nil
}

func (instance *MacrosMongo) Proto() (*Macros, error) {
	if instance == nil {
		return nil, ErrNilInstance
	}
	return &Macros{
		Calories: instance.Calories,
		Proteins: instance.Proteins,
		Carbs:    instance.Carbs,
		Fats:     instance.Fats,
	}, nil
}

type IngredientMongo struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	Title            string             `bson:"title"`
	MacrosNormalized *MacrosMongo       `bson:"macros_normalized"`
}

func (instance *Ingredient) Mongo() (*IngredientMongo, error) {
	macrosNormalized, err := instance.GetMacrosNormalized().Mongo()
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(instance.GetId())
	if err != nil {
		return nil, ErrInvalidId
	}

	return &IngredientMongo{
		Id:               id,
		Title:            instance.GetTitle(),
		MacrosNormalized: macrosNormalized,
	}, nil
}

func (instance *IngredientMongo) Proto() (*Ingredient, error) {
	if instance == nil {
		return nil, ErrNilInstance
	}

	macrosForWeight, err := instance.MacrosNormalized.Proto()
	if err != nil {
		return nil, err
	}

	return &Ingredient{
		Id:               instance.Id.Hex(),
		Title:            instance.Title,
		MacrosNormalized: macrosForWeight,
	}, nil
}

type WeightedIngredientMongo struct {
	IngredientId primitive.ObjectID `bson:"ingredient_id"`
	Gramms       float32            `bson:"gramms"`
}

func (instance *WeightedIngredient) Mongo() (*WeightedIngredientMongo, error) {
	ingredientId, err := primitive.ObjectIDFromHex(instance.GetIngredientId())
	if err != nil {
		return nil, ErrInvalidId
	}

	return &WeightedIngredientMongo{
		IngredientId: ingredientId,
		Gramms:       instance.GetGramms(),
	}, nil
}

func (instance *WeightedIngredientMongo) Proto() (*WeightedIngredient, error) {
	if instance == nil {
		return nil, ErrNilInstance
	}

	return &WeightedIngredient{
		IngredientId: instance.IngredientId.Hex(),
		Gramms:       instance.Gramms,
	}, nil
}

type RecipeMongo struct {
	Id            primitive.ObjectID   `bson:"_id,omitempty"`
	Title         string               `bson:"title"`
	Description   string               `bson:"description"`
	IngredientIds []primitive.ObjectID `bson:"ingredient_ids"`
}

func (instance *Recipe) Mongo() (*RecipeMongo, error) {
	id, err := primitive.ObjectIDFromHex(instance.GetId())
	if err != nil {
		return nil, ErrInvalidId
	}

	ingredientIds := make([]primitive.ObjectID, 0, len(instance.GetIngredientIds()))
	for _, ingredientHexId := range instance.GetIngredientIds() {
		ingredientId, err := primitive.ObjectIDFromHex(ingredientHexId)
		if err != nil {
			return nil, ErrInvalidId
		}
		ingredientIds = append(ingredientIds, ingredientId)
	}

	return &RecipeMongo{
		Id:            id,
		Title:         instance.GetTitle(),
		Description:   instance.GetDescription(),
		IngredientIds: ingredientIds,
	}, nil
}

func (instance *RecipeMongo) Proto() (*Recipe, error) {
	if instance == nil {
		return nil, ErrNilInstance
	}

	ingredientIds := make([]string, 0, len(instance.IngredientIds))
	for _, mongoIngredientId := range instance.IngredientIds {
		ingredientIds = append(ingredientIds, mongoIngredientId.Hex())
	}

	return &Recipe{
		Id:            instance.Id.Hex(),
		Title:         instance.Title,
		Description:   instance.Description,
		IngredientIds: ingredientIds,
	}, nil
}

type MealMongo struct {
	Id                  primitive.ObjectID         `bson:"_id,omitempty"`
	Username            string                     `bson:"username"`
	RecipeId            primitive.ObjectID         `bson:"recipe_id,omitempty"`
	WeightedIngredients []*WeightedIngredientMongo `bson:"weighted_ingredients"`
	Timestamp           int64                      `bson:"timestamp"`
	MealStatus          uint8                      `bson:"meal_status"`
}

func (instance *Meal) Mongo() (*MealMongo, error) {
	id, err := primitive.ObjectIDFromHex(instance.GetId())
	if err != nil {
		return nil, ErrInvalidId
	}

	recipeId, err := primitive.ObjectIDFromHex(instance.GetRecipeId().GetValue())
	if err != nil {
		return nil, ErrInvalidId
	}

	ingredients := make([]*WeightedIngredientMongo, 0, len(instance.GetWeightedIngredients()))
	for _, weightedIngredient := range instance.GetWeightedIngredients() {
		weightedIngredientMongo, err := weightedIngredient.Mongo()
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, weightedIngredientMongo)
	}

	return &MealMongo{
		Id:                  id,
		Username:            instance.GetUsername(),
		RecipeId:            recipeId,
		WeightedIngredients: ingredients,
		Timestamp:           instance.GetTimestamp().GetSeconds(),
		MealStatus:          uint8(instance.GetMealStatus().Number()),
	}, nil
}

func (instance *MealMongo) Proto() (*Meal, error) {
	if instance == nil {
		return nil, ErrNilInstance
	}

	ingredients := make([]*WeightedIngredient, 0, len(instance.WeightedIngredients))
	for _, weightedIngredientMongo := range instance.WeightedIngredients {
		weightedIngredient, err := weightedIngredientMongo.Proto()
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, weightedIngredient)
	}

	return &Meal{
		Id:       instance.Id.Hex(),
		Username: instance.Username,
		RecipeId: &wrapperspb.StringValue{
			Value: instance.RecipeId.Hex(),
		},
		WeightedIngredients: ingredients,
		Timestamp: &timestamppb.Timestamp{
			Seconds: instance.Timestamp,
		},
		MealStatus: MealStatus(instance.MealStatus),
	}, nil
}
