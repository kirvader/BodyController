package src

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/kirvader/BodyController/services/nutrition/proto"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (svc *Service) CreateRecipe(ctx context.Context, req *pb.CreateRecipeRequest) (*pb.CreateRecipeResponse, error) {
	if req == nil || req.GetEntity() == nil || req.GetEntity().GetId() == "" { // TODO add real validation
		return nil, errors.New("nil instance")
	}

	body, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = svc.workerChannelMQ.PublishWithContext(ctx,
		"",       // exchange
		"recipe", // routing key
		false,    // mandatory
		false,    // immediate
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

	return &pb.CreateRecipeResponse{
		EntityId: req.GetEntity().GetId(),
	}, nil
}

func (svc *Service) DeleteRecipe(ctx context.Context, req *pb.DeleteRecipeRequest) (*pb.DeleteRecipeResponse, error) {
	if req == nil || req.EntityId == "" { // TODO add real validation
		return nil, errors.New("nil instance")
	}

	body, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = svc.workerChannelMQ.PublishWithContext(ctx,
		"",       // exchange
		"recipe", // routing key
		false,    // mandatory
		false,    // immediate
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

	return &pb.DeleteRecipeResponse{}, nil
}

func (svc *Service) ListRecipes(ctx context.Context, req *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
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

	coll := svc.mongoClient.Database("BodyController").Collection("Recipes")

	options := options.Find()
	options.SetSort(bson.M{"name": 1})
	options.SetSkip(int64(pageOffset))
	options.SetLimit(int64(pageSize))

	cursor, err := coll.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*pb.Recipe, 0, req.GetPageSize())

	for cursor.Next(ctx) {
		var mongoInstance pb.RecipeMongo
		err := cursor.Decode(&mongoInstance)
		if err != nil {
			return nil, fmt.Errorf("cursor decode error: %v", err)
		}

		instance, err := mongoInstance.Proto()
		if err != nil {
			return nil, fmt.Errorf("failed to parse entity: %v", err)
		}

		result = append(result, instance)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	if int32(len(result)) < pageSize {
		return &pb.ListRecipesResponse{
			Entities: result,
		}, nil
	}

	nextPageToken, err := constructPageToken(pageSize, pageOffset)
	if err != nil {
		return nil, err
	}

	return &pb.ListRecipesResponse{
		Entities: result,
		NextPageToken: &wrapperspb.StringValue{
			Value: nextPageToken,
		},
	}, nil
}
