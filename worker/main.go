package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO: get name from the backend server
const DEFAULT_NAME string = "worker"
const DEFAULT_RABBITMQ_URL string = "amqp://guest:guest@localhost:5672/"
const TASK_QUEUE_NAME string = "tasks"
const NOTIF_QUEUE_NAME string = "notifications"

var OUTPUT_DIR string = ""
var WORKER_NAME string
var RABBITMQ_URL string

type Task struct {
	Name     string `json:"name"`
	TaskType string `json:"taskType"`
	Arg      string `json:"arg"`
}

// Download data from the specified URL.
// The downloaded data is stored inside of the OUTPUT_DIR directory.
func Downloader(fullURLFile string) error {
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Printf("error while downloading file: %s\n", err)
		return err
	}

	filePath := fileURL.Path
	segments := strings.Split(filePath, "/")
	var fileName = segments[len(segments)-1]

	file, err := os.Create(fmt.Sprintf("./%s/%s", OUTPUT_DIR, fileName))
	if err != nil {
		log.Printf("error while downloading file: %s\n", err)
		return err
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Printf("error while downloading file: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	log.Printf("[Downloader]  downloaded file %s with size %d", fileName, size)

	return nil
}

func main() {
	flag.StringVar(&WORKER_NAME, "name", DEFAULT_NAME, "the name of the worker, must be unique")
	flag.StringVar(&RABBITMQ_URL, "rabbitmq", DEFAULT_RABBITMQ_URL, "the connection url to RabbitMQ")

	flag.Parse()

	log.Printf("argument `name` set to %s\n", WORKER_NAME)
	log.Printf("argument `rabbitmq` set to %s\n", RABBITMQ_URL)

	// Create the output directory for downloaded files.
	OUTPUT_DIR = fmt.Sprintf("%s_output", WORKER_NAME)
	err := os.Mkdir(OUTPUT_DIR, 0755)
	defer os.RemoveAll(OUTPUT_DIR)

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
	messages, err := ch.Consume(
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
		for message := range messages {
			log.Printf("[%s] received a message: %s\n", WORKER_NAME, message.Body)

			var task Task

			err := json.Unmarshal(message.Body, &task)

			if err != nil {
				log.Printf("[%s] failed to decode message to task: %s\n", WORKER_NAME, err)
			} else {
				log.Printf("[%s] received task `%s` of type `%s` with args `%s`\n", WORKER_NAME, task.Name, task.TaskType, task.Arg)

				Downloader(task.Arg)
			}
		}
	}()

	log.Printf("[%s] waiting for messages. To exit press CTRL+C\n", WORKER_NAME)
	<-forever
}
