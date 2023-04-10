// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: alphabill.proto

package alphabill

import (
	block "github.com/alphabill-org/alphabill/internal/block"
	txsystem "github.com/alphabill-org/alphabill/internal/txsystem"
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

type GetBlocksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNumber uint64 `protobuf:"varint,1,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	BlockCount  uint64 `protobuf:"varint,2,opt,name=block_count,json=blockCount,proto3" json:"block_count,omitempty"`
}

func (x *GetBlocksRequest) Reset() {
	*x = GetBlocksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBlocksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBlocksRequest) ProtoMessage() {}

func (x *GetBlocksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBlocksRequest.ProtoReflect.Descriptor instead.
func (*GetBlocksRequest) Descriptor() ([]byte, []int) {
	return file_alphabill_proto_rawDescGZIP(), []int{0}
}

func (x *GetBlocksRequest) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *GetBlocksRequest) GetBlockCount() uint64 {
	if x != nil {
		return x.BlockCount
	}
	return 0
}

type GetBlocksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMessage        string         `protobuf:"bytes,1,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	MaxBlockNumber      uint64         `protobuf:"varint,2,opt,name=max_block_number,json=maxBlockNumber,proto3" json:"max_block_number,omitempty"`
	Blocks              []*block.Block `protobuf:"bytes,3,rep,name=blocks,proto3" json:"blocks,omitempty"`
	MaxRoundNumber      uint64         `protobuf:"varint,4,opt,name=max_round_number,json=maxRoundNumber,proto3" json:"max_round_number,omitempty"`
	BatchMaxBlockNumber uint64         `protobuf:"varint,5,opt,name=batch_max_block_number,json=batchMaxBlockNumber,proto3" json:"batch_max_block_number,omitempty"`
}

func (x *GetBlocksResponse) Reset() {
	*x = GetBlocksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBlocksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBlocksResponse) ProtoMessage() {}

func (x *GetBlocksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBlocksResponse.ProtoReflect.Descriptor instead.
func (*GetBlocksResponse) Descriptor() ([]byte, []int) {
	return file_alphabill_proto_rawDescGZIP(), []int{1}
}

func (x *GetBlocksResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *GetBlocksResponse) GetMaxBlockNumber() uint64 {
	if x != nil {
		return x.MaxBlockNumber
	}
	return 0
}

func (x *GetBlocksResponse) GetBlocks() []*block.Block {
	if x != nil {
		return x.Blocks
	}
	return nil
}

func (x *GetBlocksResponse) GetMaxRoundNumber() uint64 {
	if x != nil {
		return x.MaxRoundNumber
	}
	return 0
}

func (x *GetBlocksResponse) GetBatchMaxBlockNumber() uint64 {
	if x != nil {
		return x.BatchMaxBlockNumber
	}
	return 0
}

type GetBlockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNo uint64 `protobuf:"varint,1,opt,name=block_no,json=blockNo,proto3" json:"block_no,omitempty"`
}

func (x *GetBlockRequest) Reset() {
	*x = GetBlockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBlockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBlockRequest) ProtoMessage() {}

func (x *GetBlockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBlockRequest.ProtoReflect.Descriptor instead.
func (*GetBlockRequest) Descriptor() ([]byte, []int) {
	return file_alphabill_proto_rawDescGZIP(), []int{2}
}

func (x *GetBlockRequest) GetBlockNo() uint64 {
	if x != nil {
		return x.BlockNo
	}
	return 0
}

type GetBlockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMessage string       `protobuf:"bytes,1,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	Block        *block.Block `protobuf:"bytes,2,opt,name=block,proto3" json:"block,omitempty"`
}

func (x *GetBlockResponse) Reset() {
	*x = GetBlockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBlockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBlockResponse) ProtoMessage() {}

func (x *GetBlockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBlockResponse.ProtoReflect.Descriptor instead.
func (*GetBlockResponse) Descriptor() ([]byte, []int) {
	return file_alphabill_proto_rawDescGZIP(), []int{3}
}

func (x *GetBlockResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *GetBlockResponse) GetBlock() *block.Block {
	if x != nil {
		return x.Block
	}
	return nil
}

type GetMaxBlockNoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetMaxBlockNoRequest) Reset() {
	*x = GetMaxBlockNoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMaxBlockNoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMaxBlockNoRequest) ProtoMessage() {}

func (x *GetMaxBlockNoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMaxBlockNoRequest.ProtoReflect.Descriptor instead.
func (*GetMaxBlockNoRequest) Descriptor() ([]byte, []int) {
	return file_alphabill_proto_rawDescGZIP(), []int{4}
}

type GetMaxBlockNoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMessage   string `protobuf:"bytes,1,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	BlockNo        uint64 `protobuf:"varint,2,opt,name=block_no,json=blockNo,proto3" json:"block_no,omitempty"`
	MaxRoundNumber uint64 `protobuf:"varint,3,opt,name=max_round_number,json=maxRoundNumber,proto3" json:"max_round_number,omitempty"`
}

func (x *GetMaxBlockNoResponse) Reset() {
	*x = GetMaxBlockNoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMaxBlockNoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMaxBlockNoResponse) ProtoMessage() {}

