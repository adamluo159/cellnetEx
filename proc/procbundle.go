package proc

import (
	"github.com/adamluo159/cellnetEx"
)

// 处理器设置接口，由各Peer实现
type ProcessorBundle interface {

	// 设置 传输器，负责收发消息
	SetTransmitter(v cellnetEx.MessageTransmitter)

	// 设置 接收后，发送前的事件处理流程
	SetHooker(v cellnetEx.EventHooker)

	// 设置 接收后最终处理回调
	SetCallback(v cellnetEx.EventCallback)
}

// 让EventCallback保证放在ses的队列里，而不是并发的
func NewQueuedEventCallback(callback cellnetEx.EventCallback) cellnetEx.EventCallback {

	return func(ev cellnetEx.Event) {
		if callback != nil {
			cellnetEx.SessionQueuedCall(ev.Session(), func() {

				callback(ev)
			})
		}
	}
}

// 当需要多个Hooker时，使用NewMultiHooker将多个hooker合并成1个hooker处理
type MultiHooker []cellnetEx.EventHooker

func (self MultiHooker) OnInboundEvent(input cellnetEx.Event) (output cellnetEx.Event) {

	for _, h := range self {

		input = h.OnInboundEvent(input)

		if input == nil {
			break
		}
	}

	return input
}

func (self MultiHooker) OnOutboundEvent(input cellnetEx.Event) (output cellnetEx.Event) {

	for _, h := range self {

		input = h.OnOutboundEvent(input)

		if input == nil {
			break
		}
	}

	return input
}

func NewMultiHooker(h ...cellnetEx.EventHooker) cellnetEx.EventHooker {

	return MultiHooker(h)
}
