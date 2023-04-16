package sigma

import "google.golang.org/protobuf/runtime/protoimpl"

type Configuration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NatsUrl      string `protobuf:"bytes,1,opt,name=NatsUrl,proto3" json:"NatsUrl,omitempty"`
	NatsUsername string `protobuf:"bytes,2,opt,name=NatsUsername,proto3" json:"NatsUsername,omitempty"`
	NatsPassword string `protobuf:"bytes,3,opt,name=NatsPassword,proto3" json:"NatsPassword,omitempty"`
}
