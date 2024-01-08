# Middleman

A tiny middleman for RabbitMQ written in Go.

## Usage

Build the project into an executable:

```sh
go build -o middleman *.go
```

Launch the middleman:

```sh
./middleman -ruser=[RABBITMQ_USER] -rpwd=[RABBITMQ_PWD] -rhost=[RABBITMQ_HOSTNAME] -rport=[RABBITMQ_PORT]
```
