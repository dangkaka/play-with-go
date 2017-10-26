## What is this project for?

* Understand know how kafka + Go + MongoDB work together
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

## Notes
* Bootstrap create 3 kafka nodes by default. If you want to add more, you can edit in bootstrap.sh

## Reference
* Use `https://github.com/wurstmeister/kafka-docker` as your Kafka cluster docker
