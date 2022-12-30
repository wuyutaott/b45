package main

import (
	"fmt"
	"net"
)

func reader(conn net.Conn) {

}

func writer(conn net.Conn) {
	//for {
	//	fmt.Println("写入数据ping")
	//	conn.Write([]byte("ping"))
	//	time.Sleep(1 * time.Second)
	//}
}

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}
	fmt.Println("connection success")

	go reader(conn)
	go writer(conn)

	select {}
}
