package jtnet

import (
	"github.com/dushu1105/jt808/protocal"
	"github.com/dushu1105/jt808/utils"
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer *Server
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler *MsgHandle
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	heartBeatResetChan chan bool

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
	Buf          []byte
	Fo           *os.File
	W            *bufio.Writer
	Sim          string
	TerminalVer uint8
	WW           *bufio.Writer
	NextSeq      uint16
	responseResult  map[uint16]map[uint16]RespResult
	TcpTimeout   int64
	HeartBeatTimeout int64
}

type RespResult struct {
	t int64
	result byte
}

//创建连接的方法
func NewConntion(server *Server, conn *net.TCPConn, connID uint32, msgHandler *MsgHandle) *Connection {
	//初始化Conn属性
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
		Buf:          make([]byte, 0),
		responseResult: make(map[uint16]map[uint16]RespResult),
		TcpTimeout: -1,
		HeartBeatTimeout: -1,
		heartBeatResetChan: make(chan bool, 1),
	}

	//c.IsHaiKangHeader = true
	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartHeartBeater(){
	fmt.Println("[HeartBeater Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn HeartBeater exit!]")

	t := time.Duration(36000)
	for {
		select {
		case <-time.After(t * time.Second):
			if c.HeartBeatTimeout > 0{
				fmt.Println(c.RemoteAddr().String(), "[conn HeartBeater timeout!]")
				c.Stop()
			}
		case <-c.heartBeatResetChan:
			if c.HeartBeatTimeout > 0{
				t = time.Duration(c.HeartBeatTimeout)
			} else {
				t = time.Duration(36000)
			}

			break
		case <-c.ExitBuffChan:
			return
		}
	}

}

/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			c.NextSeq += 1

			fmt.Printf("Send data succ! data = %d\n", len(data))
		case <-c.ExitBuffChan:
			return
		}
	}
}

func splitPackage(data []byte) [][2]int {
	beginArray := make([][2]int, 0)
	num := 0
	begin := -1
	for i, d := range data {
		if d == protocal.JT808Sign {
			num += 1
			if num % 2 == 1{
				begin = i
			} else {
				a := [2]int{begin, i}
				beginArray = append(beginArray, a)
				begin = -1
			}
		}
	}

	if begin != -1 {
		beginArray = append(beginArray, [2]int{begin, 0})
	}
	return beginArray
}

/*
	读消息Goroutine，用于从客户端中读取数据
*/
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()
	var leftData []byte
	for {
		//读取客户端的Msg head
		data := make([]byte, 2048)

		if n, err := c.Conn.Read(data); err != nil {
			fmt.Println("read msg head error ", err)
			break
		} else {
			//if utils.GlobalObject.FakeData{
			//	c.W.Write(headData)
			//}
			if len(leftData) != 0{
				data = append(leftData, data[:n]...)
				leftData = make([]byte, 0)
			}
		}

		if c.HeartBeatTimeout > 0 {
			c.heartBeatResetChan<-true
		}

		dataArrayPos := splitPackage(data)
		var begin, end int
		for _, posArray := range dataArrayPos {
			begin = posArray[0]
			end = posArray[1]

			if end == 0 {
				leftData = data[begin:]
				break
			}

			var jtMsg protocal.JT808Msg
			jtMsg.ConnId = c.GetConnID()
			err := jtMsg.Parse(data[begin : end+1])
			if err != nil {
				jtMsg.Print("jtMsg.Parse 有错，退出", err)
				return
			}

			if !jtMsg.IsCompleted {
				continue
			}


			ret, err := c.MsgHandler.DoMsgHandler(&jtMsg)
			if err != nil {
				jtMsg.Print("DoMsgHandler 有错，退出", err)
				return
			}

			if !ret.NeedFeedBack{
				//表明不需要回数据
				c.SetRespResult(ret.Result.Id, ret.Result.Seq, ret.Result.Result)
				continue
			}

			data, err = ret.Msg.Packet()
			if err != nil {
				jtMsg.Print("结果打包 有错，退出", err)
				return
			}
			err = c.SendData(data)

			if jtMsg.IsVerifyMsg(){
				if jtMsg.Is2019Ver(){

				}
				c.TcpServer.ConnMgr.RelatedSim(jtMsg.Header.Sim, c)
			}
			//fmt.Println("send msg id,", ret.Header.Id, "seq,", ret.Header.Seq)
		}
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
	go c.StartHeartBeater()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer heartbeat
	c.ExitBuffChan <- true
	if utils.GlobalObject.FakeData {
		c.Fo.Close()
	}

	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(id uint16, sim string, fragflag uint8, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	d, err := protocal.MakeJt808Msg(id, c.NextSeq, c.TerminalVer, fragflag, sim, data)
	if err != nil{
		return nil
	}
	//写回客户端
	n, err := c.Conn.Write(d)
	if err != nil {
		fmt.Println("Send Data error:, ", err, " Conn Writer exit")
		return err
	}
	c.NextSeq += 1

	fmt.Printf("Send data succ! data = %d\n", n)

	return nil
}

func (c *Connection) SendData(data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	//写回客户端
	c.msgChan <- data

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) GetRespResult(id uint16, seq uint16) (byte, bool){
	if v, ok := c.responseResult[id];ok{
		if v1, ok := v[seq]; ok{
			if c.TcpTimeout > 0 {
				if time.Now().UTC().Unix() - v1.t > c.TcpTimeout{
					return protocal.Failed, true
				}
			}

			return v1.result, true
		}
	}

	return 0, false
}

func (c *Connection) SetRespResult(id uint16, seq uint16, r byte) {
	if _, ok := c.responseResult[id];!ok{
		c.responseResult[id] = make(map[uint16]RespResult)
	}

	c.responseResult[id][seq] = RespResult{result:r, t:time.Now().UTC().Unix()}
}