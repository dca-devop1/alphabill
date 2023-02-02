// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: fee_credit_txs.proto

package fc

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

type TransferFeeCreditOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// amount to transfer
	Amount uint64 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	// system_identifier of the target partition
	TargetSystemIdentifier []byte `protobuf:"bytes,2,opt,name=target_system_identifier,json=targetSystemIdentifier,proto3" json:"target_system_identifier,omitempty"`
	// unit id of the corresponding “add fee credit” transaction
	TargetRecordId []byte `protobuf:"bytes,3,opt,name=target_record_id,json=targetRecordId,proto3" json:"target_record_id,omitempty"`
	// earliest round when the corresponding “add fee credit” transaction can be executed in the target system
	EarliestAdditionTime uint64 `protobuf:"varint,4,opt,name=earliest_addition_time,json=earliestAdditionTime,proto3" json:"earliest_addition_time,omitempty"`
	// latest round when the corresponding “add fee credit” transaction can be executed in the target system
	LatestAdditionTime uint64 `protobuf:"varint,5,opt,name=latest_addition_time,json=latestAdditionTime,proto3" json:"latest_addition_time,omitempty"`
	// the current state hash of the target credit record if the record exists, or to nil if the record does not exist yet
	Nonce []byte `protobuf:"bytes,6,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// hash of this unit's previous transacton
	Backlink []byte `protobuf:"bytes,7,opt,name=backlink,proto3" json:"backlink,omitempty"`
}

func (x *TransferFeeCreditOrder) Reset() {
	*x = TransferFeeCreditOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fee_credit_txs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferFeeCreditOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferFeeCreditOrder) ProtoMessage() {}

func (x *TransferFeeCreditOrder) ProtoReflect() protoreflect.Message {
	mi := &file_fee_credit_txs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferFeeCreditOrder.ProtoReflect.Descriptor instead.
func (*TransferFeeCreditOrder) Descriptor() ([]byte, []int) {
	return file_fee_credit_txs_proto_rawDescGZIP(), []int{0}
}

func (x *TransferFeeCreditOrder) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *TransferFeeCreditOrder) GetTargetSystemIdentifier() []byte {
	if x != nil {
		return x.TargetSystemIdentifier
	}
	return nil
}

func (x *TransferFeeCreditOrder) GetTargetRecordId() []byte {
	if x != nil {
		return x.TargetRecordId
	}
	return nil
}

func (x *TransferFeeCreditOrder) GetEarliestAdditionTime() uint64 {
	if x != nil {
		return x.EarliestAdditionTime
	}
	return 0
}

func (x *TransferFeeCreditOrder) GetLatestAdditionTime() uint64 {
	if x != nil {
		return x.LatestAdditionTime
	}
	return 0
}

func (x *TransferFeeCreditOrder) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *TransferFeeCreditOrder) GetBacklink() []byte {
	if x != nil {
		return x.Backlink
	}
	return nil
}

type AddFeeCreditOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// target fee credit record owner condition
	FeeCreditOwnerCondition []byte `protobuf:"bytes,1,opt,name=fee_credit_owner_condition,json=feeCreditOwnerCondition,proto3" json:"fee_credit_owner_condition,omitempty"`
	// bill transfer record of type "transfer fee credit"
	FeeCreditTransfer *txsystem.Transaction `protobuf:"bytes,2,opt,name=fee_credit_transfer,json=feeCreditTransfer,proto3" json:"fee_credit_transfer,omitempty"`
	// block proof of "transfer fee credit" transaction
	FeeCreditTransferProof *block.BlockProof `protobuf:"bytes,3,opt,name=fee_credit_transfer_proof,json=feeCreditTransferProof,proto3" json:"fee_credit_transfer_proof,omitempty"`
}

func (x *AddFeeCreditOrder) Reset() {
	*x = AddFeeCreditOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fee_credit_txs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddFeeCreditOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddFeeCreditOrder) ProtoMessage() {}

func (x *AddFeeCreditOrder) ProtoReflect() protoreflect.Message {
	mi := &file_fee_credit_txs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddFeeCreditOrder.ProtoReflect.Descriptor instead.
func (*AddFeeCreditOrder) Descriptor() ([]byte, []int) {
	return file_fee_credit_txs_proto_rawDescGZIP(), []int{1}
}

func (x *AddFeeCreditOrder) GetFeeCreditOwnerCondition() []byte {
	if x != nil {
		return x.FeeCreditOwnerCondition
	}
	return nil
}

func (x *AddFeeCreditOrder) GetFeeCreditTransfer() *txsystem.Transaction {
	if x != nil {
		return x.FeeCreditTransfer
	}
	return nil
}

func (x *AddFeeCreditOrder) GetFeeCreditTransferProof() *block.BlockProof {
	if x != nil {
		return x.FeeCreditTransferProof
	}
	return nil
}

type CloseFeeCreditOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// current balance of the fee credit record
	Amount uint64 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	// unit id of the fee credit record in money partition
	TargetUnitId []byte `protobuf:"bytes,2,opt,name=target_unit_id,json=targetUnitId,proto3" json:"target_unit_id,omitempty"`
	// the current state hash of the target unit in money partition
	Nonce []byte `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (x *CloseFeeCreditOrder) Reset() {
	*x = CloseFeeCreditOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fee_credit_txs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseFeeCreditOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseFeeCreditOrder) ProtoMessage() {}

