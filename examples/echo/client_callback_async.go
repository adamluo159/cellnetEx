package main

import (
	"fmt"
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	"github.com/adamluo159/cellnetEx/proc"
)

func clientAsyncCallback() {

	// 等待服务器返回数据
	done := make(chan struct{})

	queue := cellnetEx.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "clientAsyncCallback", peerAddress, queue)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnetEx.Event) {

		switch msg := ev.Message().(type) {
		case *cellnetEx.SessionConnected: // 已经连接上
			fmt.Println("clientAsyncCallback connected")
			ev.Session().Send(&TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *TestEchoACK: //收到服务器发送的消息

			fmt.Printf("clientAsyncCallback recv %+v\n", msg)

			// 完成操作
			done <- struct{}{}

		case *cellnetEx.SessionClosed:
			fmt.Println("clientAsyncCallback closed")
		}
	})

	p.Start()

	queue.StartLoop()

	// 等待客户端收到消息
	<-done
}
