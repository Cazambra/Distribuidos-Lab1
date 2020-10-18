package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"strconv"
	"container/list"
	proto "../proto"
	"google.golang.org/grpc"
)

//Server ...
type Server struct{
	paquetes []proto.Packet
	registros []proto.Register
	qR *list.List
	qP *list.List
	qN *list.List
}

type Paquete struct{
	id string
	seguimiento int64
	tipo string
	valor int64
	intentos int
	estado string
}

var cont int = 0 //contador para seguimientos
var registros []proto.Register //slice con registros
var paquetes []proto.Packet //slice con los paquetes

var qR = list.New()
var qP = list.New()
var qN = list.New()

//Ready ...
func (s *Server) Ready(ctx context.Context,adv *proto.ReadyAdvice) (*proto.Deliver, error){
	deli := proto.Deliver{}

	for deli.Segundo == nil { //agregar OR para tiempo de espera
		switch adv.Tipo {
		case "Retail 1", "Retail 2":
			if adv.GetUltRet() {
				if qR.Len() == 0 {
					packet := qP.Front()
					if deli.Primero == nil {
						deli.Primero = packet.Value.(*proto.Packet)
					} else{
						deli.Segundo = packet.Value.(*proto.Packet)
					}
					qP.Remove(packet)
				} else {
					packet := qR.Front()
					if deli.Primero == nil {
						deli.Primero = packet.Value.(*proto.Packet)
					} else{
						deli.Segundo = packet.Value.(*proto.Packet)
					}
					qR.Remove(packet)
				}
			} else {
				packet := qR.Front()
				if deli.Primero == nil{
					deli.Primero = packet.Value.(*proto.Packet)
				} else{
					deli.Segundo = packet.Value.(*proto.Packet)
				}
				qR.Remove(packet)
			}
		case "Normal":
			if qP.Len() == 0 {
				packet := qN.Front()
				if deli.Primero == nil {
					deli.Primero = packet.Value.(*proto.Packet)
				} else{
					deli.Segundo = packet.Value.(*proto.Packet)
				}
				qN.Remove(packet)	
			} else {
				packet := qP.Front()
				if deli.Primero == nil {
					deli.Primero = packet.Value.(*proto.Packet)
				} else{
					deli.Segundo = packet.Value.(*proto.Packet)
				}
				qP.Remove(packet)	
			}
		}
	}
	log.Println("Saliendo paquetes: %s", deli.Primero.GetIdPacket(), deli.Segundo.GetIdPacket())
	//HAY QUE ACTUALIZAR EL ESTADO DE LOS PAQUETES QUE SALEN (!)
	for i := range paquetes{
		if  paquetes[i].GetIdPacket() == deli.Primero.GetIdPacket(){
			paquetes[i].Estado = "en camino"
			deli.Primero.Estado = "en camino"
		} else if paquetes[i].GetIdPacket() == deli.Segundo.GetIdPacket(){
			paquetes[i].Estado = "en camino"
			deli.Segundo.Estado = "en camino"
		} else{
			continue
		}
	}
	return &deli, nil
}

//Delievered ...
func (s *Server) Delivered(ctx context.Context,deli *proto.Deliver) (*proto.ReplySeguimiento, error){
	// actualizar []paquetes
	for i := range paquetes {
		if paquetes[i].GetIdPacket() == deli.Primero.GetIdPacket() {
			paquetes[i] = *deli.Primero
		} else if paquetes[i].GetIdPacket() == deli.Segundo.GetIdPacket(){
			paquetes[i] = *deli.Segundo
		} else{
			continue
		}
	}
	//mandar a financiero

	return &proto.ReplySeguimiento{Estado: ""}, nil
}


//Request ...
func (s *Server) Request(ctx context.Context, num_seguimiento *proto.QuerySeguimiento) (*proto.ReplySeguimiento, error) {
	fmt.Println("Enviando seguimiento ")
	//aqui hay que consultar por el .estado de los PAQUETES
	var status string
	for i := range paquetes {
		if paquetes[i].Seguimiento == num_seguimiento.Seguimiento {
			status = paquetes[i].Estado
		} 
	}
	return &proto.ReplySeguimiento{Estado: status}, nil
}
//Send order ...
func (s *Server) SendOrder(ctx context.Context, orden *proto.Order) (*proto.QuerySeguimiento, error){
	fmt.Println("Paquete recibido: %s ", orden.GetId())
	//aqui se debe generar el codigo de seguimiento 
	cont = cont + 1
	seguimiento := int64(cont)
	
	//armar el REGISTRO (timestamp, id-paquete, etc) para despues mapearlo y armar los PAQUETES
	new_reg := proto.Register{
		Timestamp: time.Now().Format("02-01-2006 15:04"),
		IdPacket: strconv.Itoa(cont),
		Tipo: "null",
		Nombre: orden.GetProducto(),
		Valor: orden.GetValor(),
		Origen: orden.GetTienda(),
		Destino: orden.GetDestino(),
		Seguimiento: int64(0),
	}
	
	if orden.GetTienda() == "pyme" { // si es pyme entra aqui
		if orden.GetTipo() == "0"{ // si es normal
			new_reg.Tipo = "normal"
		} else { //si es prioritario
			new_reg.Tipo = "prioritario"
		}
		new_reg.Seguimiento = seguimiento
		//fmt.Println(new_reg)
	} else { // si es retail
		new_reg.Tipo = "retail"
		//fmt.Println("%+v", new_reg)
	}
	registros = append(registros, new_reg)

	// mapear los registros -> paquetes
	new_pack := proto.Packet{
		IdPacket: new_reg.IdPacket,
		Seguimiento: seguimiento,
		Tipo: new_reg.Tipo,
		Valor: new_reg.Valor,
		Intentos: int64(0),
		Estado: "en bodega",
	}
	paquetes = append(paquetes, new_pack)
	//fmt.Println(new_pack, "\n")

	//se agregan a la cola los paquetes
	switch new_pack.Tipo {
	case "retail":
		qR.PushBack(&new_pack)
		//fmt.Println("%+v", qR.Back())
	case "prioritario":
		qP.PushBack(&new_pack)
	case "normal":
		qN.PushBack(&new_pack)
	}
	return &proto.QuerySeguimiento{Seguimiento: seguimiento}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterLogisticaServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	lis1, err1 := net.Listen("tcp", ":9001")
	if err1 != nil {
		log.Fatalf("failed to listen: %v", err1)
	}

	grpcServer1 := grpc.NewServer()
	proto.RegisterLogisticaServiceServer(grpcServer1, &Server{})

	if err1 := grpcServer1.Serve(lis1); err1 != nil {
		log.Fatalf("failed to serve: %s", err1)
	}	

}
