package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	proto "../proto"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	p := proto.NewRequestSeguimientoClient(conn)

	response, err := p.Request(context.Background(), &proto.QuerySeguimiento{Seguimiento: "123456789"})
	if err != nil {
		log.Fatalf("Error when calling Request: %s", err)
	}
	log.Printf("Response from server: %s", response.Estado)

}