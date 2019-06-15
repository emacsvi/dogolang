package main

import (
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 定义一个队列，定义队列操作具备幂等性，也就是说多次重复定义，相同名称的队列只会创建一个队列。
	// name: 队列名称
	// durable: 是否持久化
	// autoDelete: 自动删除
	// exclusive: 排他队列,如果你希望创建一个队列,并且只有你当前这个程序(或进程)进行消费处理.不希望别的客户端读取到这个队列.用这个方法甚好.而同时如果当进程断开连接.这个队列也会被销毁.不管是否设置了持久化或者自动删除.
	// noWait: 不等待处理结果
	// arguments: 额外的参数

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World! 您好，这是第一条golang测试消息。"

	// mandatory: 当mandatory标志位设置为true时，如果exchange根据自身类型和消息routeKey无法找到一个符合条件的queue，那么会调用basic.return方法将消息返回给生产者（Basic.Return + Content-Header + Content-Body）；当mandatory设置为false时，出现上述情形broker会直接将消息扔掉。
	// immediate: 当immediate标志位设置为true时，如果exchange在将消息路由到queue(s)时发现对于的queue上么有消费者，那么这条消息不会放入队列中。当与消息routeKey关联的所有queue（一个或者多个）都没有消费者时，该消息会通过basic.return方法返还给生产者。
	// 概括来说，mandatory标志告诉服务器至少将该消息route到一个队列中，否则将消息返还给生产者；immediate标志告诉服务器如果该消息关联的queue上有消费者，则马上将消息投递给它，如果所有queue都没有消费者，直接把消息返还给生产者，不用将消息入队列等待消费者了。
	err = ch.Publish(
		"", // exchange: 这里交换机为空，则会走默认的一个交换机，而默认的交换机会将消息送给hello的队列
		q.Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
