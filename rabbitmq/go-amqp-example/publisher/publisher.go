package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/emacsvi/dogolang/rabbitmq/go-amqp-example/contracts"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

var (
	amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func init() {
	flag.Parse()
	initAmqp()
}

var conn *amqp.Connection
var ch *amqp.Channel

func initAmqp() {
	var err error

	conn, err = amqp.Dial(*amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"go-test-exchange", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // noWait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare the Exchange")

	q, err := ch.QueueDeclare(
		"go-test-queue", // name, leave empty to generate a unique name
		true,            // durable
		false,           // delete when usused
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
	failOnError(err, "Error declaring the Queue")

	err = ch.QueueBind(
		q.Name,             // name of the queue
		"go-test-key",      // bindingKey
		"go-test-exchange", // sourceExchange
		false,              // noWait
		nil,                // arguments
	)
	failOnError(err, "Error binding to the Queue")
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func publishMessages(messages int) {
	for i := 0; i < messages; i++ {
		user := contracts.User{}
		user.FirstName = randomString(randInt(3, 10))
		user.LastName = randomString(randInt(3, 10))

		payload, err := json.Marshal(user)
		failOnError(err, "Failed to marshal JSON")

		err = ch.Publish(
			"go-test-exchange", // exchange
			"go-test-key",      // routing key
			false,              // mandatory
			false,              // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Transient,
				ContentType:  "application/json",
				Body:         payload,
				Timestamp:    time.Now(),
			})

		failOnError(err, "Failed to Publish on RabbitMQ")
	}
}

func main() {
	log.Println("Starting publisher...")

	// Close Channel
	defer ch.Close()

	// Close Connection
	defer conn.Close()

	for {
		// Publish messages
		publishMessages(10000)
		time.Sleep(500 * time.Millisecond)
		/*
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			if i == 4 {
				log.Println("start to resend message....")
			}
		}
		*/
	}

}
