package tcp

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	"net"
	"time"
)

type tcpSyncConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	defaultSes *tcpSession
}

func (self *tcpSyncConnector) Port() int {
	conn := self.defaultSes.Conn()

	if conn == nil {
		return 0
	}

	return conn.LocalAddr().(*net.TCPAddr).Port
}

func (self *tcpSyncConnector) Start() cellnetEx.Peer {

	// 尝试用Socket连接地址
	conn, err := net.Dial("tcp", self.Address())

	// 发生错误时退出
	if err != nil {

		log.Debugf("#tcp.connect failed(%s)@%d address: %s", self.Name(), self.defaultSes.ID(), self.Address())

		self.ProcEvent(&cellnetEx.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnetEx.SessionConnectError{}})
		return self
	}

	self.defaultSes.setConn(conn)

	self.ApplySocketOption(conn)

	self.defaultSes.Start()

	self.ProcEvent(&cellnetEx.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnetEx.SessionConnected{}})

	return self
}

func (self *tcpSyncConnector) Session() cellnetEx.Session {
	return self.defaultSes
}

func (self *tcpSyncConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *tcpSyncConnector) ReconnectDuration() time.Duration {
	return 0
}

func (self *tcpSyncConnector) SetReconnectDuration(v time.Duration) {

}

func (self *tcpSyncConnector) Stop() {

	if self.defaultSes != nil {
		self.defaultSes.Close()
	}

}

func (self *tcpSyncConnector) IsReady() bool {

	return self.SessionCount() != 0
}

func (self *tcpSyncConnector) TypeName() string {
	return "tcp.SyncConnector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnetEx.Peer {
		self := &tcpSyncConnector{
			SessionManager: new(peer.CoreSessionManager),
		}

		self.defaultSes = newSession(nil, self, nil)

		self.CoreTCPSocketOption.Init()

		return self
	})
}
