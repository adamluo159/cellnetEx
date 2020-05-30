package tests

import (
	"fmt"
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	_ "github.com/adamluo159/cellnetEx/peer/gorillaws"
	_ "github.com/adamluo159/cellnetEx/peer/tcp"
	_ "github.com/adamluo159/cellnetEx/peer/udp"
	"github.com/adamluo159/cellnetEx/proc"
	_ "github.com/adamluo159/cellnetEx/proc/gorillaws"
	_ "github.com/adamluo159/cellnetEx/proc/tcp"
	_ "github.com/adamluo159/cellnetEx/proc/udp"
	"testing"
	"time"
)

type echoContext struct {
	Address   string
	Protocol  string
	Processor string
	Tester    *SignalTester
	Acceptor  cellnetEx.GenericPeer
}

var (
	echoContexts = []*echoContext{
		{
			Address:   "127.0.0.1:7701",
			Protocol:  "tcp",
			Processor: "tcp.ltv",
		},
		{
			Address:   "127.0.0.1:7702",
			Protocol:  "udp",
			Processor: "udp.ltv",
		},

		{
			Address:   "127.0.0.1:7703",
			Protocol:  "gorillaws",
			Processor: "gorillaws.ltv",
		},
	}
)

func echo_StartServer(context *echoContext) {
	queue := cellnetEx.NewEventQueue()

	context.Acceptor = peer.NewGenericPeer(context.Protocol+".Acceptor", context.Protocol+"server", context.Address, queue)

	proc.BindProcessorHandler(context.Acceptor, context.Processor, func(ev cellnetEx.Event) {

		switch msg := ev.Message().(type) {
		case *cellnetEx.SessionAccepted:
			fmt.Println("server accepted")
		case *TestEchoACK:

			fmt.Printf("server recv %+v\n", msg)

			ev.Session().Send(&TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			})

		case *cellnetEx.SessionClosed:
			fmt.Println("session closed: ", ev.Session().ID())
		}

	})

	context.Acceptor.Start()

	queue.StartLoop()
}

func echo_StartClient(echoContext *echoContext) {
	queue := cellnetEx.NewEventQueue()

	p := peer.NewGenericPeer(echoContext.Protocol+".Connector", echoContext.Protocol+"client", echoContext.Address, queue)

	proc.BindProcessorHandler(p, echoContext.Processor, func(ev cellnetEx.Event) {

		switch msg := ev.Message().(type) {
		case *cellnetEx.SessionConnected:
			fmt.Println("client connected")
			ev.Session().Send(&TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *TestEchoACK:

			fmt.Printf("client recv %+v\n", msg)

			echoContext.Tester.Done(1)

		case *cellnetEx.SessionClosed:
			fmt.Println("client closed")
		}
	})

	p.Start()

	queue.StartLoop()

	echoContext.Tester.WaitAndExpect("not recv data", 1)
}

func runEcho(t *testing.T, index int) {

	ctx := echoContexts[index]

	ctx.Tester = NewSignalTester(t)
	ctx.Tester.SetTimeout(time.Hour)

	echo_StartServer(ctx)

	echo_StartClient(ctx)

	ctx.Acceptor.Stop()
}

func TestEchoTCP(t *testing.T) {

	runEcho(t, 0)
}

func TestEchoUDP(t *testing.T) {

	runEcho(t, 1)
}

func TestEchoWS(t *testing.T) {

	runEcho(t, 2)
}
