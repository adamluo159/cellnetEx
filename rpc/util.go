package rpc

import (
	"errors"
	"github.com/adamluo159/cellnetEx"
)

var (
	ErrInvalidPeerSession = errors.New("rpc: Invalid peer type, require cellnetEx.RPCSessionGetter or cellnetEx.Session")
	ErrEmptySession       = errors.New("rpc: Empty session")
)

type RPCSessionGetter interface {
	RPCSession() cellnetEx.Session
}

// 从peer获取rpc使用的session
func getPeerSession(ud interface{}) (ses cellnetEx.Session, err error) {

	if ud == nil {
		return nil, ErrInvalidPeerSession
	}

	switch i := ud.(type) {
	case RPCSessionGetter:
		ses = i.RPCSession()
	case cellnetEx.Session:
		ses = i
	case cellnetEx.TCPConnector:
		ses = i.Session()
	default:
		err = ErrInvalidPeerSession
		return
	}

	if ses == nil {
		return nil, ErrEmptySession
	}

	return
}
