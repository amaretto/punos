package server

import (
	"fmt"
	"log"
	"net"

	pb "github.com/amaretto/punos/pkg/cmd/cli/pb"
	"github.com/amaretto/punos/pkg/cmd/cli/service"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

// NewCommand create command
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "server",
		Short: "start punos server",
		Long:  `start punos server`,
		Run: func(cmd *cobra.Command, args []string) {
			server()
		},
	}
	return c
}

func server() {
	port := 19003
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Run server port: %d", port)

	grpcServer := grpc.NewServer()
	pb.RegisterCtrlServer(grpcServer, &service.DJService{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
