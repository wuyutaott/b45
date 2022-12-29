package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgID uint32, router IRouter)
	GetConnMgr() IConnManager
	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
	Packet() IDataPack
}
