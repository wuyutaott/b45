package zpack

import (
	"github.com/wuyutaott/b45/ziface"
	"sync"
)

var packOnce sync.Once

type packFactory struct{}

var factoryInstance *packFactory

func Factory() *packFactory {
	packOnce.Do(func() {
		factoryInstance = new(packFactory)
	})
	return factoryInstance
}

func (f *packFactory) NewPack(kind string) ziface.IDataPack {
	var dataPack ziface.IDataPack
	switch kind {
	case ziface.B45DataPack:
		dataPack = NewDataPack()
		break
	default:
		dataPack = NewDataPack()
	}
	return dataPack
}
