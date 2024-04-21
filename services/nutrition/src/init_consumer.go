package src

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func createConsumer() {

}

func InitConsumer(ctx context.Context) error {
	conn, err := amqp.Dial("amqp://guest:guest@message-broker:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	inventoryChannel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err)
	}
	defer inventoryChannel.Close()

	inventoryConsumerChannel, err := inventoryChannel.ConsumeWithContext(
		ctx,
		"ingredient",
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %s", err)
	}

	fmt.Println("inventory consumer initialized")

	for item := range inventoryConsumerChannel {
		log.Print(item)
	}
	return nil
}
