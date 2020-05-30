package tests

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	"sync"
	"testing"
)

func TestContextSet(t *testing.T) {

	p := peer.NewPeer("tcp.Acceptor").(cellnetEx.TCPAcceptor)
	p.(cellnetEx.ContextSet).SetContext("sd", nil)

	if v, ok := p.(cellnetEx.ContextSet).GetContext("sd"); ok && v == nil {

	} else {
		t.FailNow()
	}

	var connMap = new(sync.Map)
	if p.(cellnetEx.ContextSet).FetchContext("sd", &connMap) && connMap == nil {

	} else {
		t.FailNow()
	}
}

func TestAutoAllocPort(t *testing.T) {

	p := peer.NewGenericPeer("tcp.Acceptor", "autoacc", ":0", nil)
	p.Start()

	t.Log("auto alloc port:", p.(cellnetEx.TCPAcceptor).Port())
}
