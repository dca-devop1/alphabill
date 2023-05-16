// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: block_proof.proto

package block

import (
	certificates "github.com/alphabill-org/alphabill/internal/certificates"
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

type ProofType int32

const (
	ProofType_PRIM       ProofType = 0
	ProofType_SEC        ProofType = 1
	ProofType_ONLYSEC    ProofType = 2
	ProofType_NOTRANS    ProofType = 3
	ProofType_EMPTYBLOCK ProofType = 4
)

// Enum value maps for ProofType.
var (
	ProofType_name = map[int32]string{
		0: "PRIM",
		1: "SEC",
		2: "ONLYSEC",
		3: "NOTRANS",
		4: "EMPTYBLOCK",
	}
	ProofType_value = map[string]int32{
		"PRIM":       0,
		"SEC":        1,
		"ONLYSEC":    2,
		"NOTRANS":    3,
		"EMPTYBLOCK": 4,
	}
)

func (x ProofType) Enum() *ProofType {
	p := new(ProofType)
	*p = x
	return p
}

func (x ProofType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProofType) Descriptor() protoreflect.EnumDescriptor {
	return file_block_proof_proto_enumTypes[0].Descriptor()
}

func (ProofType) Type() protoreflect.EnumType {
	return &file_block_proof_proto_enumTypes[0]
}

func (x ProofType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProofType.Descriptor instead.
func (ProofType) EnumDescriptor() ([]byte, []int) {
	return file_block_proof_proto_rawDescGZIP(), []int{0}
}

// wrapper around block proof with transaction data
type TxProof struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// deprecated: Use proof.UnicityCertificate.GetRoundNumber() instead; will be removed
	BlockNumber uint64                `protobuf:"varint,1,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	Tx          *txsystem.Transaction `protobuf:"bytes,2,opt,name=tx,proto3" json:"tx,omitempty"`
	Proof       *BlockProof           `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
}

func (x *TxProof) Reset() {
	*x = TxProof{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proof_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxProof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxProof) ProtoMessage() {}

func (x *TxProof) ProtoReflect() protoreflect.Message {
	mi := &file_block_proof_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxProof.ProtoReflect.Descriptor instead.
func (*TxProof) Descriptor() ([]byte, []int) {
	return file_block_proof_proto_rawDescGZIP(), []int{0}
}

func (x *TxProof) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *TxProof) GetTx() *txsystem.Transaction {
	if x != nil {
		return x.Tx
	}
	return nil
}

func (x *TxProof) GetProof() *BlockProof {
	if x != nil {
		return x.Proof
	}
	return nil
}

type BlockProof struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProofType        ProofType `protobuf:"varint,1,opt,name=proof_type,json=proofType,proto3,enum=ProofType" json:"proof_type,omitempty"`
	BlockHeaderHash  []byte    `protobuf:"bytes,2,opt,name=block_header_hash,json=blockHeaderHash,proto3" json:"block_header_hash,omitempty"`
	TransactionsHash []byte    `protobuf:"bytes,3,opt,name=transactions_hash,json=transactionsHash,proto3" json:"transactions_hash,omitempty"`
	// hash value of either primary tx or secondary txs or zero hash, depending on proof type
	HashValue          []byte                           `protobuf:"bytes,4,opt,name=hash_value,json=hashValue,proto3" json:"hash_value,omitempty"`
	BlockTreeHashChain *BlockTreeHashChain              `protobuf:"bytes,5,opt,name=block_tree_hash_chain,json=blockTreeHashChain,proto3" json:"block_tree_hash_chain,omitempty"`
	SecTreeHashChain   *SecTreeHashChain                `protobuf:"bytes,6,opt,name=sec_tree_hash_chain,json=secTreeHashChain,proto3" json:"sec_tree_hash_chain,omitempty"`
	UnicityCertificate *certificates.UnicityCertificate `protobuf:"bytes,7,opt,name=unicity_certificate,json=unicityCertificate,proto3" json:"unicity_certificate,omitempty"`
}

func (x *BlockProof) Reset() {
	*x = BlockProof{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proof_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockProof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockProof) ProtoMessage() {}

func (x *BlockProof) ProtoReflect() protoreflect.Message {
	mi := &file_block_proof_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockProof.ProtoReflect.Descriptor instead.
func (*BlockProof) Descriptor() ([]byte, []int) {
	return file_block_proof_proto_rawDescGZIP(), []int{1}
}

func (x *BlockProof) GetProofType() ProofType {
	if x != nil {
		return x.ProofType
	}
	return ProofType_PRIM
}

func (x *BlockProof) GetBlockHeaderHash() []byte {
	if x != nil {
		return x.BlockHeaderHash
	}
	return nil
}

func (x *BlockProof) GetTransactionsHash() []byte {
	if x != nil {
		return x.TransactionsHash
	}
	return nil
}

func (x *BlockProof) GetHashValue() []byte {
	if x != nil {
		return x.HashValue
	}
	return nil
}

func (x *BlockProof) GetBlockTreeHashChain() *BlockTreeHashChain {
	if x != nil {
		return x.BlockTreeHashChain
	}
	return nil
}

func (x *BlockProof) GetSecTreeHashChain() *SecTreeHashChain {
	if x != nil {
		return x.SecTreeHashChain
	}
	return nil
}

func (x *BlockProof) GetUnicityCertificate() *certificates.UnicityCertificate {
	if x != nil {
		return x.UnicityCertificate
	}
	return nil
}

type BlockTreeHashChain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*ChainItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *BlockTreeHashChain) Reset() {
	*x = BlockTreeHashChain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proof_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockTreeHashChain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockTreeHashChain) ProtoMessage() {}

func (x *BlockTreeHashChain) ProtoReflect() protoreflect.Message {
	mi := &file_block_proof_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockTreeHashChain.ProtoReflect.Descriptor instead.
func (*BlockTreeHashChain) Descriptor() ([]byte, []int) {
	return file_block_proof_proto_rawDescGZIP(), []int{2}
}

func (x *BlockTreeHashChain) GetItems() []*ChainItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type SecTreeHashChain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*MerklePathItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SecTreeHashChain) Reset() {
	*x = SecTreeHashChain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proof_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecTreeHashChain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecTreeHashChain) ProtoMessage() {}

func (x *SecTreeHashChain) ProtoReflect() protoreflect.Message {
	mi := &file_block_proof_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecTreeHashChain.ProtoReflect.Descriptor instead.
func (*SecTreeHashChain) Descriptor() ([]byte, []int) {
	return file_block_proof_proto_rawDescGZIP(), []int{3}
}

func (x *SecTreeHashChain) GetItems() []*MerklePathItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type ChainItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val  []byte `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
	Hash []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *ChainItem) Reset() {
	*x = ChainItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proof_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChainItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChainItem) ProtoMessage() {}

func (x *ChainItem) ProtoReflect() protoreflect.Message {
	mi := &file_block_proof_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChainItem.ProtoReflect.Descriptor instead.
func (*ChainItem) Descriptor() ([]byte, []int) {
	return file_block_proof_proto_rawDescGZIP(), []int{4}
}

func (x *ChainItem) GetVal() []byte {
	if x != nil {
		return x.Val
	}
	return nil
}

func (x *ChainItem) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type MerklePathItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DirectionLeft direction from parent node; left=true right=false
	DirectionLeft bool `protobuf:"varint,1,opt,name=direction_left,json=directionLeft,proto3" json:"direction_left,omitempty"`
	// PathItem Hash of Merkle Tree node
	PathItem []byte `protobuf:"bytes,2,opt,name=path_item,json=pathItem,proto3" json:"path_item,omitempty"`
}

func (x *MerklePathItem) Reset() {
	*x = MerklePathItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proof_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MerklePathItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MerklePathItem) ProtoMessage() {}

func (x *MerklePathItem) ProtoReflect() protoreflect.Message {
	mi := &file_block_proof_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MerklePathItem.ProtoReflect.Descriptor instead.
func (*MerklePathItem) Descriptor() ([]byte, []int) {
	return file_block_proof_proto_rawDescGZIP(), []int{5}
}

func (x *MerklePathItem) GetDirectionLeft() bool {
	if x != nil {
		return x.DirectionLeft
	}
	return false
}

func (x *MerklePathItem) GetPathItem() []byte {
	if x != nil {
		return x.PathItem
	}
	return nil
}

var File_block_proof_proto protoreflect.FileDescriptor

var file_block_proof_proto_rawDesc = []byte{
	0x0a, 0x11, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x07, 0x54, 0x78,
	0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x02, 0x74, 0x78, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x02, 0x74, 0x78, 0x12, 0x21, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f,
	0x6f, 0x66, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x22, 0xff, 0x02, 0x0a, 0x0a, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x29, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6f,
	0x66, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x50,
	0x72, 0x6f, 0x6f, 0x66, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x2b, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a,
	0x68, 0x61, 0x73, 0x68, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x09, 0x68, 0x61, 0x73, 0x68, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x46, 0x0a, 0x15, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x72, 0x65, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x63,
	0x68, 0x61, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x54, 0x72, 0x65, 0x65, 0x48, 0x61, 0x73, 0x68, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x52,
	0x12, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x72, 0x65, 0x65, 0x48, 0x61, 0x73, 0x68, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x12, 0x40, 0x0a, 0x13, 0x73, 0x65, 0x63, 0x5f, 0x74, 0x72, 0x65, 0x65, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x5f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x53, 0x65, 0x63, 0x54, 0x72, 0x65, 0x65, 0x48, 0x61, 0x73, 0x68, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x52, 0x10, 0x73, 0x65, 0x63, 0x54, 0x72, 0x65, 0x65, 0x48, 0x61, 0x73, 0x68,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x44, 0x0a, 0x13, 0x75, 0x6e, 0x69, 0x63, 0x69, 0x74, 0x79,
	0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x55, 0x6e, 0x69, 0x63, 0x69, 0x74, 0x79, 0x43, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x12, 0x75, 0x6e, 0x69, 0x63, 0x69, 0x74, 0x79,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x22, 0x36, 0x0a, 0x12, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x72, 0x65, 0x65, 0x48, 0x61, 0x73, 0x68, 0x43, 0x68, 0x61, 0x69,
	0x6e, 0x12, 0x20, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x22, 0x39, 0x0a, 0x10, 0x53, 0x65, 0x63, 0x54, 0x72, 0x65, 0x65, 0x48, 0x61,
	0x73, 0x68, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x50,
	0x61, 0x74, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x31,
	0x0a, 0x09, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x76,
	0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x22, 0x54, 0x0a, 0x0e, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x25, 0x0a, 0x0e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x6c, 0x65, 0x66, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x65, 0x66, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61,
	0x74, 0x68, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x70,
	0x61, 0x74, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x2a, 0x48, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x6f, 0x66,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x52, 0x49, 0x4d, 0x10, 0x00, 0x12, 0x07,
	0x0a, 0x03, 0x53, 0x45, 0x43, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4f, 0x4e, 0x4c, 0x59, 0x53,
	0x45, 0x43, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x4f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x10,
	0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x10,
	0x04, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x3b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_block_proof_proto_rawDescOnce sync.Once
	file_block_proof_proto_rawDescData = file_block_proof_proto_rawDesc
)

func file_block_proof_proto_rawDescGZIP() []byte {
	file_block_proof_proto_rawDescOnce.Do(func() {
		file_block_proof_proto_rawDescData = protoimpl.X.CompressGZIP(file_block_proof_proto_rawDescData)
	})
	return file_block_proof_proto_rawDescData
}

var file_block_proof_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_block_proof_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_block_proof_proto_goTypes = []interface{}{
	(ProofType)(0),                          // 0: ProofType
	(*TxProof)(nil),                         // 1: TxProof
	(*BlockProof)(nil),                      // 2: BlockProof
	(*BlockTreeHashChain)(nil),              // 3: BlockTreeHashChain
	(*SecTreeHashChain)(nil),                // 4: SecTreeHashChain
	(*ChainItem)(nil),                       // 5: ChainItem
	(*MerklePathItem)(nil),                  // 6: MerklePathItem
	(*txsystem.Transaction)(nil),            // 7: Transaction
	(*certificates.UnicityCertificate)(nil), // 8: UnicityCertificate
}
var file_block_proof_proto_depIdxs = []int32{
	7, // 0: TxProof.tx:type_name -> Transaction
	2, // 1: TxProof.proof:type_name -> BlockProof
	0, // 2: BlockProof.proof_type:type_name -> ProofType
	3, // 3: BlockProof.block_tree_hash_chain:type_name -> BlockTreeHashChain
	4, // 4: BlockProof.sec_tree_hash_chain:type_name -> SecTreeHashChain
	8, // 5: BlockProof.unicity_certificate:type_name -> UnicityCertificate
	5, // 6: BlockTreeHashChain.items:type_name -> ChainItem
	6, // 7: SecTreeHashChain.items:type_name -> MerklePathItem
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_block_proof_proto_init() }
func file_block_proof_proto_init() {
	if File_block_proof_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_block_proof_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxProof); i {
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
		file_block_proof_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockProof); i {
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
		file_block_proof_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockTreeHashChain); i {
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
		file_block_proof_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecTreeHashChain); i {
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
		file_block_proof_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChainItem); i {
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
		file_block_proof_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MerklePathItem); i {
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
			RawDescriptor: file_block_proof_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_block_proof_proto_goTypes,
		DependencyIndexes: file_block_proof_proto_depIdxs,
		EnumInfos:         file_block_proof_proto_enumTypes,
		MessageInfos:      file_block_proof_proto_msgTypes,
	}.Build()
	File_block_proof_proto = out.File
	file_block_proof_proto_rawDesc = nil
	file_block_proof_proto_goTypes = nil
	file_block_proof_proto_depIdxs = nil
}
