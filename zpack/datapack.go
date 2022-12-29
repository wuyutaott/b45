package zpack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/wuyutaott/b45/utils"
	"github.com/wuyutaott/b45/ziface"
)

var defaultHeaderLen uint32 = 8

type DataPack struct{}

func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return defaultHeaderLen
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	if utils.GCfg.MaxPacketSize > 0 && msg.DataLen > utils.GCfg.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}
	return msg, nil
}
