package main

import (
	"log"
	"net"

	pb "github.com/amaretto/punos/cmd/test/grpc/pb"
	"github.com/amaretto/punos/cmd/test/grpc/service"

	"google.golang.org/grpc"
)

func main() {
	listenPort, err := net.Listen("tcp", ":19003")
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	catService := &service.MyCatService{}

	pb.RegisterCatServer(server, catService)
	server.Serve(listenPort)
}