func (x *CloseFeeCreditOrder) ProtoReflect() protoreflect.Message {
	mi := &file_fee_credit_txs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseFeeCreditOrder.ProtoReflect.Descriptor instead.
func (*CloseFeeCreditOrder) Descriptor() ([]byte, []int) {
	return file_fee_credit_txs_proto_rawDescGZIP(), []int{2}
}

func (x *CloseFeeCreditOrder) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *CloseFeeCreditOrder) GetTargetUnitId() []byte {
	if x != nil {
		return x.TargetUnitId
	}
	return nil
}

func (x *CloseFeeCreditOrder) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

type ReclaimFeeCreditOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// bill transfer record of type "close fee credit"
	CloseFeeCreditTransfer *txsystem.Transaction `protobuf:"bytes,1,opt,name=close_fee_credit_transfer,json=closeFeeCreditTransfer,proto3" json:"close_fee_credit_transfer,omitempty"`
	// block proof of "close fee credit" transaction
	CloseFeeCreditProof *block.BlockProof `protobuf:"bytes,2,opt,name=close_fee_credit_proof,json=closeFeeCreditProof,proto3" json:"close_fee_credit_proof,omitempty"`
	// hash of this unit's previous transacton
	Backlink []byte `protobuf:"bytes,3,opt,name=backlink,proto3" json:"backlink,omitempty"`
}

func (x *ReclaimFeeCreditOrder) Reset() {
	*x = ReclaimFeeCreditOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fee_credit_txs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReclaimFeeCreditOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReclaimFeeCreditOrder) ProtoMessage() {}

