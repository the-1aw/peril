package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

const amqpConnectionString = "amqp://guest:guest@localhost:5672/"

func hangUntilInterrupt() {
	terminationChan := make(chan os.Signal, 1)
	signal.Notify(terminationChan, os.Interrupt)
	<-terminationChan
}

func main() {
	amqpConnection, err := amqp.Dial(amqpConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()
	fmt.Println("Connection to amqp was successful...")
	hangUntilInterrupt()
	fmt.Println("\nShutting down amqp server...")
}
