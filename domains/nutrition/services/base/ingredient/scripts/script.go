package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kirvader/BodyController/domains/nutrition/models"
	pbIngredient "github.com/kirvader/BodyController/domains/nutrition/services/base/ingredient/proto"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:10000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pbIngredient.NewIngredientClient(conn)

	resp, err := client.CreateIngredient(context.Background(), &pbIngredient.CreateIngredientRequest{
		Ingredient: &models.Ingredient{
			Macros100G: &models.Macros100G{
				Proteins: 40,
				Carbs:    40,
				Fats:     20,
				Calories: 400,
			},
			Title:       "ingredient",
			Description: "good one",
		},
	})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("response got: %v", resp)
}
