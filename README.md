## Ringier Golang Exercise

This project is a simple exercise of an event bus system. It consists of the following components:

- **Receiver Service** - receives events from a HTTP API and queues it to an MSB.
- **Consumer Service** - consumes events from the MSB and calls the dispatcher & persister services via a gRPC api.
- **Dispatcher Service** - simulates an event dispatcher by responding with a success rate of 80%.
- **Persister Service** - persists events to a mongodb database.
- **Monitoring Service** - provides a Prometheus metrics endpoint.

The following were the design choices I made when developing this project.

- MSB (RabbitMQ and Apache Kafka)
- Database (MongoDB)
- Service to Service calls (gRPC)
- Monitoring (Prometheus and Grafana)

## Requirements

The following requirements are necessary to run the project:

- Docker with Docker Compose

## Installation

Simply run the following command to install the project:

```bash
$ docker compose up --build
```

The receiver services exposes an HTTP API on localhost:9000/api/v1/events, send a POST request to it with the following
body:

```json
{
  "event": "UserCreated",
  "route": "general",
  "from": "47e9b5a0-c5af-4ae2-817d-9635dc1a3bec",
  "to": [
    {
      "to": "47e9b5a0-c5af-4ae2-817d-9635dc1a3bec",
      "endpoints": [
        "/api/users"
      ]
    }
  ],
  "reference": "some_reference",
  "created_at": "2022-08-19T11:02:00.511Z",
  "payload": {
    "id": "957ab08e-3188-4af0-be19-1e8e47cdfec9",
    "name": "John Doe",
    "email": "john.doe@example.com"
  }
}
```


