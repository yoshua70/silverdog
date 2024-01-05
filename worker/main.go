package main

import (
	"flag"
	"log"
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
}
