package src

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/kirvader/BodyController/services/nutrition/proto"
)

func (svc *Service) CreateMeal(ctx context.Context, req *pb.CreateMealRequest) (*pb.CreateMealResponse, error) {
	if req == nil || req.GetEntity() == nil || req.GetEntity().GetId() == "" { // TODO add real validation
		return nil, errors.New("nil instance")
	}

	body, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = svc.workerChannelMQ.PublishWithContext(ctx,
		"",     // exchange
		"meal", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Type:        OperationCreate,
			Timestamp:   time.Now(),
			Body:        []byte(body),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to publish event: %s", err)
	}
	log.Println("published CREATE event with id: ", req.GetEntity().GetId())

	return &pb.CreateMealResponse{
		EntityId: req.GetEntity().GetId(),
	}, nil

	// coll := svc.mongoClient.Database("BodyController").Collection("Meals")

	// ingredients := make([]weightedIngredientMongo, 0, len(req.GetEntity().GetWeightedIngredients()))
	// for _, item := range req.GetEntity().GetWeightedIngredients() {
	// 	ingredients = append(ingredients, convertweightedIngredientProtoToMongo(item))
	// }
	// var recipeId primitive.ObjectID

	// if req.GetEntity().GetRecipeId().GetValue() == "" {
	// 	recipeId = primitive.NilObjectID
	// } else {
	// 	parsedRecipeId, err := primitive.ObjectIDFromHex(req.GetEntity().GetRecipeId().GetValue())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	recipeId = parsedRecipeId
	// }

	// resp, err := coll.InsertOne(ctx, mealMongo{
	// 	Username:            req.GetEntity().GetUsername(),
	// 	RecipeId:            recipeId,
	// 	WeightedIngredients: ingredients,
	// 	Timestamp:           req.GetEntity().GetTimestamp().GetSeconds(),
	// 	MealStatus:          uint8(req.GetEntity().GetMealStatus()),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// return &pb.CreateMealResponse{
	// 	EntityId: resp.InsertedID.(primitive.ObjectID).Hex(),
	// }, nil
}

// DeleteMeal implements proto.NutritionServer.
func (svc *Service) DeleteMeal(ctx context.Context, req *pb.DeleteMealRequest) (*pb.DeleteMealResponse, error) {
	if req == nil || req.EntityId == "" { // TODO add real validation
		return nil, errors.New("nil instance")
	}

	body, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = svc.workerChannelMQ.PublishWithContext(ctx,
		"",     // exchange
		"meal", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Type:        OperationDelete,
			Timestamp:   time.Now(),
			Body:        []byte(body),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to publish event: %s", err)
	}
	log.Println("published DELETE event with id: ", req.EntityId)

	return &pb.DeleteMealResponse{}, nil
	// coll := svc.mongoClient.Database("BodyController").Collection("Meals")

	// objectId, err := primitive.ObjectIDFromHex(req.GetEntityId())
	// if err != nil {
	// 	return nil, err
	// }

	// resp, err := coll.DeleteOne(ctx, bson.M{"_id": objectId})
	// if err != nil {
	// 	return nil, err
	// }

	// if resp.DeletedCount != 1 {
	// 	return nil, errors.New("delete operation failed")
	// }
	// return &pb.DeleteMealResponse{}, nil
}

// ListMeals implements proto.NutritionServer.
func (svc *Service) ListMeals(ctx context.Context, req *pb.ListMealsRequest) (*pb.ListMealsResponse, error) {
	var pageSize, pageOffset int32
	if req.GetPageToken() != nil {
		pageSizeFromToken, pageOffsetFromToken, err := deconstructPageToken(req.GetPageToken().GetValue())
		if err != nil {
			return nil, err
		}
		pageSize, pageOffset = pageSizeFromToken, pageOffsetFromToken
	} else {
		pageOffset = 0
		if req.GetPageSize() <= 0 {
			pageSize = 100
		} else if req.GetPageSize() >= 500 {
			pageSize = 500
		} else {
			pageSize = req.GetPageSize()
		}
	}

	coll := svc.mongoClient.Database("BodyController").Collection("Meals")

	options := options.Find()
	options.SetSort(bson.M{"username": 1, "timestamp": -1})
	options.SetSkip(int64(pageOffset))
	options.SetLimit(int64(pageSize))

	cursor, err := coll.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*pb.Meal, 0, req.GetPageSize())

	for cursor.Next(ctx) {
		var mongoInstance pb.MealMongo
		err := cursor.Decode(&mongoInstance)
		if err != nil {
			return nil, fmt.Errorf("cursor decode error: %v", err)
		}
		instance, err := mongoInstance.Proto()
		if err != nil {
			return nil, fmt.Errorf("failed to parse mongo record: %v", err)
		}

		result = append(result, instance)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}
	if int32(len(result)) < pageSize {
		return &pb.ListMealsResponse{
			Entities: result,
		}, nil
	}

	nextPageToken, err := constructPageToken(pageSize, pageOffset)
	if err != nil {
		return nil, err
	}

	return &pb.ListMealsResponse{
		Entities: result,
		NextPageToken: &wrapperspb.StringValue{
			Value: nextPageToken,
		},
	}, nil
}
