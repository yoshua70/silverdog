package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const DEFAULT_PORT int = 3333
const DEFAULT_RABBITMQ_USER string = "guest"
const DEFAULT_RABBITMQ_PWD string = "guest"
const DEFAULT_RABBITMQ_HOSTNAME string = "localhost"
const DEFAULT_RABBITMQ_PORT int = 5672

// "amqp://guest:guest@localhost:5672/"
var DEFAULT_RABBITMQ_URL string = fmt.Sprintf("amqp://%s:%s@%s:%d", DEFAULT_RABBITMQ_USER, DEFAULT_RABBITMQ_PWD, DEFAULT_RABBITMQ_HOSTNAME, DEFAULT_RABBITMQ_PORT)

const TASK_QUEUE_NAME string = "tasks"

var PORT int
var RABBITMQ_URL string
var RABBITMQ_USER string
var RABBITMQ_PWD string
var RABBITMQ_HOSTNAME string
var RABBITMQ_PORT int

type Task struct {
	Name     string `json:"name"`
	TaskType string `json:"taskType"`
	Arg      string `json:"arg"`
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Printf("%s /\n", r.Method)
	// An empty method corresponds to a GET response for a client.
	case "":
		log.Printf("%s /\n", r.Method)
	default:
		log.Printf("%s / unsupported method\n", r.Method)
		io.WriteString(w, "error")
		return
	}

	io.WriteString(w, "ok")
}

func HandleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		task, err := HandleTaskPostRequest(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		SendMessageToQeue(TASK_QUEUE_NAME, body)

		log.Printf("POST /task sent message\n")

		io.WriteString(w, "ok")
	default:
		log.Printf("%s / unsupported method\n", r.Method)
		io.WriteString(w, "error")
		return
	}
}

func HandleTaskPostRequest(body io.ReadCloser) (Task, error) {
	var task Task
	err := json.NewDecoder(body).Decode(&task)

	if err != nil {
		return task, errors.New(err.Error())
	}

	err = CheckTaskObject(task)

	if err != nil {
		return task, errors.New(err.Error())
	}

	log.Printf("POST /task %v\n", task)

	return task, nil
}

func HandleTaskGetRequest(w http.ResponseWriter, r *http.Request) {

}

func CheckTaskObject(task Task) error {
	if !(len(task.Name) > 0) || !(len(task.TaskType) > 0) || !(len(task.Arg) > 0) {
		return errors.New("some fields are missing. Please ensure that your request body contains the following fields: `name`, `taskType`, `arg`")
	}
	return nil
}

func main() {

	flag.IntVar(&PORT, "port", DEFAULT_PORT, "the listening port for the server")
	flag.StringVar(&RABBITMQ_USER, "ruser", DEFAULT_RABBITMQ_USER, "the user to connect to RabbitMQ with")
	flag.StringVar(&RABBITMQ_PWD, "rpwd", DEFAULT_RABBITMQ_PWD, "the password for the RabbitMQ user")
	flag.StringVar(&RABBITMQ_HOSTNAME, "rhost", DEFAULT_RABBITMQ_HOSTNAME, "the hostname of the RabbitMQ instance")
	flag.IntVar(&RABBITMQ_PORT, "rport", DEFAULT_RABBITMQ_PORT, "the port of the RabbitMQ instance")

	flag.Parse()

	log.Printf("argument `port` set to %d\n", PORT)
	log.Printf("argument `ruser` set to %s\n", RABBITMQ_USER)
	log.Printf("argument `rpwd` set to %s\n", RABBITMQ_PWD)
	log.Printf("argument `rhost` set to %s\n", RABBITMQ_HOSTNAME)
	log.Printf("argument `rport` set to %d\n", RABBITMQ_PORT)

	RABBITMQ_URL = fmt.Sprintf("amqp://%s:%s@%s:%d", RABBITMQ_USER, RABBITMQ_PWD, RABBITMQ_HOSTNAME, RABBITMQ_PORT)
	log.Printf("rabbitmq connection url is: %s\n", RABBITMQ_URL)

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/task", HandleTask)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
