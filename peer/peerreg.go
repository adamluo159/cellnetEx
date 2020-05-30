package peer

import (
	"fmt"
	"github.com/adamluo159/cellnetEx"
	"sort"
)

type PeerCreateFunc func() cellnetEx.Peer

var peerByName = map[string]PeerCreateFunc{}

// 注册Peer创建器
func RegisterPeerCreator(f PeerCreateFunc) {

	// 临时实例化一个，获取类型
	dummyPeer := f()

	if _, ok := peerByName[dummyPeer.TypeName()]; ok {
		panic("duplicate peer type: " + dummyPeer.TypeName())
	}

	peerByName[dummyPeer.TypeName()] = f
}

// Peer创建器列表
func PeerCreatorList() (ret []string) {

	for name := range peerByName {
		ret = append(ret, name)
	}

	sort.Strings(ret)
	return
}

// cellnet自带的peer对应包
func getPackageByPeerName(name string) string {
	switch name {
	case "tcp.Connector", "tcp.Acceptor", "tcp.SyncConnector":
		return "github.com/adamluo159/cellnetEx/peer/tcp"
	case "udp.Connector", "udp.Acceptor":
		return "github.com/adamluo159/cellnetEx/peer/udp"
	case "gorillaws.Acceptor", "gorillaws.Connector", "gorillaws.SyncConnector":
		return "github.com/adamluo159/cellnetEx/peer/gorillaws"
	case "http.Connector", "http.Acceptor":
		return "github.com/adamluo159/cellnetEx/peer/http"
	case "redix.Connector":
		return "github.com/adamluo159/cellnetEx/peer/redix"
	case "mysql.Connector":
		return "github.com/adamluo159/cellnetEx/peer/mysql"
	default:
		return "package/to/your/peer"
	}
}

// 创建一个Peer
func NewPeer(peerType string) cellnetEx.Peer {
	peerCreator := peerByName[peerType]
	if peerCreator == nil {
		panic(fmt.Sprintf("peer type not found '%s'\ntry to add code below:\nimport (\n  _ \"%s\"\n)\n\n",
			peerType,
			getPackageByPeerName(peerType)))
	}

	return peerCreator()
}

// 创建Peer后，设置基本属性
func NewGenericPeer(peerType, name, addr string, q cellnetEx.EventQueue) cellnetEx.GenericPeer {

	p := NewPeer(peerType)
	gp := p.(cellnetEx.GenericPeer)
	gp.SetName(name)
	gp.SetAddress(addr)
	gp.SetQueue(q)
	return gp
}
