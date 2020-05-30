package tests

import (
	"fmt"
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/codec"
	_ "github.com/adamluo159/cellnetEx/codec/binary"
	"github.com/adamluo159/cellnetEx/util"
	"reflect"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnetEx.RegisterMessageMeta(&cellnetEx.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(util.StringHash("tests.TestEchoACK")),
	})
}
