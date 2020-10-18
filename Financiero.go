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
var  inicial = int64(0)
var enviostot = int64(0)
var nocomplt = int64(0)
type packet struct {
	IdPacket string
	Seguimiento int64
	Valor int64
	Tipo string
	Intentos int64
	Estado string
}

func financiero(inicial int64, enviostot int64, nocomplt int64,paquete *packet) (int64,int64,int64) {
	fmt.Println(enviostot)
	 enviostot +=  1	
	 fmt.Println(paquete.Estado)
	if paquete.Estado == "no recibido"{
			nocomplt +=1

	}

	 inicial += paquete.Valor - 10*(paquete.Intentos-1)
	return inicial,enviostot,nocomplt

}
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Balance total: %d", inicial)
		fmt.Println("Envios totales: %d", enviostot)
		fmt.Println("Envios no Recibidos: %d", nocomplt)
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

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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