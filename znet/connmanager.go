package znet

import (
	"errors"
	"fmt"
	"github.com/wuyutaott/b45/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	cm.connections[conn.GetConnID()] = conn
	cm.connLock.Unlock()
	fmt.Println("connection add to ConnManager success! conn num =", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	delete(cm.connections, conn.GetConnID())
	cm.connLock.Unlock()
	fmt.Printf("connection remove connID = %d success! conn num = %d \n", conn.GetConnID(), cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (cm *ConnManager) Len() int {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	for _, conn := range cm.connections {
		conn.Stop()
	}
	cm.connLock.Unlock()
	fmt.Println("Clear all connection success! conn num =", cm.Len())
}

func (cm *ConnManager) ClearOneConn(connID uint32) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	connections := cm.connections
	if conn, ok := connections[connID]; ok {
		conn.Stop()
		fmt.Printf("Clear connection id: %d success! \n", connID)
		return
	}
	fmt.Printf("Clear connection id: %d err! \n", connID)
}
