package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/the-1aw/peril/internal/gamelogic"
	"github.com/the-1aw/peril/internal/pubsub"
	"github.com/the-1aw/peril/internal/routing"
)

const amqpConnectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	amqpConnection, err := amqp.Dial(amqpConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()
	channel, err := amqpConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = pubsub.DeclareAndBind(amqpConnection, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.QueueTypeDurable)
	if err != nil {
		log.Fatal(err)
	}
	gamelogic.PrintServerHelp()
	for {
		userInputs := gamelogic.GetInput()
		if len(userInputs) == 0 {
			continue
		}
		switch userInputs[0] {
		case "pause":
			fmt.Println("Pausing the game")
			pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		case "resume":
			fmt.Println("Resuming the game")
			pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		case "quit":
			fmt.Println("Quitting the game")
			return
		default:
			fmt.Println("Unknown command: ", userInputs[0])
		}
	}
}
