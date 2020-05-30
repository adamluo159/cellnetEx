package db

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/peer"
	_ "github.com/adamluo159/cellnetEx/peer/redix"
	"github.com/mediocregopher/radix.v2/redis"
	"testing"
	"time"
)

func prepareOperator(t *testing.T) cellnetEx.RedisPoolOperator {
	p := peer.NewGenericPeer("redix.Connector", "redis", "127.0.0.1:6379", nil)
	p.(cellnetEx.RedisConnector).SetConnectionCount(1)
	p.Start()

	for i := 0; i < 5; i++ {

		if p.(cellnetEx.PeerReadyChecker).IsReady() {
			break
		}

		time.Sleep(time.Millisecond * 200)
	}

	if !p.(cellnetEx.PeerReadyChecker).IsReady() {
		t.FailNow()
	}

	return p.(cellnetEx.RedisPoolOperator)
}

func TestRedix(t *testing.T) {

	prepareOperator(t).Operate(func(rawClient interface{}) interface{} {

		client := rawClient.(*redis.Client)
		client.Cmd("SET", "mykey", "myvalue")

		v, err := client.Cmd("GET", "mykey").Str()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if v != "myvalue" {
			t.FailNow()
		}

		client.Cmd("DEL", "mykey")

		return nil
	})

}
