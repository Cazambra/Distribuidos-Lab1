package main
import "fmt"
import (
	"log"
	"encoding/json"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"syscall"
)
var  inicial = 0
var enviostot = 0
var nocomplt = 0
type packet struct {
	Id string
	Seguimiento string
	Valor int
	Tipo string
	Intentos int
	Estado string
}

func financiero(inicial int, enviostot int, nocomplt int,paquete *packet) (int,int,int) {

	 enviostot +=  1	
	 fmt.Println(paquete.Estado)
	if paquete.Estado == "No Recibido"{
			nocomplt +=1

	}

	 inicial += paquete.Valor
	return inicial,enviostot,nocomplt

}
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Balance total: %d", inicial)
		fmt.Println("Envios totales: %d", enviostot)
		fmt.Println("Envios no Completados: %d", nocomplt)
		os.Exit(0)
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	SetupCloseHandler()

	conn, err := amqp.Dial("amqp://tete:moraga@localhost:5672/")
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var m *packet
			err := json.Unmarshal(d.Body,&m)
			inicial,enviostot,nocomplt = financiero(inicial,enviostot,nocomplt,m)

			fmt.Println("Received a message: %+v", m)
			failOnError(err,"FaileD to receive message")
		}
	}()


	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	

}