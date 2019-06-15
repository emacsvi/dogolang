package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
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

	body := bodyFrom(os.Args)
	// 消息的持久化
	// 如过将queue的持久化标识durable设置为true,则代表是一个持久的队列，那么在服务重启之后，也会存在，因为服务会把持久化的queue存放在硬盘上，当服务重启的时候，会重新什么之前被持久化的queue。队列是可以被持久化，但是里面的消息是否为持久化那还要看消息的持久化设置。也就是说，重启之前那个queue里面还没有发出去的消息的话，重启之后那队列里面是不是还存在原来的消息，这个就要取决于发生着在发送消息时对消息的设置了。
	// 如果要在重启后保持消息的持久化必须设置消息是持久化的标识。
	err = ch.Publish(
		"",
		q.Name, // routingKey
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: []byte(body),
		},
		)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x]Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
