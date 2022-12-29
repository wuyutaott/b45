package znet

import "github.com/wuyutaott/b45/ziface"

type Request struct {
	conn   ziface.IConnection
	msg    ziface.IMessage
	router ziface.IRouter
	index  int8
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

func (r *Request) Next() {
	r.index++
	for r.index < 4 {
		switch r.index {
		case 1:
			r.router.PreHandle(r)
		case 2:
			r.router.Handle(r)
		case 3:
			r.router.PostHandle(r)
		}
		r.index++
	}
}

func (r *Request) Abort() {
	r.index = 4
}
