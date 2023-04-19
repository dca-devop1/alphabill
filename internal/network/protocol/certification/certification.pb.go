// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: certification.proto

package certification

import (
	certificates "github.com/alphabill-org/alphabill/internal/certificates"
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

type BlockCertificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SystemIdentifier []byte                    `protobuf:"bytes,1,opt,name=system_identifier,json=systemIdentifier,proto3" json:"system_identifier,omitempty"`
	NodeIdentifier   string                    `protobuf:"bytes,2,opt,name=node_identifier,json=nodeIdentifier,proto3" json:"node_identifier,omitempty"`
	InputRecord      *certificates.InputRecord `protobuf:"bytes,3,opt,name=input_record,json=inputRecord,proto3" json:"input_record,omitempty"`
	Signature        []byte                    `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *BlockCertificationRequest) Reset() {
	*x = BlockCertificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_certification_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockCertificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockCertificationRequest) ProtoMessage() {}

func (x *BlockCertificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_certification_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockCertificationRequest.ProtoReflect.Descriptor instead.
func (*BlockCertificationRequest) Descriptor() ([]byte, []int) {
	return file_certification_proto_rawDescGZIP(), []int{0}
}

func (x *BlockCertificationRequest) GetSystemIdentifier() []byte {
	if x != nil {
		return x.SystemIdentifier
	}
	return nil
}

func (x *BlockCertificationRequest) GetNodeIdentifier() string {
	if x != nil {
		return x.NodeIdentifier
	}
	return ""
}

func (x *BlockCertificationRequest) GetInputRecord() *certificates.InputRecord {
	if x != nil {
		return x.InputRecord
	}
	return nil
}

func (x *BlockCertificationRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

var File_certification_proto protoreflect.FileDescriptor

var file_certification_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc0, 0x01, 0x0a, 0x19, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x11, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x10, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e,
	0x6f, 0x64, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x2f, 0x0a,
	0x0c, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x52, 0x0b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x5a, 0x5a, 0x58,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69,
	0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x63, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3b, 0x63, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_certification_proto_rawDescOnce sync.Once
	file_certification_proto_rawDescData = file_certification_proto_rawDesc
)

func file_certification_proto_rawDescGZIP() []byte {
	file_certification_proto_rawDescOnce.Do(func() {
		file_certification_proto_rawDescData = protoimpl.X.CompressGZIP(file_certification_proto_rawDescData)
	})
	return file_certification_proto_rawDescData
}

var file_certification_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_certification_proto_goTypes = []interface{}{
	(*BlockCertificationRequest)(nil), // 0: BlockCertificationRequest
	(*certificates.InputRecord)(nil),  // 1: InputRecord
}
var file_certification_proto_depIdxs = []int32{
	1, // 0: BlockCertificationRequest.input_record:type_name -> InputRecord
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_certification_proto_init() }
func file_certification_proto_init() {
	if File_certification_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_certification_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockCertificationRequest); i {
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
			RawDescriptor: file_certification_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_certification_proto_goTypes,
		DependencyIndexes: file_certification_proto_depIdxs,
		MessageInfos:      file_certification_proto_msgTypes,
	}.Build()
	File_certification_proto = out.File
	file_certification_proto_rawDesc = nil
	file_certification_proto_goTypes = nil
	file_certification_proto_depIdxs = nil
}
