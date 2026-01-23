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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(state routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(state)
	}
}

func handleMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(move gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(move)
	}
}

func main() {
	conn, err := amqp091.Dial(amqpConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
	}
	fmt.Println("Starting Peril client...")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}
	pauseQueueName := routing.PauseKey + "." + username
	armyMovesQueueName := routing.ArmyMovesPrefix + "." + username
	gameState := gamelogic.NewGameState(username)
	pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, pauseQueueName, routing.PauseKey, pubsub.QueueTypeTransient, handlerPause(gameState))
	pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, armyMovesQueueName, routing.ArmyMovesPrefix+".*", pubsub.QueueTypeTransient, handleMove(gameState))
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
			if move, err := gameState.CommandMove(userInput); err != nil {
				fmt.Fprintf(os.Stderr, "failed to move troup %v", err)
				continue
			} else {
				pubsub.PublishJSON(channel, routing.ExchangePerilTopic, armyMovesQueueName, move)
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
