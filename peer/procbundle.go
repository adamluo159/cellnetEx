package peer

import (
	"errors"
	"github.com/adamluo159/cellnetEx"
)

// 手动投递消息， 兼容v2的设计
type MessagePoster interface {

	// 投递一个消息到Hooker之前
	ProcEvent(ev cellnetEx.Event)
}

type CoreProcBundle struct {
	transmit cellnetEx.MessageTransmitter
	hooker   cellnetEx.EventHooker
	callback cellnetEx.EventCallback
}

func (self *CoreProcBundle) GetBundle() *CoreProcBundle {
	return self
}

func (self *CoreProcBundle) SetTransmitter(v cellnetEx.MessageTransmitter) {
	self.transmit = v
}

func (self *CoreProcBundle) SetHooker(v cellnetEx.EventHooker) {
	self.hooker = v
}

func (self *CoreProcBundle) SetCallback(v cellnetEx.EventCallback) {
	self.callback = v
}

var notHandled = errors.New("Processor: Transimitter nil")

func (self *CoreProcBundle) ReadMessage(ses cellnetEx.Session) (msg interface{}, err error) {

	if self.transmit != nil {
		return self.transmit.OnRecvMessage(ses)
	}

	return nil, notHandled
}

func (self *CoreProcBundle) SendMessage(ev cellnetEx.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnOutboundEvent(ev)
	}

	if self.transmit != nil && ev != nil {
		self.transmit.OnSendMessage(ev.Session(), ev.Message())
	}
}

func (self *CoreProcBundle) ProcEvent(ev cellnetEx.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnInboundEvent(ev)
	}

	if self.callback != nil && ev != nil {
		self.callback(ev)
	}
}
