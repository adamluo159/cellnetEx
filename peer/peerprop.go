package peer

import "github.com/adamluo159/cellnetEx"

type CorePeerProperty struct {
	name  string
	queue cellnetEx.EventQueue
	addr  string
}

// 获取通讯端的名称
func (self *CorePeerProperty) Name() string {
	return self.name
}

// 获取队列
func (self *CorePeerProperty) Queue() cellnetEx.EventQueue {
	return self.queue
}

// 获取SetAddress中的侦听或者连接地址
func (self *CorePeerProperty) Address() string {

	return self.addr
}

func (self *CorePeerProperty) SetName(v string) {
	self.name = v
}

func (self *CorePeerProperty) SetQueue(v cellnetEx.EventQueue) {
	self.queue = v
}

func (self *CorePeerProperty) SetAddress(v string) {
	self.addr = v
}
