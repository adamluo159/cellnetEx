package rpc

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/codec"
)

type RecvMsgEvent struct {
	ses    cellnetEx.Session
	Msg    interface{}
	callid int64
}

func (self *RecvMsgEvent) Session() cellnetEx.Session {
	return self.ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *RecvMsgEvent) Queue() cellnetEx.EventQueue {
	return self.ses.Peer().(interface {
		Queue() cellnetEx.EventQueue
	}).Queue()
}

func (self *RecvMsgEvent) Reply(msg interface{}) {

	data, meta, err := codec.EncodeMessage(msg, nil)

	if err != nil {
		log.Errorf("rpc reply message encode error: %s", err)
		return
	}

	self.ses.Send(&RemoteCallACK{
		MsgID:  uint32(meta.ID),
		Data:   data,
		CallID: self.callid,
	})
}
