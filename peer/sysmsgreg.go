package peer

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/codec"
	_ "github.com/adamluo159/cellnetEx/codec/binary"
	"github.com/adamluo159/cellnetEx/util"
	"reflect"
)

func init() {
	cellnetEx.RegisterMessageMeta(&cellnetEx.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*cellnetEx.SessionAccepted)(nil)).Elem(),
		ID:    int(util.StringHash("cellnetEx.SessionAccepted")),
	})
	cellnetEx.RegisterMessageMeta(&cellnetEx.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*cellnetEx.SessionConnected)(nil)).Elem(),
		ID:    int(util.StringHash("cellnetEx.SessionConnected")),
	})
	cellnetEx.RegisterMessageMeta(&cellnetEx.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*cellnetEx.SessionConnectError)(nil)).Elem(),
		ID:    int(util.StringHash("cellnetEx.SessionConnectError")),
	})
	cellnetEx.RegisterMessageMeta(&cellnetEx.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*cellnetEx.SessionClosed)(nil)).Elem(),
		ID:    int(util.StringHash("cellnetEx.SessionClosed")),
	})
	cellnetEx.RegisterMessageMeta(&cellnetEx.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*cellnetEx.SessionCloseNotify)(nil)).Elem(),
		ID:    int(util.StringHash("cellnetEx.SessionCloseNotify")),
	})
	cellnetEx.RegisterMessageMeta(&cellnetEx.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*cellnetEx.SessionInit)(nil)).Elem(),
		ID:    int(util.StringHash("cellnetEx.SessionInit")),
	})
}
