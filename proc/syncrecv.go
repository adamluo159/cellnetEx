package proc

import (
	"github.com/adamluo159/cellnetEx"
	"reflect"
	"sync"
)

// 同步接收消息器, 可选件，可作为流程测试辅助工具
type SyncReceiver struct {
	evChan chan cellnetEx.Event

	callback func(ev cellnetEx.Event)
}

// 将处理回调返回给BindProcessorHandler用于注册
func (self *SyncReceiver) EventCallback() cellnetEx.EventCallback {

	return self.callback
}

// 持续阻塞，直到某个消息到达后，使用回调返回消息
func (self *SyncReceiver) Recv(callback cellnetEx.EventCallback) *SyncReceiver {
	callback(<-self.evChan)
	return self
}

// 持续阻塞，直到某个消息到达后，返回消息
func (self *SyncReceiver) WaitMessage(msgName string) (msg interface{}) {

	var wg sync.WaitGroup

	meta := cellnetEx.MessageMetaByFullName(msgName)
	if meta == nil {
		panic("unknown message name:" + msgName)
	}

	wg.Add(1)

	self.Recv(func(ev cellnetEx.Event) {

		inMeta := cellnetEx.MessageMetaByType(reflect.TypeOf(ev.Message()))
		if inMeta == meta {
			msg = ev.Message()
			wg.Done()
		}

	})

	wg.Wait()
	return
}

func NewSyncReceiver(p cellnetEx.Peer) *SyncReceiver {

	self := &SyncReceiver{
		evChan: make(chan cellnetEx.Event),
	}

	self.callback = func(ev cellnetEx.Event) {

		self.evChan <- ev
	}

	return self
}