func (x *ReclaimFeeCreditOrder) ProtoReflect() protoreflect.Message {
	mi := &file_fee_credit_txs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReclaimFeeCreditOrder.ProtoReflect.Descriptor instead.
func (*ReclaimFeeCreditOrder) Descriptor() ([]byte, []int) {
	return file_fee_credit_txs_proto_rawDescGZIP(), []int{3}
}

func (x *ReclaimFeeCreditOrder) GetCloseFeeCreditTransfer() *txsystem.Transaction {
	if x != nil {
		return x.CloseFeeCreditTransfer
	}
	return nil
}

func (x *ReclaimFeeCreditOrder) GetCloseFeeCreditProof() *block.BlockProof {
	if x != nil {
		return x.CloseFeeCreditProof
	}
	return nil
}

func (x *ReclaimFeeCreditOrder) GetBacklink() []byte {
	if x != nil {
		return x.Backlink
	}
	return nil
}

var File_fee_credit_txs_proto protoreflect.FileDescriptor

var file_fee_credit_txs_proto_rawDesc = []byte{
	0x0a, 0x14, 0x66, 0x65, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x74, 0x78, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xae, 0x02, 0x0a,
	0x16, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x46, 0x65, 0x65, 0x43, 0x72, 0x65, 0x64,
	0x69, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x38, 0x0a, 0x18, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x16, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x16, 0x65, 0x61, 0x72, 0x6c, 0x69, 0x65, 0x73, 0x74, 0x5f,
	0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x14, 0x65, 0x61, 0x72, 0x6c, 0x69, 0x65, 0x73, 0x74, 0x41, 0x64, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x6c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x12, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x41,
	0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6e,
	0x6f, 0x6e, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x22, 0xd6, 0x01,
	0x0a, 0x11, 0x41, 0x64, 0x64, 0x46, 0x65, 0x65, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x12, 0x3b, 0x0a, 0x1a, 0x66, 0x65, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x69,
	0x74, 0x5f, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x17, 0x66, 0x65, 0x65, 0x43, 0x72, 0x65, 0x64,
	0x69, 0x74, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x3c, 0x0a, 0x13, 0x66, 0x65, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x66, 0x65, 0x65,
	0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x12, 0x46,
	0x0a, 0x19, 0x66, 0x65, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x66, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0b, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x52, 0x16,
	0x66, 0x65, 0x65, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x22, 0x69, 0x0a, 0x13, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x46,
	0x65, 0x65, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f,
	0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x6e, 0x69, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6e,
	0x6f, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63,
	0x65, 0x22, 0xbe, 0x01, 0x0a, 0x15, 0x52, 0x65, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x46, 0x65, 0x65,
	0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x47, 0x0a, 0x19, 0x63,
	0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x16, 0x63, 0x6c,
	0x6f, 0x73, 0x65, 0x46, 0x65, 0x65, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x12, 0x40, 0x0a, 0x16, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x66, 0x65,
	0x65, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x6f,
	0x66, 0x52, 0x13, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x46, 0x65, 0x65, 0x43, 0x72, 0x65, 0x64, 0x69,
	0x74, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69,
	0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69,
	0x6e, 0x6b, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x74, 0x78, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x66, 0x63, 0x3b, 0x66, 0x63,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fee_credit_txs_proto_rawDescOnce sync.Once
	file_fee_credit_txs_proto_rawDescData = file_fee_credit_txs_proto_rawDesc
)

func file_fee_credit_txs_proto_rawDescGZIP() []byte {
	file_fee_credit_txs_proto_rawDescOnce.Do(func() {
		file_fee_credit_txs_proto_rawDescData = protoimpl.X.CompressGZIP(file_fee_credit_txs_proto_rawDescData)
	})
	return file_fee_credit_txs_proto_rawDescData
}

var file_fee_credit_txs_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_fee_credit_txs_proto_goTypes = []interface{}{
	(*TransferFeeCreditOrder)(nil), // 0: TransferFeeCreditOrder
	(*AddFeeCreditOrder)(nil),      // 1: AddFeeCreditOrder
	(*CloseFeeCreditOrder)(nil),    // 2: CloseFeeCreditOrder
	(*ReclaimFeeCreditOrder)(nil),  // 3: ReclaimFeeCreditOrder
	(*txsystem.Transaction)(nil),   // 4: Transaction
	(*block.BlockProof)(nil),       // 5: BlockProof
}
var file_fee_credit_txs_proto_depIdxs = []int32{
	4, // 0: AddFeeCreditOrder.fee_credit_transfer:type_name -> Transaction
	5, // 1: AddFeeCreditOrder.fee_credit_transfer_proof:type_name -> BlockProof
	4, // 2: ReclaimFeeCreditOrder.close_fee_credit_transfer:type_name -> Transaction
	5, // 3: ReclaimFeeCreditOrder.close_fee_credit_proof:type_name -> BlockProof
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_fee_credit_txs_proto_init() }
func file_fee_credit_txs_proto_init() {
	if File_fee_credit_txs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fee_credit_txs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferFeeCreditOrder); i {
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
		file_fee_credit_txs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddFeeCreditOrder); i {
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
		file_fee_credit_txs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CloseFeeCreditOrder); i {
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
		file_fee_credit_txs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReclaimFeeCreditOrder); i {
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
			RawDescriptor: file_fee_credit_txs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fee_credit_txs_proto_goTypes,
		DependencyIndexes: file_fee_credit_txs_proto_depIdxs,
		MessageInfos:      file_fee_credit_txs_proto_msgTypes,
	}.Build()
	File_fee_credit_txs_proto = out.File
	file_fee_credit_txs_proto_rawDesc = nil
	file_fee_credit_txs_proto_goTypes = nil
	file_fee_credit_txs_proto_depIdxs = nil
}