func (x *GetMaxBlockNoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMaxBlockNoResponse.ProtoReflect.Descriptor instead.
func (*GetMaxBlockNoResponse) Descriptor() ([]byte, []int) {
	return file_alphabill_proto_rawDescGZIP(), []int{5}
}

func (x *GetMaxBlockNoResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *GetMaxBlockNoResponse) GetBlockNo() uint64 {
	if x != nil {
		return x.BlockNo
	}
	return 0
}

func (x *GetMaxBlockNoResponse) GetMaxRoundNumber() uint64 {
	if x != nil {
		return x.MaxRoundNumber
	}
	return 0
}

var File_alphabill_proto protoreflect.FileDescriptor

var file_alphabill_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x61, 0x62, 0x72, 0x70, 0x63, 0x1a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x1f, 0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0xe1, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x6d,
	0x61, 0x78, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x6d, 0x61, 0x78, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x06, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x28, 0x0a, 0x10, 0x6d, 0x61, 0x78, 0x5f, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0e, 0x6d, 0x61, 0x78, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x33, 0x0a, 0x16, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x13, 0x62, 0x61, 0x74, 0x63, 0x68, 0x4d, 0x61, 0x78, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x22, 0x2c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x5f, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x4e, 0x6f, 0x22, 0x55, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x05, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x52, 0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x22, 0x16, 0x0a, 0x14, 0x47, 0x65, 0x74,
	0x4d, 0x61, 0x78, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x81, 0x01, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x78, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x4e, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x19, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x6f, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x12, 0x28, 0x0a, 0x10, 0x6d,
	0x61, 0x78, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x6d, 0x61, 0x78, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x32, 0x9d, 0x02, 0x0a, 0x10, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x62,
	0x69, 0x6c, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x12, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x0c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x14,
	0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x12, 0x16, 0x2e, 0x61, 0x62, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x61, 0x62, 0x72,
	0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x73, 0x12, 0x17, 0x2e, 0x61, 0x62, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x62,
	0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4d, 0x61,
	0x78, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x12, 0x1b, 0x2e, 0x61, 0x62, 0x72, 0x70, 0x63,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x78, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x62, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65,
	0x74, 0x4d, 0x61, 0x78, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x6f, 0x72,
	0x67, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x3b, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69,
	0x6c, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_alphabill_proto_rawDescOnce sync.Once
	file_alphabill_proto_rawDescData = file_alphabill_proto_rawDesc
)

func file_alphabill_proto_rawDescGZIP() []byte {
	file_alphabill_proto_rawDescOnce.Do(func() {
		file_alphabill_proto_rawDescData = protoimpl.X.CompressGZIP(file_alphabill_proto_rawDescData)
	})
	return file_alphabill_proto_rawDescData
}

var file_alphabill_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_alphabill_proto_goTypes = []interface{}{
	(*GetBlocksRequest)(nil),             // 0: abrpc.GetBlocksRequest
	(*GetBlocksResponse)(nil),            // 1: abrpc.GetBlocksResponse
	(*GetBlockRequest)(nil),              // 2: abrpc.GetBlockRequest
	(*GetBlockResponse)(nil),             // 3: abrpc.GetBlockResponse
	(*GetMaxBlockNoRequest)(nil),         // 4: abrpc.GetMaxBlockNoRequest
	(*GetMaxBlockNoResponse)(nil),        // 5: abrpc.GetMaxBlockNoResponse
	(*block.Block)(nil),                  // 6: Block
	(*txsystem.Transaction)(nil),         // 7: Transaction
	(*txsystem.TransactionResponse)(nil), // 8: TransactionResponse
}
var file_alphabill_proto_depIdxs = []int32{
	6, // 0: abrpc.GetBlocksResponse.blocks:type_name -> Block
	6, // 1: abrpc.GetBlockResponse.block:type_name -> Block
	7, // 2: abrpc.AlphabillService.ProcessTransaction:input_type -> Transaction
	2, // 3: abrpc.AlphabillService.GetBlock:input_type -> abrpc.GetBlockRequest
	0, // 4: abrpc.AlphabillService.GetBlocks:input_type -> abrpc.GetBlocksRequest
	4, // 5: abrpc.AlphabillService.GetMaxBlockNo:input_type -> abrpc.GetMaxBlockNoRequest
	8, // 6: abrpc.AlphabillService.ProcessTransaction:output_type -> TransactionResponse
	3, // 7: abrpc.AlphabillService.GetBlock:output_type -> abrpc.GetBlockResponse
	1, // 8: abrpc.AlphabillService.GetBlocks:output_type -> abrpc.GetBlocksResponse
	5, // 9: abrpc.AlphabillService.GetMaxBlockNo:output_type -> abrpc.GetMaxBlockNoResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_alphabill_proto_init() }
func file_alphabill_proto_init() {
	if File_alphabill_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_alphabill_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBlocksRequest); i {
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
		file_alphabill_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBlocksResponse); i {
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
		file_alphabill_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBlockRequest); i {
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
		file_alphabill_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBlockResponse); i {
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
		file_alphabill_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMaxBlockNoRequest); i {
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
		file_alphabill_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMaxBlockNoResponse); i {
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
			RawDescriptor: file_alphabill_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_alphabill_proto_goTypes,
		DependencyIndexes: file_alphabill_proto_depIdxs,
		MessageInfos:      file_alphabill_proto_msgTypes,
	}.Build()
	File_alphabill_proto = out.File
	file_alphabill_proto_rawDesc = nil
	file_alphabill_proto_goTypes = nil
	file_alphabill_proto_depIdxs = nil
}
