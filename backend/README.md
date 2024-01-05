# A simple task queue in Go using RabbitMQ

## Usage

Build the project into an executable:
```sh
go build -o tasker *.go 
```

Launch the http server:
```sh
./tasker
```

You can now make requests to the server.

## Run RabbitMQ

You should have docker installed on your machine.

```sh
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
```