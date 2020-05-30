package gogopb

import (
	"github.com/adamluo159/cellnetEx"
	"github.com/adamluo159/cellnetEx/codec"
	"github.com/gogo/protobuf/proto"
)

type gogopbCodec struct {
}

// 编码器的名称
func (self *gogopbCodec) Name() string {
	return "gogopb"
}

func (self *gogopbCodec) MimeType() string {
	return "application/x-protobuf"
}

func (self *gogopbCodec) Encode(msgObj interface{}, ctx cellnetEx.ContextSet) (data interface{}, err error) {

	return proto.Marshal(msgObj.(proto.Message))

}

func (self *gogopbCodec) Decode(data interface{}, msgObj interface{}) error {

	return proto.Unmarshal(data.([]byte), msgObj.(proto.Message))
}

func init() {

	// 注册编码器
	codec.RegisterCodec(new(gogopbCodec))
}
