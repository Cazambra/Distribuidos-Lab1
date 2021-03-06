package main

import (
	"log"
	"sync"
	"fmt"
	"time"
	"math/rand"

	context "context"
	"google.golang.org/grpc"

	proto "../../logistica/proto"
)

type Camion struct {
	tipo string //r1, r2 o normal
	status bool
	UltRet bool
}

var wg sync.WaitGroup

func Reparto(camion Camion) int {

	for true{

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.10.28.42:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	p := proto.NewLogisticaServiceClient(conn)

	adv := proto.ReadyAdvice{
		Tipo: camion.tipo,
		UltRet: camion.UltRet,
	}
	log.Println("Camión listo: %s", camion.tipo)
	time.Sleep(time.Second)

	packets, err := p.Ready(context.Background(), &adv)
	if err != nil {
		log.Fatalf("Error when calling Request: %s", err)
	}
	//log.Printf("Response from server: %+v", packets)
	
	if packets.Primero.GetValor() >= packets.Segundo.GetValor(){
		switch packets.Primero.GetTipo() {
		case "normal", "prioritario": 
			pen := int64(0)
			for (packets.Primero.Estado == "no recibido" ||packets.Primero.Estado == "en camino") && packets.Primero.GetValor() > pen && packets.Primero.Intentos < 2{
				prob := rand.Intn(10)
				if prob <= 7 {
					packets.Primero.Intentos += 1
					packets.Primero.Estado = "recibido" 
				} else {
					packets.Primero.Intentos += 1
					packets.Primero.Estado = "no recibido" 
					pen += 10 
					if packets.Primero.GetValor() + 10 < pen {
						pen -= 10
						break
					}
				}
			}
			pen = 0
			for (packets.Segundo.Estado == "no recibido" ||packets.Segundo.Estado == "en camino") && packets.Segundo.GetValor() > pen && packets.Segundo.Intentos < 2{
				prob := rand.Intn(10)
				if prob <= 7 {
					packets.Segundo.Intentos += 1
					packets.Segundo.Estado = "recibido" 
				} else {
					packets.Segundo.Intentos += 1
					packets.Segundo.Estado = "no recibido" 
					pen += 10 
					if packets.Segundo.GetValor() + 10 < pen {
						pen -= 10
						break
					}
				}
			}
		case "retail":
			pen := int64(0)
			for (packets.Primero.Estado == "no recibido" ||packets.Primero.Estado == "en camino")	&& packets.Primero.Intentos < 3{
				prob := rand.Intn(10)
				if prob <= 7 {
					packets.Primero.Intentos += 1
					packets.Primero.Estado = "recibido" 
				} else {
					packets.Primero.Intentos += 1
					packets.Primero.Estado = "no recibido" 
					pen += 10 
					if packets.Primero.GetValor() + 10 < pen {
						pen -= 10
						break
					}
				}
			}
			pen = 0
			for (packets.Segundo.Estado == "no recibido" ||packets.Segundo.Estado == "en camino")	&& packets.Segundo.Intentos < 3{
				prob := rand.Intn(10)
				if prob <= 7 {
					packets.Segundo.Intentos += 1
					packets.Segundo.Estado = "recibido" 
				} else {
					packets.Segundo.Intentos += 1
					packets.Segundo.Estado = "no recibido" 
					pen += 10 
					if packets.Segundo.GetValor() + 10 < pen {
						pen -= 10
						break
					}
				}
			}
		}
	} else {
		switch packets.Primero.GetTipo() {
			case "normal", "prioritario": 
				pen := int64(0)
				for (packets.Primero.Estado == "no recibido" ||packets.Primero.Estado == "en camino")	&& packets.Primero.GetValor() > pen && packets.Primero.Intentos < 2{
					prob := rand.Intn(10)
					if prob <= 7 {
						packets.Primero.Intentos += 1
						packets.Primero.Estado = "recibido" 
					} else {
						packets.Primero.Intentos += 1
						packets.Primero.Estado = "no recibido" 
						pen += 10 
						if packets.Primero.GetValor() + 10 < pen {
							pen -= 10
							break
						}
					}
				}
				pen = 0
				for (packets.Segundo.Estado == "no recibido" ||packets.Segundo.Estado == "en camino")	&& packets.Segundo.GetValor() > pen && packets.Segundo.Intentos < 2{
					prob := rand.Intn(10)
					if prob <= 7 {
						packets.Segundo.Intentos += 1
						packets.Segundo.Estado = "recibido" 
					} else {
						packets.Segundo.Intentos += 1
						packets.Segundo.Estado = "no recibido" 
						pen += 10 
						if packets.Segundo.GetValor() + 10 < pen {
							pen -= 10
							break
						}
					}
				}
			case "retail":
				pen := int64(0)
				for (packets.Primero.Estado == "no recibido" ||packets.Primero.Estado == "en camino")	&& packets.Primero.Intentos < 3{
					prob := rand.Intn(10)
					if prob <= 7 {
						packets.Primero.Intentos += 1
						packets.Primero.Estado = "recibido" 
					} else {
						packets.Primero.Intentos += 1
						packets.Primero.Estado = "no recibido" 
						pen += 10 
						if packets.Primero.GetValor() + 10 < pen {
							pen -= 10
							break
						}
					}
				}
				pen = 0
				for (packets.Segundo.Estado == "no recibido" ||packets.Segundo.Estado == "en camino")	&& packets.Segundo.Intentos < 3{
					prob := rand.Intn(10)
					if prob <= 7 {
						packets.Segundo.Intentos += 1
						packets.Segundo.Estado = "recibido" 
					} else {
						packets.Segundo.Intentos += 1
						packets.Segundo.Estado = "no recibido" 
						pen += 10 
						if packets.Segundo.GetValor() + 10 < pen {
							pen -= 10
							break
						}
					}
				}
			}
	}
	_, err = p.Delivered(context.Background(), packets)
	if err != nil {
		log.Fatalf("Error when calling Request: %s", err)
	}
	if packets.Primero.GetTipo() == "retail" || packets.Segundo.GetTipo()== "retail"{
		camion.UltRet = true
	}
	time.Sleep(5*time.Second)
}
	defer wg.Done()
	return 0
}

func main()  {

	rand.Seed(time.Now().UnixNano())

	camion_r1 := Camion{
		tipo: "Retail 1",
		status: true,
	}
	camion_r2 := Camion{
		tipo: "Retail 2",
		status: true,
	}
	camion_n := Camion{
		tipo: "Normal",
		status: true,
	}


		wg.Add(3)
		go Reparto(camion_r1)
		time.Sleep(time.Second)
		go Reparto(camion_r2)
		time.Sleep(time.Second)
		go Reparto(camion_n)
		wg.Wait()
		time.Sleep(3*time.Second)
	
	
	fmt.Println("main finish")
}