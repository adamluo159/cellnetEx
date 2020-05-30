package relay

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/codec"
	"github.com/adamluo159/cellnetEx/msglog"
)

type PassthroughContent struct {
	Int64      int64   // 透传int64
	Int64Slice []int64 // 透传int64切片
	Str        string
}

// 处理入站的relay消息
func ResoleveInboundEvent(inputEvent cellnetEx.Event) (ouputEvent cellnetEx.Event, handled bool, err error) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:

		ev := &RecvMsgEvent{
			Ses: inputEvent.Session(),
			ack: relayMsg,
		}

		if relayMsg.MsgID != 0 {

			ev.Msg, _, err = codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Msg)
			if err != nil {
				return
			}
		}

		if msglog.IsMsgLogValid(int(relayMsg.MsgID)) {

			peerInfo := inputEvent.Session().Peer().(cellnetEx.PeerProperty)

			log.Debugf("#relay.recv(%s)@%d len: %d %s {%s}| %s",
				peerInfo.Name(),
				inputEvent.Session().ID(),
				cellnetEx.MessageSize(ev.Message()),
				cellnetEx.MessageToName(ev.Message()),
				cellnetEx.MessageToString(relayMsg),
				cellnetEx.MessageToString(ev.Message()))
		}

		if bcFunc != nil {
			// 转到对应线程中调用
			cellnetEx.SessionQueuedCall(inputEvent.Session(), func() {
				bcFunc(ev)
			})
		}

		return ev, true, nil
	}

	return inputEvent, false, nil
}

// 处理relay.Relay出站消息的日志
func ResolveOutboundEvent(inputEvent cellnetEx.Event) (handled bool, err error) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:
		if msglog.IsMsgLogValid(int(relayMsg.MsgID)) {

			var payload interface{}
			if relayMsg.MsgID != 0 {

				payload, _, err = codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Msg)
				if err != nil {
					return
				}
			}

			peerInfo := inputEvent.Session().Peer().(cellnetEx.PeerProperty)

			log.Debugf("#relay.send(%s)@%d len: %d %s {%s}| %s",
				peerInfo.Name(),
				inputEvent.Session().ID(),
				cellnetEx.MessageSize(payload),
				cellnetEx.MessageToName(payload),
				cellnetEx.MessageToString(relayMsg),
				cellnetEx.MessageToString(payload))
		}

		return true, nil

	}

	return
}
