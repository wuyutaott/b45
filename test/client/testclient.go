package main

import (
	"fmt"
	"github.com/wuyutaott/b45/zpack"
	"net"
	"sync"
	"time"
)

var pack = zpack.NewDataPack()

func reader(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer fmt.Println("reader exit")

	for {
		headerBuff := make([]byte, pack.GetHeadLen())
		_, err := conn.Read(headerBuff)
		if err != nil {
			fmt.Println("read header err:", err)
			return
		}
		msg, err := pack.Unpack(headerBuff)
		if err != nil {
			fmt.Println("unpack err:", err)
			return
		}
		if msg.GetDataLen() > 0 {
			bodyBuff := make([]byte, msg.GetDataLen())
			_, err := conn.Read(bodyBuff)
			if err != nil {
				fmt.Println("read body err:", err)
				return
			}
			msg.SetData(bodyBuff)
		}
		fmt.Println("收到服务器消息", msg.GetMsgID(), string(msg.GetData()))
	}
}

func writer(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer fmt.Println("writer exit")

	for {
		msg := zpack.NewMsgPackage(1, []byte("ping"))
		data, err := pack.Pack(msg)
		if err != nil {
			fmt.Println("pack err:", err)
			return
		}

		if _, err := conn.Write(data); err != nil {
			fmt.Println("write err:", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}
	fmt.Println("connection success")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go reader(wg, conn)
	go writer(wg, conn)
	wg.Wait()

	fmt.Println("main exit")
}
