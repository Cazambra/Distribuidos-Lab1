package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	proto "../../logistica/proto"
)

type camion struct {
	tipo string //r1, r2 o normal
	status bool
}


func main()  {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
    /*
	camion_r1 := camion{
		tipo: "Retail 1",
		status: true,
	}
	camion_r2 := camion{
		tipo: "Retail 2",
		status: true,
	}
	camion_n := camion{
		tipo: "Normal",
		status: true,
	}
	*/
	p := proto.NewLogisticaServiceClient(conn)

	response, err := p.Ready(context.Background(), &proto.ReadyAdvice{Tipo: "Retail 1", UltRet: true})
	if err != nil {
		log.Fatalf("Error when calling Request: %s", err)
	}
	log.Printf("Response from server: %+v", response)
}