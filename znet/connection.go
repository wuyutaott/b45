package znet

import (
	"context"
	"github.com/wuyutaott/b45/ziface"
	"net"
)

type Connection struct {
	TCPServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return nil
}

func (c *Connection) StartWriter() {

}

func (c *Connection) StartReader() {

}

func (c *Connection) Start() {

}

func (c *Connection) Stop() {

}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return nil
}

func (c *Connection) GetConnID() uint32 {
	return 0
}

func (c *Connection) RemoteAddr() net.Addr {
	return nil
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	return nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {

}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	return nil, nil
}

func (c *Connection) RemoveProperty(key string) {

}

func (c *Connection) Context() context.Context {
	return nil
}

func (c *Connection) finalizer() {

}
