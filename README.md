# Messaging Campaign System

This project is a scalable messaging campaign management system designed to handle large-scale messaging using Kafka for message queuing and Go for backend services. The system allows users to create and schedule campaigns, manage recipient lists, and send personalized messages in bulk.

## Table of Contents

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation & Running](#installation--running)
- [Configuration](#Configuration)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Tolist](#todolist)



## Introduction

This system is built to support the creation, scheduling, and broadcasting of messaging campaigns using Kafka for queuing and Go for backend processing. The platform enables large-scale message delivery to recipients from a CSV file containing their phone numbers and names, with customizable message templates.


## Features

List the main features of the project:

- Create and schedule messaging campaigns
- Import recipient lists from CSV files
- Generate and queue personalized messages for delivery
- Kafka-based message queuing for high throughput
- Consumer service for processing and simulating message delivery
- MySQL for campaign and recipient data storage

## Installation & Running

### Prerequisites

Make sure the following tools are installed:

- [Go 1.23.1](https://golang.org/doc/install)
- [MySQL 8.0+](https://dev.mysql.com/downloads/mysql/)
- [Kafka 2.8.0+](https://kafka.apache.org/downloads) (with Zookeeper)

### Running Locally

1. Install and start mysql:
    ```bash
    brew install mysql
    brew services start mysql
    ```
2. Install and start kafka:
    ```bash
    brew install kafka
    brew services start kafka   

3. Install mysql scripts:
    ```bash
    cd camp-mgr/app/deploy
    ./install_camp_db.sh
    ```
   
4. Install dependencies:
    ```bash
    go mod tidy
    ```   

5. Run the service as a producer only:
    ```bash
    cd /camp-mgr/app/campmgr/cmd
    go run main.go -c ../etc/campmgr.yaml
    ```
6. Run the service as producer and consumer:
    ```bash
    cd /camp-mgr/app/campmgr/cmd
    go run main.go -c ../etc/campmgr.yaml -w
    ```

The service will be running at `http://localhost:10001`.

## Configuration

Refer to camp-mgr/app/campmgr/etc/campmgr.yaml

## API Documentation
todo
### Example

#### Create campaign

- **URL**: `http://127.0.0.1:10001/v1/camp/create`
- **Method**: `POST`
- **Request Body**:
    ```json
    {
    "campaign_name": "game0",
    "campaign_id": 1000000000,
    "message_template": "hi, this is a test message.",
    "scheduled_time": 1727351999,
    "csv_file_path": "recipients.csv"
    }
    ```
- **Response**:
    ```json
    {
    "code": 0,
    "msg": "success"
    }
    ```

## Tolist
1. Unit test code 
2. Support API document, eg. go-swagger
3. Consumer Support for Cluster Deployment
4. Optimize MySQL queries for better efficiency
5. Use memory pools to optimize frequent memory allocations
6. Use a goroutine pool to optimize frequent task starts

