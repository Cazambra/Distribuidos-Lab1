package main

import (
	"log"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
	"strconv"
	"math/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	proto "../../logistica/proto"
)

type order struct {
	id string
	producto string
	valor int
	tienda string
	destino string
	tipo string
}

var seguimientos []int64;
var wg sync.WaitGroup

func leer (path string) ([][]string, error) {
	recordFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Ocurrió un error: ", err)
	}

	reader := csv.NewReader(recordFile)

	//ignorar la primera linea
	if _, err := reader.Read(); err != nil {
		panic(err)
	}

	records, err := reader.ReadAll()

	return records, err
}

func send_orders() {
	//conexión
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	fmt.Println("Ingrese tipo de cliente (pyme/retail): ")
	var input string;
	fmt.Scanln(&input)

	for true {
		switch input {
		case "pyme":
			recordspyme, _ := leer("../../pymes.csv")
			for _, line := range recordspyme {
				val, _ := strconv.ParseInt(line[2], 10, 64)
				o := proto.Order{
					Id: line[0],
					Producto: line[1],
					Valor: val,
					Tienda: line[3],
					Destino: line [4],
					Tipo: line[5], //es 0 o 1 porque es pyme, RETAIL PARA EL OTRO CSV
				}

				p := proto.NewLogisticaServiceClient(conn)

				//se envia la orden
				response, err := p.SendOrder(context.Background(), &o)
				if err !=nil {
					log.Fatalf("Error when calling Request: %s", err)
				}
				log.Printf("Orden enviada: %s ", o.GetId(), " número de seguimiento: %d", response.Seguimiento)

				seguimientos = append(seguimientos, response.Seguimiento)
				time.Sleep(time.Second)

			}
		case "retail":
			recordsret, _ := leer("../../retail.csv")
			//aqui se envían las órdenes retail 1x1
			for _, line := range recordsret {
				val, _ := strconv.ParseInt(line[2], 10, 64)
				o := proto.Order{
					Id: line[0],
					Producto: line[1],
					Valor: val,
					Tienda: line[3],
					Destino: line [4],
					Tipo: "retail",
				}

				p := proto.NewLogisticaServiceClient(conn)

				//se envia la orden
				response, err := p.SendOrder(context.Background(), &o)
				if err !=nil {
					log.Fatalf("Error when calling Request: %s", err)
				}
				log.Printf("Orden enviada: %s ", o.GetId())
				time.Sleep(time.Second)
			}
		}
	}
	defer wg.Done()
}

func consultar_estado() {
	time.Sleep(10*time.Second)
	//conexión
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	p := proto.NewLogisticaServiceClient(conn)

	//se consulta por el estado
	for true{
		i := rand.Intn(range seguimientos)
		response2, err2 := p.Request(context.Background(), seguimientos[i])
		if err2 !=nil {
			log.Fatalf("Error when calling Request: %s", err2)
		}
		log.Printf("Seguimiento: %d", seguimientos[i], " Estado: %+v ", response2.Estado)
		time.Sleep(2*time.Second)
	}
}

func main() {
	//goroutines
	wg.Add(2)
	go send_orders()
	go consultar_estado()
	log.Println("[x] Ctrl + C para detener")
	wg.Wait()
}
