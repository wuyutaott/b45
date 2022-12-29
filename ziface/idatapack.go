package ziface

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}

const (
	B45DataPack string = "b45_pack"
)

const (
	B45Message string = "b45_message"
)
