package rpc

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/codec"
	"github.com/adamluo159/cellnetEx/msglog"
)

type RemoteCallMsg interface {
	GetMsgID() uint16
	GetMsgData() []byte
	GetCallID() int64
}

func ResolveInboundEvent(inputEvent cellnetEx.Event) (ouputEvent cellnetEx.Event, handled bool, err error) {

	if _, ok := inputEvent.(*RecvMsgEvent); ok {
		return inputEvent, false, nil
	}

	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return inputEvent, false, nil
	}

	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return inputEvent, false, err
	}

	if msglog.IsMsgLogValid(int(rpcMsg.GetMsgID())) {
		peerInfo := inputEvent.Session().Peer().(cellnetEx.PeerProperty)

		log.Debugf("#rpc.recv(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			cellnetEx.MessageSize(userMsg),
			cellnetEx.MessageToName(userMsg),
			cellnetEx.MessageToString(userMsg))
	}

	switch inputEvent.Message().(type) {
	case *RemoteCallREQ: // 服务端收到客户端的请求

		return &RecvMsgEvent{
			inputEvent.Session(),
			userMsg,
			rpcMsg.GetCallID(),
		}, true, nil

	case *RemoteCallACK: // 客户端收到服务器的回应
		request := getRequest(rpcMsg.GetCallID())
		if request != nil {
			request.RecvFeedback(userMsg)
		}

		return inputEvent, true, nil
	}

	return inputEvent, false, nil
}

func ResolveOutboundEvent(inputEvent cellnetEx.Event) (handled bool, err error) {
	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return false, nil
	}

	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return false, err
	}

	if msglog.IsMsgLogValid(int(rpcMsg.GetMsgID())) {
		peerInfo := inputEvent.Session().Peer().(cellnetEx.PeerProperty)

		log.Debugf("#rpc.send(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			cellnetEx.MessageSize(userMsg),
			cellnetEx.MessageToName(userMsg),
			cellnetEx.MessageToString(userMsg))
	}

	// 避免后续环节处理

	return true, nil
}
