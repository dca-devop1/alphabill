// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: genesis.proto

package genesis

import (
	certificates "gitdc.ee.guardtime.com/alphabill/alphabill/internal/certificates"
	p1 "gitdc.ee.guardtime.com/alphabill/alphabill/internal/protocol/p1"
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

type RootGenesis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Partitions []*GenesisPartitionRecord `protobuf:"bytes,1,rep,name=partitions,proto3" json:"partitions,omitempty"`
	// root chain public key
	TrustBase     []byte `protobuf:"bytes,2,opt,name=trust_base,json=trustBase,proto3" json:"trust_base,omitempty"`
	HashAlgorithm uint32 `protobuf:"varint,3,opt,name=hash_algorithm,json=hashAlgorithm,proto3" json:"hash_algorithm,omitempty"`
}

func (x *RootGenesis) Reset() {
	*x = RootGenesis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RootGenesis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RootGenesis) ProtoMessage() {}

func (x *RootGenesis) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RootGenesis.ProtoReflect.Descriptor instead.
func (*RootGenesis) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{0}
}

func (x *RootGenesis) GetPartitions() []*GenesisPartitionRecord {
	if x != nil {
		return x.Partitions
	}
	return nil
}

func (x *RootGenesis) GetTrustBase() []byte {
	if x != nil {
		return x.TrustBase
	}
	return nil
}

func (x *RootGenesis) GetHashAlgorithm() uint32 {
	if x != nil {
		return x.HashAlgorithm
	}
	return 0
}

type GenesisPartitionRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes                   []*PartitionNode                 `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Certificate             *certificates.UnicityCertificate `protobuf:"bytes,2,opt,name=certificate,proto3" json:"certificate,omitempty"`
	SystemDescriptionRecord *SystemDescriptionRecord         `protobuf:"bytes,3,opt,name=system_description_record,json=systemDescriptionRecord,proto3" json:"system_description_record,omitempty"`
}

func (x *GenesisPartitionRecord) Reset() {
	*x = GenesisPartitionRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenesisPartitionRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenesisPartitionRecord) ProtoMessage() {}

func (x *GenesisPartitionRecord) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenesisPartitionRecord.ProtoReflect.Descriptor instead.
func (*GenesisPartitionRecord) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{1}
}

func (x *GenesisPartitionRecord) GetNodes() []*PartitionNode {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *GenesisPartitionRecord) GetCertificate() *certificates.UnicityCertificate {
	if x != nil {
		return x.Certificate
	}
	return nil
}

func (x *GenesisPartitionRecord) GetSystemDescriptionRecord() *SystemDescriptionRecord {
	if x != nil {
		return x.SystemDescriptionRecord
	}
	return nil
}

type PartitionGenesis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SystemDescriptionRecord *SystemDescriptionRecord         `protobuf:"bytes,1,opt,name=system_description_record,json=systemDescriptionRecord,proto3" json:"system_description_record,omitempty"`
	Certificate             *certificates.UnicityCertificate `protobuf:"bytes,2,opt,name=certificate,proto3" json:"certificate,omitempty"`
	TrustBase               []byte                           `protobuf:"bytes,3,opt,name=trust_base,json=trustBase,proto3" json:"trust_base,omitempty"`
	EncryptionKey           []byte                           `protobuf:"bytes,4,opt,name=encryption_key,json=encryptionKey,proto3" json:"encryption_key,omitempty"`
	Keys                    []*PublicKeyInfo                 `protobuf:"bytes,5,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *PartitionGenesis) Reset() {
	*x = PartitionGenesis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartitionGenesis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartitionGenesis) ProtoMessage() {}

func (x *PartitionGenesis) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartitionGenesis.ProtoReflect.Descriptor instead.
func (*PartitionGenesis) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{2}
}

func (x *PartitionGenesis) GetSystemDescriptionRecord() *SystemDescriptionRecord {
	if x != nil {
		return x.SystemDescriptionRecord
	}
	return nil
}

func (x *PartitionGenesis) GetCertificate() *certificates.UnicityCertificate {
	if x != nil {
		return x.Certificate
	}
	return nil
}

func (x *PartitionGenesis) GetTrustBase() []byte {
	if x != nil {
		return x.TrustBase
	}
	return nil
}

func (x *PartitionGenesis) GetEncryptionKey() []byte {
	if x != nil {
		return x.EncryptionKey
	}
	return nil
}

func (x *PartitionGenesis) GetKeys() []*PublicKeyInfo {
	if x != nil {
		return x.Keys
	}
	return nil
}

type PublicKeyInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeIdentifier      string `protobuf:"bytes,1,opt,name=node_identifier,json=nodeIdentifier,proto3" json:"node_identifier,omitempty"`
	SigningPublicKey    []byte `protobuf:"bytes,2,opt,name=signing_public_key,json=signingPublicKey,proto3" json:"signing_public_key,omitempty"`
	EncryptionPublicKey []byte `protobuf:"bytes,3,opt,name=encryption_public_key,json=encryptionPublicKey,proto3" json:"encryption_public_key,omitempty"`
}

