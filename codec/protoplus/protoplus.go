package protoplus

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/codec"
	"github.com/davyxu/protoplus/proto"
)

type protoplus struct {
}

func (self *protoplus) Name() string {
	return "protoplus"
}

func (self *protoplus) MimeType() string {
	return "application/binary"
}

func (self *protoplus) Encode(msgObj interface{}, ctx cellnetEx.ContextSet) (data interface{}, err error) {

	return proto.Marshal(msgObj)

}

func (self *protoplus) Decode(data interface{}, msgObj interface{}) error {

	return proto.Unmarshal(data.([]byte), msgObj)
}

func init() {

	codec.RegisterCodec(new(protoplus))
}
