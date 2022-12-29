package znet

import "github.com/wuyutaott/b45/ziface"

type Request struct {
	conn   ziface.IConnection
	msg    ziface.IMessage
	router ziface.IRouter
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

func (r *Request) BindRouter(router ziface.IRouter) {
	r.router = router
}
