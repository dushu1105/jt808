package jtnet

import (
	"errors"
	"fmt"
	"sync"
	cmap "github.com/orcaman/concurrent-map"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[uint32]*Connection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
	validConnects cmap.ConcurrentMap //sim号标记的连接
	idSimMap map[uint32]string
}

/*
	创建一个链接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]*Connection),
		validConnects: cmap.New(),
		idSimMap:make(map[uint32]string),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn *Connection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

//删除连接
func (connMgr *ConnManager) Remove(conn *Connection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	if sim, ok := connMgr.idSimMap[conn.GetConnID()];ok{
		delete(connMgr.idSimMap, conn.GetConnID())
		if connMgr.validConnects.Has(sim){
			connMgr.validConnects.Remove(sim)
		}
	}

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

//利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (*Connection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
		if sim, ok := connMgr.idSimMap[connID];ok{
			delete(connMgr.idSimMap, connID)
			if connMgr.validConnects.Has(sim){
				connMgr.validConnects.Remove(sim)
			}
		}
	}

	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) RelatedSim(sim string, conn *Connection) {
	if connMgr.validConnects.Has(sim){
		c, _ := connMgr.validConnects.Get(sim)
		c.(*Connection).Stop()
	}

	connMgr.validConnects.Set(sim, conn)
	connMgr.idSimMap[conn.ConnID] = sim
	fmt.Println(sim, "tag connection", conn.ConnID)
}

func (connMgr *ConnManager) GetBySim(sim string) (*Connection, error) {
	if conn, ok := connMgr.validConnects.Get(sim); ok {
		return conn.(*Connection), nil
	} else {
		return nil, errors.New("connection not found")
	}
}