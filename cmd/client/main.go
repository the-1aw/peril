package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/the-1aw/peril/internal/gamelogic"
	"github.com/the-1aw/peril/internal/pubsub"
	"github.com/the-1aw/peril/internal/routing"
)

const amqpConnectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	conn, err := amqp091.Dial(amqpConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Starting Peril client...")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}
	pauseQueueName := routing.PauseKey + "." + username
	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, pauseQueueName, routing.PauseKey, pubsub.QueueTypeTransient)
	gameState := gamelogic.NewGameState(username)
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		switch userInput[0] {
		case "spawn":
			if err := gameState.CommandSpawn(userInput); err != nil {
				fmt.Fprintf(os.Stderr, "failed to spawn troup %v", err)
				continue
			}
		case "move":
			if _, err := gameState.CommandMove(userInput); err != nil {
				fmt.Fprintf(os.Stderr, "failed to move troup %v", err)
				continue
			}
		case "help":
			gamelogic.PrintClientHelp()
		case "status":
			gameState.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Printf("Unknown command %s\n", userInput[0])
		}
	}
}
