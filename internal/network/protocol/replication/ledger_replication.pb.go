// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.9
// source: ledger_replication.proto

package replication

import (
	block "github.com/alphabill-org/alphabill/internal/block"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LedgerReplicationResponse_Status int32

const (
	LedgerReplicationResponse_OK                         LedgerReplicationResponse_Status = 0
	LedgerReplicationResponse_INVALID_REQUEST_PARAMETERS LedgerReplicationResponse_Status = 1
	LedgerReplicationResponse_UNKNOWN_SYSTEM_IDENTIFIER  LedgerReplicationResponse_Status = 2
	LedgerReplicationResponse_BLOCKS_NOT_FOUND           LedgerReplicationResponse_Status = 3
	LedgerReplicationResponse_UNKNOWN                    LedgerReplicationResponse_Status = 4
)

// Enum value maps for LedgerReplicationResponse_Status.
var (
	LedgerReplicationResponse_Status_name = map[int32]string{
		0: "OK",
		1: "INVALID_REQUEST_PARAMETERS",
		2: "UNKNOWN_SYSTEM_IDENTIFIER",
		3: "BLOCKS_NOT_FOUND",
		4: "UNKNOWN",
	}
	LedgerReplicationResponse_Status_value = map[string]int32{
		"OK":                         0,
		"INVALID_REQUEST_PARAMETERS": 1,
		"UNKNOWN_SYSTEM_IDENTIFIER":  2,
		"BLOCKS_NOT_FOUND":           3,
		"UNKNOWN":                    4,
	}
)

func (x LedgerReplicationResponse_Status) Enum() *LedgerReplicationResponse_Status {
	p := new(LedgerReplicationResponse_Status)
	*p = x
	return p
}

func (x LedgerReplicationResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LedgerReplicationResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_ledger_replication_proto_enumTypes[0].Descriptor()
}

func (LedgerReplicationResponse_Status) Type() protoreflect.EnumType {
	return &file_ledger_replication_proto_enumTypes[0]
}

func (x LedgerReplicationResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LedgerReplicationResponse_Status.Descriptor instead.
func (LedgerReplicationResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_ledger_replication_proto_rawDescGZIP(), []int{1, 0}
}

type LedgerReplicationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SystemIdentifier []byte `protobuf:"bytes,1,opt,name=system_identifier,json=systemIdentifier,proto3" json:"system_identifier,omitempty"`
	NodeIdentifier   string `protobuf:"bytes,2,opt,name=node_identifier,json=nodeIdentifier,proto3" json:"node_identifier,omitempty"`
	BeginBlockNumber uint64 `protobuf:"varint,3,opt,name=begin_block_number,json=beginBlockNumber,proto3" json:"begin_block_number,omitempty"`
	EndBlockNumber   uint64 `protobuf:"varint,4,opt,name=end_block_number,json=endBlockNumber,proto3" json:"end_block_number,omitempty"`
}

func (x *LedgerReplicationRequest) Reset() {
	*x = LedgerReplicationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ledger_replication_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedgerReplicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedgerReplicationRequest) ProtoMessage() {}

func (x *LedgerReplicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_replication_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedgerReplicationRequest.ProtoReflect.Descriptor instead.
func (*LedgerReplicationRequest) Descriptor() ([]byte, []int) {
	return file_ledger_replication_proto_rawDescGZIP(), []int{0}
}

func (x *LedgerReplicationRequest) GetSystemIdentifier() []byte {
	if x != nil {
		return x.SystemIdentifier
	}
	return nil
}

func (x *LedgerReplicationRequest) GetNodeIdentifier() string {
	if x != nil {
		return x.NodeIdentifier
	}
	return ""
}

func (x *LedgerReplicationRequest) GetBeginBlockNumber() uint64 {
	if x != nil {
		return x.BeginBlockNumber
	}
	return 0
}

func (x *LedgerReplicationRequest) GetEndBlockNumber() uint64 {
	if x != nil {
		return x.EndBlockNumber
	}
	return 0
}

type LedgerReplicationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  LedgerReplicationResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=LedgerReplicationResponse_Status" json:"status,omitempty"`
	Message string                           `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Blocks  []*block.Block                   `protobuf:"bytes,3,rep,name=blocks,proto3" json:"blocks,omitempty"`
}

func (x *LedgerReplicationResponse) Reset() {
	*x = LedgerReplicationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ledger_replication_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedgerReplicationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedgerReplicationResponse) ProtoMessage() {}

func (x *LedgerReplicationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_replication_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedgerReplicationResponse.ProtoReflect.Descriptor instead.
func (*LedgerReplicationResponse) Descriptor() ([]byte, []int) {
	return file_ledger_replication_proto_rawDescGZIP(), []int{1}
}

func (x *LedgerReplicationResponse) GetStatus() LedgerReplicationResponse_Status {
	if x != nil {
		return x.Status
	}
	return LedgerReplicationResponse_OK
}

func (x *LedgerReplicationResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *LedgerReplicationResponse) GetBlocks() []*block.Block {
	if x != nil {
		return x.Blocks
	}
	return nil
}

var File_ledger_replication_proto protoreflect.FileDescriptor

var file_ledger_replication_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8, 0x01, 0x0a, 0x18, 0x4c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x11, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x10, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x12, 0x62, 0x65,
	0x67, 0x69, 0x6e, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x62, 0x65, 0x67, 0x69, 0x6e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x10, 0x65, 0x6e, 0x64, 0x5f,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0e, 0x65, 0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x22, 0x84, 0x02, 0x0a, 0x19, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x39, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x21, 0x2e, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x06, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x22, 0x72, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x49, 0x4e, 0x56, 0x41, 0x4c,
	0x49, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d,
	0x45, 0x54, 0x45, 0x52, 0x53, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x55, 0x4e, 0x4b, 0x4e, 0x4f,
	0x57, 0x4e, 0x5f, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x49,
	0x46, 0x49, 0x45, 0x52, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x53,
	0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x04, 0x42, 0x56, 0x5a, 0x54, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c,
	0x6c, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3b, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ledger_replication_proto_rawDescOnce sync.Once
	file_ledger_replication_proto_rawDescData = file_ledger_replication_proto_rawDesc
)

func file_ledger_replication_proto_rawDescGZIP() []byte {
	file_ledger_replication_proto_rawDescOnce.Do(func() {
		file_ledger_replication_proto_rawDescData = protoimpl.X.CompressGZIP(file_ledger_replication_proto_rawDescData)
	})
	return file_ledger_replication_proto_rawDescData
}

var file_ledger_replication_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_ledger_replication_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ledger_replication_proto_goTypes = []interface{}{
	(LedgerReplicationResponse_Status)(0), // 0: LedgerReplicationResponse.Status
	(*LedgerReplicationRequest)(nil),      // 1: LedgerReplicationRequest
	(*LedgerReplicationResponse)(nil),     // 2: LedgerReplicationResponse
	(*block.Block)(nil),                   // 3: Block
}
var file_ledger_replication_proto_depIdxs = []int32{
	0, // 0: LedgerReplicationResponse.status:type_name -> LedgerReplicationResponse.Status
	3, // 1: LedgerReplicationResponse.blocks:type_name -> Block
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_ledger_replication_proto_init() }
func file_ledger_replication_proto_init() {
	if File_ledger_replication_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ledger_replication_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LedgerReplicationRequest); i {
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
		file_ledger_replication_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LedgerReplicationResponse); i {
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
			RawDescriptor: file_ledger_replication_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ledger_replication_proto_goTypes,
		DependencyIndexes: file_ledger_replication_proto_depIdxs,
		EnumInfos:         file_ledger_replication_proto_enumTypes,
		MessageInfos:      file_ledger_replication_proto_msgTypes,
	}.Build()
	File_ledger_replication_proto = out.File
	file_ledger_replication_proto_rawDesc = nil
	file_ledger_replication_proto_goTypes = nil
	file_ledger_replication_proto_depIdxs = nil
}
