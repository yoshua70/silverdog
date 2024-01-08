# A simple task queue in Go using RabbitMQ

## Usage

Build the project into an executable:

```sh
go build -o tasker *.go
```

Launch the http server:

```sh
./tasker -ruser=[RABBITMQ_USER] -rpwd=[RABBITMQ_PWD] -rhost=[RABBITMQ_HOSTNAME] -rport=[RABBITMQ_PORT]
```

You can now make requests to the server.

## Run RabbitMQ

You should have docker installed on your machine.

```sh
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
```

## Endpoints

The server exposes different endpoints:

- `/`
- `/task`

### `/`

#### `GET`

The root of the API. It just returns a `ok` response.

### `/task`

#### `POST`

Make a post request to the route `/task` to create a new task. The body of the request should be of the following form:

```json
{
  "name": "[the name of task]",
  "taskType": "[the task to be executed]",
  "arg": "[the needed arguments for the task]"
}
```

The available `taskType` are the following:

- download: download a resource from an URL. For that you must give the url in the `arg` key of the body.
