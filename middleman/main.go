package main

import (
	"flag"
	"log"
)

const DEFAULT_PORT int = 8090
const DEFAULT_RABBITMQ_URL string = "amqp://guest:guest@localhost:5672/"
const TASK_QUEUE_NAME string = "tasks"
const NOTIF_QUEUE_NAME string = "notifications"

var PORT int
var RABBITMQ_URL string

func main() {

	flag.IntVar(&PORT, "port", DEFAULT_PORT, "the listening port for the server")
	flag.StringVar(&RABBITMQ_URL, "rabbitmq", DEFAULT_RABBITMQ_URL, "the connection url to RabbitMQ")

	flag.Parse()

	log.Printf("argument `port` set to %d\n", PORT)
	log.Printf("argument `rabbitmq` set to %s\n", RABBITMQ_URL)
}
