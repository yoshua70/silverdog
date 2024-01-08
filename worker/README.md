# TaskWorker

A worker to execute a dedicated task.

## Usage

Build the project into an executable:

```sh
go build -o worker *.go
```

You can launch as many workers as you want in multiple shells, for that open a new shell windows and run the following command:

```sh
./worker -name=[WORKER_NAME] -ruser=[RABBITMQ_USER] -rpwd=[RABBITMQ_PWD] -rhost=[RABBITMQ_HOSTNAME] -rport=[RABBITMQ_PORT]
```

The workers will then start receiving tasks from the backend server.

## Tasks

Workers can execute multiple tasks. For now the available tasks are the following:

- Download

### Download

A worker download a file from a given URL and store the downloaded files in a worker name `[WORKER_NAME]_output`. That is why it is important for the workers to have unique names.
