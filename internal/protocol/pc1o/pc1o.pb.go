// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: pc1o.proto

package pc1o

import (
	certificates "gitdc.ee.guardtime.com/alphabill/alphabill/internal/certificates"
	transaction "gitdc.ee.guardtime.com/alphabill/alphabill/internal/transaction"
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

type PC1ORequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SystemIdentifier   []byte                           `protobuf:"bytes,1,opt,name=system_identifier,json=systemIdentifier,proto3" json:"system_identifier,omitempty"`
	NodeIdentifier     string                           `protobuf:"bytes,2,opt,name=node_identifier,json=nodeIdentifier,proto3" json:"node_identifier,omitempty"`
	UnicityCertificate *certificates.UnicityCertificate `protobuf:"bytes,3,opt,name=unicity_certificate,json=unicityCertificate,proto3" json:"unicity_certificate,omitempty"`
	Transactions       []*transaction.Transaction       `protobuf:"bytes,4,rep,name=transactions,proto3" json:"transactions,omitempty"`
	Signature          []byte                           `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *PC1ORequest) Reset() {
	*x = PC1ORequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pc1o_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PC1ORequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PC1ORequest) ProtoMessage() {}

func (x *PC1ORequest) ProtoReflect() protoreflect.Message {
	mi := &file_pc1o_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PC1ORequest.ProtoReflect.Descriptor instead.
func (*PC1ORequest) Descriptor() ([]byte, []int) {
	return file_pc1o_proto_rawDescGZIP(), []int{0}
}

func (x *PC1ORequest) GetSystemIdentifier() []byte {
	if x != nil {
		return x.SystemIdentifier
	}
	return nil
}

func (x *PC1ORequest) GetNodeIdentifier() string {
	if x != nil {
		return x.NodeIdentifier
	}
	return ""
}

func (x *PC1ORequest) GetUnicityCertificate() *certificates.UnicityCertificate {
	if x != nil {
		return x.UnicityCertificate
	}
	return nil
}

func (x *PC1ORequest) GetTransactions() []*transaction.Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *PC1ORequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

var File_pc1o_proto protoreflect.FileDescriptor

var file_pc1o_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x70, 0x63, 0x31, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xf9, 0x01, 0x0a, 0x0b, 0x50, 0x43, 0x31, 0x4f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x11, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e, 0x6f, 0x64, 0x65, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x44, 0x0a, 0x13, 0x75, 0x6e, 0x69,
	0x63, 0x69, 0x74, 0x79, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x55, 0x6e, 0x69, 0x63, 0x69, 0x74, 0x79,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x12, 0x75, 0x6e, 0x69,
	0x63, 0x69, 0x74, 0x79, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12,
	0x30, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42,
	0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x64, 0x63, 0x2e, 0x65, 0x65, 0x2e, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62,
	0x69, 0x6c, 0x6c, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f,
	0x70, 0x63, 0x31, 0x6f, 0x3b, 0x70, 0x63, 0x31, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_pc1o_proto_rawDescOnce sync.Once
	file_pc1o_proto_rawDescData = file_pc1o_proto_rawDesc
)

func file_pc1o_proto_rawDescGZIP() []byte {
	file_pc1o_proto_rawDescOnce.Do(func() {
		file_pc1o_proto_rawDescData = protoimpl.X.CompressGZIP(file_pc1o_proto_rawDescData)
	})
	return file_pc1o_proto_rawDescData
}

var file_pc1o_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pc1o_proto_goTypes = []interface{}{
	(*PC1ORequest)(nil),                     // 0: PC1ORequest
	(*certificates.UnicityCertificate)(nil), // 1: UnicityCertificate
	(*transaction.Transaction)(nil),         // 2: Transaction
}
var file_pc1o_proto_depIdxs = []int32{
	1, // 0: PC1ORequest.unicity_certificate:type_name -> UnicityCertificate
	2, // 1: PC1ORequest.transactions:type_name -> Transaction
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pc1o_proto_init() }
func file_pc1o_proto_init() {
	if File_pc1o_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pc1o_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PC1ORequest); i {
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
			RawDescriptor: file_pc1o_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pc1o_proto_goTypes,
		DependencyIndexes: file_pc1o_proto_depIdxs,
		MessageInfos:      file_pc1o_proto_msgTypes,
	}.Build()
	File_pc1o_proto = out.File
	file_pc1o_proto_rawDesc = nil
	file_pc1o_proto_goTypes = nil
	file_pc1o_proto_depIdxs = nil
}
