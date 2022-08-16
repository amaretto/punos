package service

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/amaretto/punos/pkg/cmd/cli/pb"
)

const (
	playPause = "1"
	forward   = "2"
	rewind    = "3"
	layout    = "2006-01-02T15:05:06"
)

// DJService is service for djing
type DJService struct {
	tts          []turntable
	ctrl         controller
	mu           sync.RWMutex
	execCommands map[string]command
}

type command struct {
	execTime int64
	cmd      string
}

type turntable struct {
	id string
}

type controller struct {
	id string
}

////--- TurnTable ---////

// RegistTT register TurnTable
func (s *DJService) RegistTT(ctx context.Context, p *pb.TTRegistRequest) (*pb.TTRegistResult, error) {
	log.Printf("Failed to start sending command :%s\n", p.Id)
	s.tts = append(s.tts, turntable{p.Id})
	if s.execCommands == nil {
		s.execCommands = make(map[string]command)
	}
	s.mu.Lock()
	s.execCommands[p.Id] = command{execTime: time.Now().UnixNano()}
	s.mu.Unlock()
	log.Printf("map :%v\n", s.execCommands)
	return &pb.TTRegistResult{Result: true}, nil
}

// GetTTCmd get controll command of TurnTable
func (s *DJService) GetTTCmd(p *pb.GetTTCmdRequest, stream pb.Ctrl_GetTTCmdServer) error {
	log.Printf("Start to send command to TurnTable :%s\n", p.Id)
	var lastExecTime int64
	var execTime int64
	var execCmd string
	for {
		s.mu.RLock()
		execTime = s.execCommands[p.Id].execTime
		execCmd = s.execCommands[p.Id].cmd
		s.mu.RUnlock()
		//log.Printf("lastExectime:%d,execCommandTime:%d", lastExecTime, execTime)
		//log.Printf("lastExectime:%d,execCommandTime:%d", lastExecTime, s.execCommands[p.Id].execTime)
		if lastExecTime < execTime {
			log.Printf("Execution command :%s at %d\n", execCmd, execTime)
			if err := stream.Send(&pb.TTCmd{Cmd: execCmd, Param: ""}); err != nil {
				return err
			}
		}
		lastExecTime = execTime
	}
}

////--- Controller ---////

// RegistCtrl register Controller
func (s *DJService) RegistCtrl(ctx context.Context, p *pb.CtrlRegistRequest) (*pb.CtrlRegistResult, error) {
	log.Printf("Regist controller :%s\n", p.Id)
	s.ctrl = controller{p.Id}
	return &pb.CtrlRegistResult{Result: true}, nil
}

// SendTTCmd accept command
func (s *DJService) SendTTCmd(stream pb.Ctrl_SendTTCmdServer) error {
	log.Printf("Start to accept command from Controller")
	for {
		m, err := stream.Recv()
		log.Printf("Accept command from Controller cmd:%s param:%s\n", m.Cmd, m.Param)

		if err == io.EOF {
			return stream.SendAndClose(&pb.TTCmdResult{
				Result: true,
			})
		}
		if err != nil {
			return err
		}
		if m.Cmd == "exit" {
			return stream.SendAndClose(&pb.TTCmdResult{
				Result: true,
			})
		}
		for _, tt := range s.tts {
			s.mu.Lock()
			s.execCommands[tt.id] = command{time.Now().UnixNano(), m.Cmd}
			s.mu.Unlock()
		}
		log.Printf("exec commands%v\n", s.execCommands)
	}
}
