package gorillaws

import (
	"sync"

	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	"github.com/gorilla/websocket"
)

// Socket会话
type wsSession struct {
	peer.CoreContextSet
	peer.CoreSessionIdentify
	*peer.CoreProcBundle

	pInterface cellnetEx.Peer

	data interface{}

	conn *websocket.Conn

	// 退出同步器
	exitSync sync.WaitGroup

	// 发送队列
	sendQueue *cellnetEx.Pipe

	cleanupGuard sync.Mutex

	endNotify func()
}

func (self *wsSession) Peer() cellnetEx.Peer {
	return self.pInterface
}

// 取原始连接
func (self *wsSession) Raw() interface{} {
	if self.conn == nil {
		return nil
	}

	return self.conn
}

func (self *wsSession) Close() {
	self.sendQueue.Add(nil)
}

// 发送封包
func (self *wsSession) Send(msg interface{}) {
	self.sendQueue.Add(msg)
}

//SetUserData 设置用户数据
func (self *wsSession) SetUserData(data interface{}) {
	self.data = data
}

//GetUserData 获取用户数据
func (self *wsSession) GetUserData() interface{} {
	return self.data
}

// 接收循环
func (self *wsSession) recvLoop() {

	for self.conn != nil {

		msg, err := self.ReadMessage(self)

		if err != nil {

			log.Debugln(err)

			// if !util.IsEOFOrNetReadError(err) {
			// 	log.Errorln("session closed:", err)
			// }

			self.ProcEvent(&cellnetEx.RecvMsgEvent{Ses: self, Msg: &cellnetEx.SessionClosed{}})
			break
		}

		self.ProcEvent(&cellnetEx.RecvMsgEvent{Ses: self, Msg: msg})
	}

	self.Close()

	// 通知完成
	self.exitSync.Done()
}

// 发送循环
func (self *wsSession) sendLoop() {

	var writeList []interface{}

	for {
		writeList = writeList[0:0]
		exit := self.sendQueue.Pick(&writeList)

		// 遍历要发送的数据
		for _, msg := range writeList {

			// TODO SendMsgEvent并不是很有意义
			self.SendMessage(&cellnetEx.SendMsgEvent{Ses: self, Msg: msg})
		}

		if exit {
			break
		}
	}

	// 关闭连接
	if self.conn != nil {
		self.conn.Close()
		self.conn = nil
	}

	// 通知完成
	self.exitSync.Done()
}

// 启动会话的各种资源
func (self *wsSession) Start() {

	// 将会话添加到管理器
	self.Peer().(peer.SessionManager).Add(self)

	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	go func() {
		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.Peer().(peer.SessionManager).Remove(self)

		if self.endNotify != nil {
			self.endNotify()
		}

	}()

	// 启动并发接收goroutine
	go self.recvLoop()

	// 启动并发发送goroutine
	go self.sendLoop()
}

func newSession(conn *websocket.Conn, p cellnetEx.Peer, endNotify func()) *wsSession {
	self := &wsSession{
		conn:       conn,
		endNotify:  endNotify,
		sendQueue:  cellnetEx.NewPipe(),
		pInterface: p,
		CoreProcBundle: p.(interface {
			GetBundle() *peer.CoreProcBundle
		}).GetBundle(),
	}

	return self
}
