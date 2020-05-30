package gorillaws

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/msglog"
)

// 带有RPC和relay功能
type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnetEx.Event) (outputEvent cellnetEx.Event) {

	msglog.WriteRecvLogger(log, "ws", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnetEx.Event) (outputEvent cellnetEx.Event) {

	msglog.WriteSendLogger(log, "ws", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}
