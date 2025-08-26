# Ethereum Block Scanner

High performance Go microservice for real-time monitoring of Ethereum addresses using native JSON-RPC. 

## Table of Contents

- [Ethereum Block Scanner](#ethereum-block-scanner)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
    - [Environment Variables](#environment-variables)
    - [Install Dependencies](#install-dependencies)
  - [Running Locally](#running-locally)
    - [With VS Code](#with-vs-code)
    - [Using binary](#using-binary)
    - [Running with Docker](#running-with-docker)

---

## Prerequisites

- Go `>= 1.24`
- Docker & Docker Compose
- Kafka

---

## Setup

### Environment Variables

Update a `.env` file available in the project root:

```env
ETH_NODE_URL=<WEBSOCKET_RPC_URL>
ADDRESSES_FILE=addresses.csv
BLOOM_FILTER_SIZE=10000000
BLOOM_FILTER_HASH=7
BATCH_SIZE=1000
CHECKPOINT_FILE=checkpoint.txt
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=ethereum-tx-events
```

### Install Dependencies
```
make install-deps
```

## Running Locally
### With VS Code
1. Open the project folder in VS Code.
2. Ensure .env is updated with your local configuration.
3. Build the binary:
```
make build
```
4. Launch application with VScode launch configuration

### Using binary
1. Build the binary:
```
make build
```
2. Run the binary
   Build the binary:
```
make run-app
```
This will:
- Start Kafka using Docker.
- Wait for Kafka to be ready.
- Launch the scanner binary.
  
### Running with Docker
```
make docker-build
make docker-start
```
This will:
- Build a Docker image for the scanner.
- Start the scanner container with Kafka dependency.

```
make docker-stop
```
This stops all Docker containers and cleans up.


**Notes**
1. .env file determines the configuration. Update Kafka brokers depending on whether you are running locally or inside Docker.
2. You can mount addresses.csv and .env in Docker using volumes.
3. Logs are printed to the console and events can be published to Kafka.
