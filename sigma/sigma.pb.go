package sigma

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParentID uint64   `protobuf:"varint,1,opt,name=ParentID,proto3" json:"ParentID,omitempty"`
	Body     []byte   `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
	Data     string   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	Entities []uint64 `protobuf:"varint,4,rep,packed,name=Entities,proto3" json:"Entities,omitempty"`
}

func (e Event) ProtoReflect() protoreflect.Message {
	//TODO implement me
	panic("implement me")
}
