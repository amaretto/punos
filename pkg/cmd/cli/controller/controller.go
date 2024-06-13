package controller

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/amaretto/punos/pkg/cmd/cli/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// NewCommand create command
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "controller",
		Short: "start punos controller",
		Long:  `start punos controller`,
		Run: func(cmd *cobra.Command, args []string) {
			controller()
		},
	}
	return c
}

const (
	address = "localhost:19003"
)

func resgistTT(ctx context.Context, c pb.CtrlClient, id string) error {
	r, err := c.RegistTT(ctx, &pb.TTRegistRequest{Id: id})
	if err != nil {
		return err
	}
	if r.Result {
		log.Printf("Register turntable %s successful!\n", id)
	} else {
		log.Printf("Register turntable %s failed!\n", id)
	}
	return nil
}

func resgistCtrl(ctx context.Context, c pb.CtrlClient, id string) error {
	r, err := c.RegistCtrl(ctx, &pb.CtrlRegistRequest{Id: id})
	if err != nil {
		return err
	}
	if r.Result {
		log.Printf("Register ctrl %s successful!\n", id)
	} else {
		log.Printf("Register ctrl %s failed!\n", id)
	}
	return nil
}

func sendTTCmd(c pb.CtrlClient) error {
	stdin := bufio.NewScanner(os.Stdin)
	stream, err := c.SendTTCmd(context.Background())
	if err != nil {
		return err
	}
	log.Printf("Start to send command")

	for {
		stdin.Scan()
		text := stdin.Text()
		if err := stream.Send(&pb.TTCmdRequest{Cmd: text}); err != nil {

		}
		if text == "exit" {
			break
		}
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}
	return nil
}

func getTTCmd(c pb.CtrlClient, id string) error {
	req := &pb.GetTTCmdRequest{Id: id}
	stream, err := c.GetTTCmd(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("Get cmd: %s\n", msg.Cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func controller() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect server\n")
	}
	defer conn.Close()
	c := pb.NewCtrlClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = resgistCtrl(ctx, c, strconv.Itoa(int(time.Now().Unix())))
	if err != nil {
		log.Printf("Failed to exec register Ctrl command :%v\n", os.Args)
	}

	err = sendTTCmd(c)
	if err != nil {
		log.Printf("Failed to start sending command :%v\n", os.Args)
	}
}