func (x *PublicKeyInfo) Reset() {
	*x = PublicKeyInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicKeyInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicKeyInfo) ProtoMessage() {}

func (x *PublicKeyInfo) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicKeyInfo.ProtoReflect.Descriptor instead.
func (*PublicKeyInfo) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{3}
}

func (x *PublicKeyInfo) GetNodeIdentifier() string {
	if x != nil {
		return x.NodeIdentifier
	}
	return ""
}

func (x *PublicKeyInfo) GetSigningPublicKey() []byte {
	if x != nil {
		return x.SigningPublicKey
	}
	return nil
}

func (x *PublicKeyInfo) GetEncryptionPublicKey() []byte {
	if x != nil {
		return x.EncryptionPublicKey
	}
	return nil
}

type SystemDescriptionRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SystemIdentifier []byte `protobuf:"bytes,1,opt,name=system_identifier,json=systemIdentifier,proto3" json:"system_identifier,omitempty"`
	T2Timeout        uint32 `protobuf:"varint,2,opt,name=t2timeout,proto3" json:"t2timeout,omitempty"`
}

func (x *SystemDescriptionRecord) Reset() {
	*x = SystemDescriptionRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SystemDescriptionRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SystemDescriptionRecord) ProtoMessage() {}

func (x *SystemDescriptionRecord) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SystemDescriptionRecord.ProtoReflect.Descriptor instead.
func (*SystemDescriptionRecord) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{4}
}

func (x *SystemDescriptionRecord) GetSystemIdentifier() []byte {
	if x != nil {
		return x.SystemIdentifier
	}
	return nil
}

func (x *SystemDescriptionRecord) GetT2Timeout() uint32 {
	if x != nil {
		return x.T2Timeout
	}
	return 0
}

type PartitionRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SystemDescriptionRecord *SystemDescriptionRecord `protobuf:"bytes,1,opt,name=system_description_record,json=systemDescriptionRecord,proto3" json:"system_description_record,omitempty"`
	Validators              []*PartitionNode         `protobuf:"bytes,2,rep,name=validators,proto3" json:"validators,omitempty"`
}

func (x *PartitionRecord) Reset() {
	*x = PartitionRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartitionRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartitionRecord) ProtoMessage() {}

func (x *PartitionRecord) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartitionRecord.ProtoReflect.Descriptor instead.
func (*PartitionRecord) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{5}
}

func (x *PartitionRecord) GetSystemDescriptionRecord() *SystemDescriptionRecord {
	if x != nil {
		return x.SystemDescriptionRecord
	}
	return nil
}

func (x *PartitionRecord) GetValidators() []*PartitionNode {
	if x != nil {
		return x.Validators
	}
	return nil
}

type PartitionNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeIdentifier      string        `protobuf:"bytes,1,opt,name=node_identifier,json=nodeIdentifier,proto3" json:"node_identifier,omitempty"`
	SigningPublicKey    []byte        `protobuf:"bytes,2,opt,name=signing_public_key,json=signingPublicKey,proto3" json:"signing_public_key,omitempty"`
	EncryptionPublicKey []byte        `protobuf:"bytes,3,opt,name=encryption_public_key,json=encryptionPublicKey,proto3" json:"encryption_public_key,omitempty"`
	P1Request           *p1.P1Request `protobuf:"bytes,4,opt,name=p1_request,json=p1Request,proto3" json:"p1_request,omitempty"`
}

func (x *PartitionNode) Reset() {
	*x = PartitionNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartitionNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartitionNode) ProtoMessage() {}

func (x *PartitionNode) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartitionNode.ProtoReflect.Descriptor instead.
func (*PartitionNode) Descriptor() ([]byte, []int) {
	return file_genesis_proto_rawDescGZIP(), []int{6}
}

func (x *PartitionNode) GetNodeIdentifier() string {
	if x != nil {
		return x.NodeIdentifier
	}
	return ""
}

func (x *PartitionNode) GetSigningPublicKey() []byte {
	if x != nil {
		return x.SigningPublicKey
	}
	return nil
}

func (x *PartitionNode) GetEncryptionPublicKey() []byte {
	if x != nil {
		return x.EncryptionPublicKey
	}
	return nil
}

func (x *PartitionNode) GetP1Request() *p1.P1Request {
	if x != nil {
		return x.P1Request
	}
	return nil
}

var File_genesis_proto protoreflect.FileDescriptor

