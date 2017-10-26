package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

var kafkaHost []string = []string{"kafka:9092"}
var kafkaTopic = "feed"
var mongoHost = []string{
	"mongo:27017",
}
type Feed struct {
	Value string `bson:"value"`
}
const (
	Database   = "testDB"
	Collection = "feed"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", writeHandler).Methods("POST")
	r.HandleFunc("/feed", viewFeed).Methods("GET")

	go readHandler()

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func viewFeed(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs : mongoHost,
		Database: Database,
	})
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(Database).C(Collection)

	//get all
	var result []Feed
	err = c.Find(nil).Sort("-_id").All(&result)
	if err != nil {
		panic(err)
	}

	var response []byte
	json, err := json.Marshal(result)
	if err != nil {
		response = []byte("{\"status\":\"not ok\",\"message\":\"unable to marshal json\"}")
	}

	response = []byte(json)
	w.Write(response)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	var value string
	if val, ok := r.URL.Query()["value"]; ok {
		value = val[0]
	} else {
		w.Write(respond(map[string]string{
			"status":  "not ok",
			"message": "parameter `value` is missing",
		}))
		return
	}

	k := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kafkaHost,
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	defer k.Close()

	k.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: []byte(value),
		},
	)

	w.Write(respond(map[string]string{
		"status": "ok",
		"value":  value,
	}))
}

func readHandler() {
	k := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   kafkaHost,
		Topic:     kafkaTopic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	defer k.Close()

	//write to mongoDB
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs : mongoHost,
		Database: Database,
	})
	if err != nil {
		panic(err)
	}
	defer session.Close()

	for {
		m, err := k.ReadMessage(context.Background())
		if err != nil {
			break
		}

		fmt.Printf(
			"[CONSUMER] message at offset %d = %s\n",
			m.Offset,
			string(m.Value),
		)

		feed := Feed{
			Value: "Processed " + string(m.Value),
		}
		c := session.DB(Database).C(Collection)

		// Insert
		if err := c.Insert(feed); err != nil {
			panic(err)
		}
	}
}

func respond(o map[string]string) []byte {
	json, err := json.Marshal(o)
	if err != nil {
		return []byte("{\"status\":\"not ok\",\"message\":\"unable to marshal json\"}")
	}

	return []byte(json)
}
