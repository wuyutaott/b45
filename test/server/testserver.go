package main

import (
	"fmt"
	"github.com/wuyutaott/b45/ziface"
	"github.com/wuyutaott/b45/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (ping *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("receive msg =", string(request.GetMsgID()))
	request.GetConnection().SendMsg(2, []byte("pong"))
}

func main() {
	server := znet.NewServer()
	server.AddRouter(1, &PingRouter{})
	server.Serve()
}