var file_genesis_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x08, 0x70, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8c, 0x01,
	0x0a, 0x0b, 0x52, 0x6f, 0x6f, 0x74, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x12, 0x37, 0x0a,
	0x0a, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x50, 0x61, 0x72, 0x74, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x74,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x72, 0x75, 0x73, 0x74, 0x5f,
	0x62, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x74, 0x72, 0x75, 0x73,
	0x74, 0x42, 0x61, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x61, 0x6c,
	0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x68,
	0x61, 0x73, 0x68, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x22, 0xcb, 0x01, 0x0a,
	0x16, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x24, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x35, 0x0a,
	0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x55, 0x6e, 0x69, 0x63, 0x69, 0x74, 0x79, 0x43, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x12, 0x54, 0x0a, 0x19, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x52, 0x17, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x22, 0x89, 0x02, 0x0a, 0x10, 0x50,
	0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x12,
	0x54, 0x0a, 0x19, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x17, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x35, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x55, 0x6e, 0x69,
	0x63, 0x69, 0x74, 0x79, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52,
	0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x74, 0x72, 0x75, 0x73, 0x74, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x09, 0x74, 0x72, 0x75, 0x73, 0x74, 0x42, 0x61, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4b,
	0x65, 0x79, 0x12, 0x22, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0x9a, 0x01, 0x0a, 0x0d, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65,
	0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x73,
	0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12,
	0x32, 0x0a, 0x15, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x13,
	0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x4b, 0x65, 0x79, 0x22, 0x64, 0x0a, 0x17, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x2b,
	0x0a, 0x11, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x32, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09,
	0x74, 0x32, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x97, 0x01, 0x0a, 0x0f, 0x50, 0x61,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x54, 0x0a,
	0x19, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x17, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x2e, 0x0a, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x73, 0x22, 0xc5, 0x01, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x2c,
	0x0a, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x73, 0x69, 0x67, 0x6e,
	0x69, 0x6e, 0x67, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x15,
	0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x13, 0x65, 0x6e, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79,
	0x12, 0x29, 0x0a, 0x0a, 0x70, 0x31, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x50, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x52, 0x09, 0x70, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x4f, 0x5a, 0x4d, 0x67,
	0x69, 0x74, 0x64, 0x63, 0x2e, 0x65, 0x65, 0x2e, 0x67, 0x75, 0x61, 0x72, 0x64, 0x74, 0x69, 0x6d,
	0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x65,
	0x73, 0x69, 0x73, 0x2f, 0x3b, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_genesis_proto_rawDescOnce sync.Once
	file_genesis_proto_rawDescData = file_genesis_proto_rawDesc
)

func file_genesis_proto_rawDescGZIP() []byte {
	file_genesis_proto_rawDescOnce.Do(func() {
		file_genesis_proto_rawDescData = protoimpl.X.CompressGZIP(file_genesis_proto_rawDescData)
	})
	return file_genesis_proto_rawDescData
}

var file_genesis_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_genesis_proto_goTypes = []interface{}{
	(*RootGenesis)(nil),                     // 0: RootGenesis
	(*GenesisPartitionRecord)(nil),          // 1: GenesisPartitionRecord
	(*PartitionGenesis)(nil),                // 2: PartitionGenesis
	(*PublicKeyInfo)(nil),                   // 3: PublicKeyInfo
	(*SystemDescriptionRecord)(nil),         // 4: SystemDescriptionRecord
	(*PartitionRecord)(nil),                 // 5: PartitionRecord
	(*PartitionNode)(nil),                   // 6: PartitionNode
	(*certificates.UnicityCertificate)(nil), // 7: UnicityCertificate
	(*p1.P1Request)(nil),                    // 8: P1Request
}
var file_genesis_proto_depIdxs = []int32{
	1,  // 0: RootGenesis.partitions:type_name -> GenesisPartitionRecord
	6,  // 1: GenesisPartitionRecord.nodes:type_name -> PartitionNode
	7,  // 2: GenesisPartitionRecord.certificate:type_name -> UnicityCertificate
	4,  // 3: GenesisPartitionRecord.system_description_record:type_name -> SystemDescriptionRecord
	4,  // 4: PartitionGenesis.system_description_record:type_name -> SystemDescriptionRecord
	7,  // 5: PartitionGenesis.certificate:type_name -> UnicityCertificate
	3,  // 6: PartitionGenesis.keys:type_name -> PublicKeyInfo
	4,  // 7: PartitionRecord.system_description_record:type_name -> SystemDescriptionRecord
	6,  // 8: PartitionRecord.validators:type_name -> PartitionNode
	8,  // 9: PartitionNode.p1_request:type_name -> P1Request
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_genesis_proto_init() }
func file_genesis_proto_init() {
	if File_genesis_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_genesis_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RootGenesis); i {
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
		file_genesis_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenesisPartitionRecord); i {
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
		file_genesis_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartitionGenesis); i {
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
		file_genesis_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublicKeyInfo); i {
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
		file_genesis_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SystemDescriptionRecord); i {
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
		file_genesis_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartitionRecord); i {
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
		file_genesis_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartitionNode); i {
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
			RawDescriptor: file_genesis_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_genesis_proto_goTypes,
		DependencyIndexes: file_genesis_proto_depIdxs,
		MessageInfos:      file_genesis_proto_msgTypes,
	}.Build()
	File_genesis_proto = out.File
	file_genesis_proto_rawDesc = nil
	file_genesis_proto_goTypes = nil
	file_genesis_proto_depIdxs = nil
}
