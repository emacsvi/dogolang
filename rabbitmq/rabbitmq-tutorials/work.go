package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"time"
)


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	// https://blog.csdn.net/weixin_34085658/article/details/85842158
	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to set Qos")
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
