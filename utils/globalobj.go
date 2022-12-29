package utils

type GlobalCfg struct {
	Name    string
	Version string
	Host    string
	TCPPort int

	MaxPacketSize    uint32
	MaxConn          int
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxMsgChanLen    uint32
}

var GCfg *GlobalCfg

func init() {
	GCfg = &GlobalCfg{
		Name:             "b45App",
		Version:          "v1.0",
		Host:             "0.0.0.0",
		TCPPort:          8999,
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}
}
