package msglog

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/davyxu/golog"
)

// 萃取消息中的消息
type PacketMessagePeeker interface {
	Message() interface{}
}

type UserData interface {
	UID() string
}

func WriteRecvLogger(log *golog.Logger, protocol string, ses cellnetEx.Session, msg interface{}) {

	if log.IsDebugEnabled() {

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(cellnetEx.MessageToID(msg)) {
			peerInfo := ses.Peer().(cellnetEx.PeerProperty)
			uid := ""
			if ses.GetUserData() != nil && ses.GetUserData().(UserData) != nil {
				uid = ses.GetUserData().(UserData).UID()
			}
			log.Debugf("#%s.recv(%s)@%d uid:%s len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ses.ID(),
				uid,
				cellnetEx.MessageSize(msg),
				cellnetEx.MessageToName(msg),
				cellnetEx.MessageToString(msg))
		}

	}
}

func WriteSendLogger(log *golog.Logger, protocol string, ses cellnetEx.Session, msg interface{}) {

	if log.IsDebugEnabled() {

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(cellnetEx.MessageToID(msg)) {
			peerInfo := ses.Peer().(cellnetEx.PeerProperty)

			log.Debugf("#%s.send(%s)@%d len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ses.ID(),
				cellnetEx.MessageSize(msg),
				cellnetEx.MessageToName(msg),
				cellnetEx.MessageToString(msg))
		}

	}

}
