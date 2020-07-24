package jtnet

import (
	"github.com/dushu1105/jt808/protocal"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]protocal.Handler //存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]protocal.Handler),
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(msg *protocal.JT808Msg) (*protocal.Jt808ResultMsg, error) {
	requestId := uint32(msg.Header.Id)
	if msg.Header.Ver == 1 {
		requestId |= 0x10000
	}

	handler, ok := mh.Apis[requestId]
	if !ok {
		msg.Printf("没有 0x%x 对应处理函数\n", msg.Header.Id)
		//panic("not found")
		return nil, nil
	} else {
		msg.Printf("处理 0x%x\n", msg.Header.Id)
	}

	//执行对应处理方法
	return handler.Do(msg)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddHandler(msgId uint32, h protocal.Handler) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = h
	fmt.Printf("Add api msgId = 0x%x\n", msgId)
}