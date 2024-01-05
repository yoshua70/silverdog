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
const DEFAULT_RABBITMQ_URL string = "amqp://guest:guest@localhost:5672/"
const TASK_QUEUE_NAME string = "tasks"

var PORT int
var RABBITMQ_URL string

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
	flag.StringVar(&RABBITMQ_URL, "rabbitmq", DEFAULT_RABBITMQ_URL, "the connection url to RabbitMQ")

	flag.Parse()

	log.Printf("argument `port` set to %d\n", PORT)
	log.Printf("argument `rabbitmq` set to %s\n", RABBITMQ_URL)

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
