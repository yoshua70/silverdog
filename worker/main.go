package main

import (
	"flag"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO: get name from the backend server
const DEFAULT_NAME string = "worker"
const DEFAULT_RABBITMQ_URL string = "amqp://guest:guest@localhost:5672/"
const TASK_QUEUE_NAME string = "tasks"
const NOTIF_QUEUE_NAME string = "notifications"

var WORKER_NAME string
var RABBITMQ_URL string

func main() {
	flag.StringVar(&WORKER_NAME, "name", DEFAULT_NAME, "the name of the worker, must be unique")
	flag.StringVar(&RABBITMQ_URL, "rabbitmq", DEFAULT_RABBITMQ_URL, "the connection url to RabbitMQ")

	flag.Parse()

	log.Printf("argument `name` set to %s\n", WORKER_NAME)
	log.Printf("argument `rabbitmq` set to %s\n", RABBITMQ_URL)

	// TODO: listen to messages from the queue.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		TASK_QUEUE_NAME, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	tasks, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "failed to register a consumer")

	var forever chan struct{}

	go func() {
		for task := range tasks {
			log.Printf("[%s] received a message: %s\n", WORKER_NAME, task.Body)
		}
	}()

	log.Printf("[%s] waiting for messages. To exit press CTRL+C\n", WORKER_NAME)
	<-forever
}
