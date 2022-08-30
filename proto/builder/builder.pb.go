// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.15.8
// source: proto/builder/builder.proto

package builder

import (
	reflect "reflect"
	sync "sync"

	github_com_prysmaticlabs_prysm_v3_consensus_types_primitives "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	_ "github.com/prysmaticlabs/prysm/v3/proto/eth/ext"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BuilderPayloadAttributes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp  uint64                                                            `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	PrevRandao []byte                                                            `protobuf:"bytes,2,opt,name=prev_randao,json=prevRandao,proto3" json:"prev_randao,omitempty" ssz-size:"32"`
	Slot       github_com_prysmaticlabs_prysm_v3_consensus_types_primitives.Slot `protobuf:"varint,3,opt,name=slot,proto3" json:"slot,omitempty" cast-type:"github.com/prysmaticlabs/prysm/v3/consensus-types/primitives.Slot"`
	BlockHash  []byte                                                            `protobuf:"bytes,4,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty" ssz-size:"32"`
}

func (x *BuilderPayloadAttributes) Reset() {
	*x = BuilderPayloadAttributes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_builder_builder_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuilderPayloadAttributes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuilderPayloadAttributes) ProtoMessage() {}

func (x *BuilderPayloadAttributes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_builder_builder_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuilderPayloadAttributes.ProtoReflect.Descriptor instead.
func (*BuilderPayloadAttributes) Descriptor() ([]byte, []int) {
	return file_proto_builder_builder_proto_rawDescGZIP(), []int{0}
}

func (x *BuilderPayloadAttributes) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *BuilderPayloadAttributes) GetPrevRandao() []byte {
	if x != nil {
		return x.PrevRandao
	}
	return nil
}

func (x *BuilderPayloadAttributes) GetSlot() github_com_prysmaticlabs_prysm_v3_consensus_types_primitives.Slot {
	if x != nil {
		return x.Slot
	}
	return github_com_prysmaticlabs_prysm_v3_consensus_types_primitives.Slot(0)
}

func (x *BuilderPayloadAttributes) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

var File_proto_builder_builder_proto protoreflect.FileDescriptor

var file_proto_builder_builder_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x2f,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x65,
	0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x1a,
	0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x65, 0x78, 0x74, 0x2f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe3, 0x01, 0x0a,
	0x18, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x27, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x5f,
	0x72, 0x61, 0x6e, 0x64, 0x61, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5,
	0x18, 0x02, 0x33, 0x32, 0x52, 0x0a, 0x70, 0x72, 0x65, 0x76, 0x52, 0x61, 0x6e, 0x64, 0x61, 0x6f,
	0x12, 0x59, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x45,
	0x82, 0xb5, 0x18, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70,
	0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x79,
	0x73, 0x6d, 0x2f, 0x76, 0x33, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2d,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x73,
	0x2e, 0x53, 0x6c, 0x6f, 0x74, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x12, 0x25, 0x0a, 0x0a, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x42,
	0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x42, 0x85, 0x01, 0x0a, 0x14, 0x6f, 0x72, 0x67, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72,
	0x65, 0x75, 0x6d, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x42, 0x0c, 0x42, 0x75, 0x69,
	0x6c, 0x64, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x37, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69,
	0x63, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x76, 0x33, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x3b, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x65, 0x72, 0xaa, 0x02, 0x10, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0xca, 0x02, 0x10, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65,
	0x75, 0x6d, 0x5c, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_builder_builder_proto_rawDescOnce sync.Once
	file_proto_builder_builder_proto_rawDescData = file_proto_builder_builder_proto_rawDesc
)

func file_proto_builder_builder_proto_rawDescGZIP() []byte {
	file_proto_builder_builder_proto_rawDescOnce.Do(func() {
		file_proto_builder_builder_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_builder_builder_proto_rawDescData)
	})
	return file_proto_builder_builder_proto_rawDescData
}

var file_proto_builder_builder_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_builder_builder_proto_goTypes = []interface{}{
	(*BuilderPayloadAttributes)(nil), // 0: ethereum.builder.BuilderPayloadAttributes
}
var file_proto_builder_builder_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_builder_builder_proto_init() }
func file_proto_builder_builder_proto_init() {
	if File_proto_builder_builder_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_builder_builder_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuilderPayloadAttributes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_builder_builder_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_builder_builder_proto_goTypes,
		DependencyIndexes: file_proto_builder_builder_proto_depIdxs,
		MessageInfos:      file_proto_builder_builder_proto_msgTypes,
	}.Build()
	File_proto_builder_builder_proto = out.File
	file_proto_builder_builder_proto_rawDesc = nil
	file_proto_builder_builder_proto_goTypes = nil
	file_proto_builder_builder_proto_depIdxs = nil
}
