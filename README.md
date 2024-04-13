This project is a template project with very basic business logic implementing Hexagonal Architecture in Go Lang

On this documentation I will not try to explain what Hexagonal Architecture is. Since that is not the target of this repository. Instead, the idea is to demonstrate how to implement it. In this case, using Go Lang. But the same idea can be implemented for any other programming language.

## Getting Started 

- clone the project 
- copy the `.env.example` file and create a `.env` file: `cp .env.example .env`
- execute `docker compose up`

This docker compose will start several containers containing an entire ecosystem to manage queues, connect to a database, monitor the applications, get logs and observability, and there it goes. So at first it is not a good idea trying to understand each of these tools because it may be a bit overwhelming.

## Stack 

This project contains an extensive stack that for sure is over engineering and adding more dependencies than it would be necessary for a production service. But the ideia here was to give and example that it is possible to add all of those components and yet, keep the service organized and working.

Here is the list of the current components:

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)
* [Go Lang](https://go.dev/)
* [Air - Live reload for Go apps](https://github.com/cosmtrek/air)
* [Postgres](https://www.postgresql.org/)
* [Apache Kafka](https://kafka.apache.org/)
* [Jaeger](https://www.jaegertracing.io/)
* [Open Telemetry](https://opentelemetry.io/)
* [Prometheus](https://prometheus.io/)
* [Grafana](https://grafana.com/)
* [Rabbit MQ](https://www.rabbitmq.com/)
* [gRPC](https://grpc.io/)
* [Protocol Buffer](https://protobuf.dev/)

## Project structure 

The initial thing to understand about the project structure is that the root folder contains 3 folders that are strictly related to Go Lang ecosystem:

- `cmd`: Contains the possible entrypoints for the applications. Specially useful if we are creating multiple services that has different entrypoints but shares the same business logic. 
- `pkg`: Contains utilities, separated by domains, that are not implementing a specific business logic. Usually the domains on this folders can be simply copied and pasted on other projects as if it was an external package. It can be seen as future external packages.
- `internal`: Contains the Bussiness logic of the application. Specific things that are the reason for this application to exist.

The remaining folders are not the focus now. They do exist for a purpose but we will reach them step by step.

> Disclaimer: The idea of this explanation is to try to be as didatic as possible. So it may not be extremely accurate in terms of theory.

