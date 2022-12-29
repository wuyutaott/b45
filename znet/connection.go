package znet

import (
	"context"
	"errors"
	"fmt"
	"github.com/wuyutaott/b45/utils"
	"github.com/wuyutaott/b45/ziface"
	"github.com/wuyutaott/b45/zpack"
	"io"
	"net"
	"sync"
	"time"
)

type Connection struct {
	sync.RWMutex
	Server       ziface.IServer
	Conn         *net.TCPConn
	ConnID       uint32
	MsgHandler   ziface.IMsgHandler
	ctx          context.Context
	cancel       context.CancelFunc
	msgBuffChan  chan []byte
	property     map[string]interface{}
	propertyLock sync.Mutex
	isClosed     bool
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Server:      server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgBuffChan: make(chan []byte, utils.GCfg.MaxMsgChanLen),
		property:    nil,
	}
	c.Server.GetConnMgr().Add(c)
	return nil
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer goroutine is running!")
	defer fmt.Println("Writer goroutine exit!")

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send buff data err:", err, "conn writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed!")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running!")
	defer fmt.Println("Reader goroutine exit!")
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			headData := make([]byte, c.Server.Packet().GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("read msg head err:", err)
				return
			}

			msg, err := c.Server.Packet().Unpack(headData)
			if err != nil {
				fmt.Println("unpack err:", err)
				return
			}

			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("read msg data err:", err)
					return
				}
			}
			msg.SetData(data)

			req := &Request{
				conn:  c,
				msg:   msg,
				index: 0,
			}

			if utils.GCfg.WorkerPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(req)
			} else {
				go c.MsgHandler.DoMsgHandler(req)
			}
		}
	}
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.Server.CallOnConnStart(c)
	go c.StartReader()
	go c.StartWriter()
	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

func (c *Connection) Stop() {
	c.cancel()
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	dp := c.Server.Packet()
	msg, err := dp.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msgID =", msgID)
		return errors.New("pack msg error")
	}

	_, err = c.Conn.Write(msg)
	return err
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	idleTimeout := time.NewTimer(5 * time.Microsecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	dp := c.Server.Packet()
	msg, err := dp.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msgID =", msgID)
		return errors.New("pack msg error")
	}

	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) finalizer() {
	c.Server.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	if c.isClosed == true {
		return
	}

	fmt.Println("Conn stop connID =", c.ConnID)

	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Conn close err:", err)
	}

	c.Server.GetConnMgr().Remove(c)

	close(c.msgBuffChan)

	c.isClosed = true
}
