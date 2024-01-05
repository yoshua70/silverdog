package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)


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

func main() {
	http.HandleFunc("/", HandleRoot)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
