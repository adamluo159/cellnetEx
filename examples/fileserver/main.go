package main

import (
	"flag"
	"fmt"
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	_ "github.com/adamluo159/cellnetEx/peer/http"
	"github.com/adamluo159/cellnetEx/proc"
	_ "github.com/adamluo159/cellnetEx/proc/http"
)

var shareDir = flag.String("share", ".", "folder to share")
var port = flag.Int("port", 9091, "listen port")

func main() {

	flag.Parse()

	queue := cellnetEx.NewEventQueue()

	p := peer.NewGenericPeer("http.Acceptor", "httpfile", fmt.Sprintf(":%d", *port), nil).(cellnetEx.HTTPAcceptor)
	p.SetFileServe(".", *shareDir)

	proc.BindProcessorHandler(p, "http", nil)

	p.Start()
	queue.StartLoop()

	queue.Wait()
}
