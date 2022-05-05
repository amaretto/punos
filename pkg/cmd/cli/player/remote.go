package player

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/amaretto/punos/pkg/cmd/cli/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func setupRemoteControl(p Player) {
	// ToDo : separate method
	address := "localhost:19003"
	// create gRPC Client
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect server\n")
	}
	defer conn.Close()
	c := pb.NewCtrlClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// ToDo : register turn table
	var flag bool
	r, err := c.RegistTT(ctx, &pb.TTRegistRequest{Id: p.playerID})
	if err != nil {
		logrus.Debug(err)
	}
	if r.Result {
		logrus.Debugf("Register turntable %s successful!\n", p.playerID)
		flag = true
	} else {
		logrus.Debugf("Register turntable %s failed!\n", p.playerID)
	}

	// ToDo : getTTCmd
	req := &pb.GetTTCmdRequest{Id: p.playerID}
	stream, err := c.GetTTCmd(context.Background(), req)
	if err != nil {
		logrus.Debug(err)
	}

	if flag {
		go func() {
			for {
				logrus.Debug("from remote controller")
				msg, err := stream.Recv()
				if err == io.EOF {
					flag = false
				}
				if msg.Cmd != "" {
					switch msg.Cmd[0] {
					case 'a':
						logrus.Debug("hogehoge")
					case 'l':
						p.Fforward()
					case 'h':
						p.Rewind()
					case 'j':
						p.Voldown()
					case 'k':
						p.Volup()
					case 'm':
						p.Spdup()
					case ',':
						p.Spddown()
					}
				}
			}
		}()
	}
}
