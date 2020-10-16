package main

import (
	"log"
	"encoding/json"
	"github.com/streadway/amqp"
)

type packet struct {
	Id string
	Seguimiento string
	Valor int
	Tipo string
	Intentos int
	Estado string
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://tete:moraga@10.10.28.41:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()

	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"financiero", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")
	for i := 1; i < 5; i++{
		body := packet{"0002","00002",24,"retail",2,"Recibido"}
		b,err := json.Marshal(body)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/json",
				Body:       []byte(b),
			})
		log.Printf(" [x] Sent %+v", body)
		failOnError(err, "Failed to publish a message")
	}
	
}