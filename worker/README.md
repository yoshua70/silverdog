# TaskWorker

A worker to execute a dedicated task.

## Usage

Build the project into an executable:
```sh
go build -o worker *.go 
```

You can launch as many workers as you want in multiple shells:
```sh
./worker -name=[WORKER_NAME] -rabbitmq=[RABBITMQ_CONNECTION_URL]
```

The workers will then start receiving tasks from the backend server.