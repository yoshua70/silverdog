package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

const REFRESH_INTERVAL int = 5
const DEFAULT_PORT int = 8090
const DEFAULT_RABBITMQ_URL string = "amqp://guest:guest@localhost:5672/"
const TASK_QUEUE_NAME string = "tasks"
const NOTIF_QUEUE_NAME string = "notifications"

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
	Sent   bool   `json:"sent"`
	Name   string `json:"name"`
}

var PORT int
var RABBITMQ_URL string
var messages []Message

// Upgrade a regular HTTP connection to a websocket connection.
var httpToWebSocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func consummer() {
	conn, err := amqp.Dial(RABBITMQ_URL)
	FailOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		NOTIF_QUEUE_NAME, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	FailOnError(err, "failed to declare a queue")

	messages, err := ch.Consume(
		NOTIF_QUEUE_NAME,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("failed to register consummer for queue %s: %s\n", NOTIF_QUEUE_NAME, err)
		return
	}

	for message := range messages {
		storeMessage(message)
	}
}

func storeMessage(message amqp.Delivery) {
	msg, err := messageParser(string(message.Body))

	if err == nil {
		log.Printf("store message: %v\n", msg)
		messages = append(messages, msg)
	}
}

func messageParser(message string) (Message, error) {
	msg := Message{Sent: false}
	err := json.Unmarshal([]byte(message), &msg)

	if err != nil {
		log.Printf("failed to decode message: %s\n", err)
		return Message{}, err
	}

	return msg, nil
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := httpToWebSocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("failed to upgrade http connection to websocket: %s\n", err)
		return
	}

	defer conn.Close()

	for {
		time.Sleep(time.Second * time.Duration(REFRESH_INTERVAL))

		for i, message := range messages {
			if !message.Sent {
				err := conn.WriteJSON(message)

				if err != nil {
					log.Printf("failed to write json in websocket pipe: %s\n", err)
					return
				}
				log.Printf("sent message: %s\n", message.Body)
				messages[i] = Message{Body: message.Body, Status: message.Status, Sent: true, Name: message.Name}
			}
		}
	}
}

func main() {

	flag.IntVar(&PORT, "port", DEFAULT_PORT, "the listening port for the server")
	flag.StringVar(&RABBITMQ_URL, "rabbitmq", DEFAULT_RABBITMQ_URL, "the connection url to RabbitMQ")

	flag.Parse()

	log.Printf("argument `port` set to %d\n", PORT)
	log.Printf("argument `rabbitmq` set to %s\n", RABBITMQ_URL)

	go consummer()

	http.HandleFunc("/ws", handleWebSocket)
	log.Printf("websoket running on: %d\n", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
