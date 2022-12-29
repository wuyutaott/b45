package znet

import (
	"errors"
	"fmt"
	"github.com/wuyutaott/b45/utils"
	"github.com/wuyutaott/b45/ziface"
	"net"
)

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	msgHandler  ziface.IMsgHandler
	ConnMgr     ziface.IConnManager
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
	exitChan    chan struct{}
}

func NewServer() ziface.IServer {
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s, listen at IP: %s, Port: %d is starting \n", s.Name, s.IP, s.Port)
	s.exitChan = make(chan struct{})

	go func() {
		s.msgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			panic(err)
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		fmt.Printf("start b45 server %s success... \n", s.Name)

		var cid uint32 = 0

		go func() {
			for {
				if s.ConnMgr.Len() >= utils.GCfg.MaxConn {
					fmt.Println("over the max connection num:", utils.GCfg.MaxConn)
					AcceptDelay.Delay()
					continue
				}

				conn, err := listener.AcceptTCP()
				if err != nil {
					//GO 1.16+
					if errors.Is(err, net.ErrClosed) {
						fmt.Println("Listener closed")
						return
					}
					fmt.Println("Accept err:", err)
					AcceptDelay.Delay()
					continue
				}

				AcceptDelay.Reset()

				dealConn := NewConnection(s, conn, cid, s.msgHandler)
				cid++

				go dealConn.Start()
			}
		}()

		select {
		case <-s.exitChan:
			err := listener.Close()
			if err != nil {
				fmt.Println("listener close err:", err)
			}
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] b45 server")
	s.ConnMgr.ClearConn()
	s.exitChan <- struct{}{}
	close(s.exitChan)
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop")
		s.OnConnStop(conn)
	}
}
