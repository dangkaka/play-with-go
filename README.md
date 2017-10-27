## What is this project for?

* Understand how kafka + Go + MongoDB work together
* Use docker to build environment: easier to install, easier to use

## Setup

1. Install docker
2. Put your `go` application in `app/`
3. Start docker
* On Mac
    ```./bootstrap-mac.sh```
* On Ubuntu
    ```./bootstrap.sh```
4. `localhost:8080` is ready to use

## Usage
* Post some value
   ```curl -X POST http://localhost:8000?value=somevalue```
* Producer will push that value to `kafka`
* There are a consumer read kafka, process that value and push `processed` value to mongoDB
* Go to `http://localhost:8000/feed` to see all processed values

## MongoDB
* How to connect to mongoDB

    ```docker exec -it tinyfeed_mongo_1 mongo```

* Some basic commands
    - Show all db 
        `show dbs`
    - Use 1 db 
        `use {dbName}`
    - Show all collection 
        `show collections`
    - Find all data in 1 collection 
        `db.{collectionName}.find({})`
    
## Kafka
* How to connect to kafka: able to connect to any kafka node. Example: node1

    ```docker exec -it tinyfeed_kafka_1 bash```
    
    ```cd /opt/kafka/bin```

* Some basic commands
    - Listen to 1 topic. Example: feed
    
        ```./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic feed```
        
    - Open producer. Example: feed
    
        ```./kafka-console-producer.sh --broker-list localhost:9092 --topic feed```
        

## Notes
* Bootstrap create 3 kafka nodes by default. If you want to add more, you can edit in bootstrap.sh

## Reference
* Use `https://github.com/wurstmeister/kafka-docker` as Kafka cluster docker
