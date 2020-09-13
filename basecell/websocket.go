package basecell

import (
	"time"

	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	_ "github.com/adamluo159/cellnetEx/peer/gorillaws"
	"github.com/adamluo159/cellnetEx/proc"
	_ "github.com/adamluo159/cellnetEx/proc/gorillaws"
)

//NewWsServer 创建ws服务器
func (bcell *BaseCell) NewWsServer(addr string) {
	bcell.peer = peer.NewGenericPeer("gorillaws.Acceptor", "", addr, bcell.queue)
	proc.BindProcessorHandler(bcell.peer, "gorillaws.ltv", bcell.msgQueue())
}

//NewWsClient 创建ws客户端
func (bcell *BaseCell) NewWsClient(addr string) {
	bcell.peer = peer.NewGenericPeer("gorillaws.Connector", "client", addr, bcell.queue)
	bcell.peer.(cellnetEx.WSConnector).SetReconnectDuration(time.Second * 5)
	proc.BindProcessorHandler(bcell.peer, "gorillaws.ltv", bcell.msgQueue())
}
