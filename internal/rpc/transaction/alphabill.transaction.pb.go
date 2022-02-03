// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: alphabill.transaction.proto

package transaction

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TransactionOrder is a generic transaction order, same for all transaction systems.
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnitId                []byte     `protobuf:"bytes,1,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	TransactionAttributes *anypb.Any `protobuf:"bytes,2,opt,name=transaction_attributes,json=transactionAttributes,proto3" json:"transaction_attributes,omitempty"`
	Timeout               uint64     `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	OwnerProof            []byte     `protobuf:"bytes,4,opt,name=owner_proof,json=ownerProof,proto3" json:"owner_proof,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_transaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_transaction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_alphabill_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetUnitId() []byte {
	if x != nil {
		return x.UnitId
	}
	return nil
}

func (x *Transaction) GetTransactionAttributes() *anypb.Any {
	if x != nil {
		return x.TransactionAttributes
	}
	return nil
}

func (x *Transaction) GetTimeout() uint64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *Transaction) GetOwnerProof() []byte {
	if x != nil {
		return x.OwnerProof
	}
	return nil
}

type TransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// True if request passed initial validation.
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	// Contains error message if ok is false.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TransactionResponse) Reset() {
	*x = TransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_transaction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionResponse) ProtoMessage() {}

func (x *TransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_transaction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionResponse.ProtoReflect.Descriptor instead.
func (*TransactionResponse) Descriptor() ([]byte, []int) {
	return file_alphabill_transaction_proto_rawDescGZIP(), []int{1}
}

func (x *TransactionResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *TransactionResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Alphabill specific transaction attributes.
type BillTransfer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewBearer   []byte `protobuf:"bytes,1,opt,name=new_bearer,json=newBearer,proto3" json:"new_bearer,omitempty"`
	TargetValue uint64 `protobuf:"varint,2,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty"`
	Backlink    []byte `protobuf:"bytes,3,opt,name=backlink,proto3" json:"backlink,omitempty"`
}

func (x *BillTransfer) Reset() {
	*x = BillTransfer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_transaction_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BillTransfer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BillTransfer) ProtoMessage() {}

func (x *BillTransfer) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_transaction_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BillTransfer.ProtoReflect.Descriptor instead.
func (*BillTransfer) Descriptor() ([]byte, []int) {
	return file_alphabill_transaction_proto_rawDescGZIP(), []int{2}
}

func (x *BillTransfer) GetNewBearer() []byte {
	if x != nil {
		return x.NewBearer
	}
	return nil
}

func (x *BillTransfer) GetTargetValue() uint64 {
	if x != nil {
		return x.TargetValue
	}
	return 0
}

func (x *BillTransfer) GetBacklink() []byte {
	if x != nil {
		return x.Backlink
	}
	return nil
}

type TransferDC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce        []byte `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	TargetBearer []byte `protobuf:"bytes,2,opt,name=target_bearer,json=targetBearer,proto3" json:"target_bearer,omitempty"`
	TargetValue  uint64 `protobuf:"varint,3,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty"`
	Backlink     []byte `protobuf:"bytes,4,opt,name=backlink,proto3" json:"backlink,omitempty"`
}

func (x *TransferDC) Reset() {
	*x = TransferDC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_transaction_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferDC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferDC) ProtoMessage() {}

func (x *TransferDC) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_transaction_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferDC.ProtoReflect.Descriptor instead.
func (*TransferDC) Descriptor() ([]byte, []int) {
	return file_alphabill_transaction_proto_rawDescGZIP(), []int{3}
}

func (x *TransferDC) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *TransferDC) GetTargetBearer() []byte {
	if x != nil {
		return x.TargetBearer
	}
	return nil
}

func (x *TransferDC) GetTargetValue() uint64 {
	if x != nil {
		return x.TargetValue
	}
	return 0
}

func (x *TransferDC) GetBacklink() []byte {
	if x != nil {
		return x.Backlink
	}
	return nil
}

type BillSplit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount         uint64 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	TargetBearer   []byte `protobuf:"bytes,2,opt,name=target_bearer,json=targetBearer,proto3" json:"target_bearer,omitempty"`
	RemainingValue uint64 `protobuf:"varint,3,opt,name=remaining_value,json=remainingValue,proto3" json:"remaining_value,omitempty"`
	Backlink       []byte `protobuf:"bytes,4,opt,name=backlink,proto3" json:"backlink,omitempty"`
}

func (x *BillSplit) Reset() {
	*x = BillSplit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_transaction_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BillSplit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BillSplit) ProtoMessage() {}

func (x *BillSplit) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_transaction_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BillSplit.ProtoReflect.Descriptor instead.
func (*BillSplit) Descriptor() ([]byte, []int) {
	return file_alphabill_transaction_proto_rawDescGZIP(), []int{4}
}

func (x *BillSplit) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *BillSplit) GetTargetBearer() []byte {
	if x != nil {
		return x.TargetBearer
	}
	return nil
}

func (x *BillSplit) GetRemainingValue() uint64 {
	if x != nil {
		return x.RemainingValue
	}
	return 0
}

func (x *BillSplit) GetBacklink() []byte {
	if x != nil {
		return x.Backlink
	}
	return nil
}

type Swap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerCondition  []byte         `protobuf:"bytes,1,opt,name=owner_condition,json=ownerCondition,proto3" json:"owner_condition,omitempty"`
	BillIdentifiers [][]byte       `protobuf:"bytes,2,rep,name=bill_identifiers,json=billIdentifiers,proto3" json:"bill_identifiers,omitempty"`
	DcTransfers     []*Transaction `protobuf:"bytes,3,rep,name=dc_transfers,json=dcTransfers,proto3" json:"dc_transfers,omitempty"`
	Proofs          [][]byte       `protobuf:"bytes,4,rep,name=proofs,proto3" json:"proofs,omitempty"`
	TargetValue     uint64         `protobuf:"varint,5,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty"`
}

func (x *Swap) Reset() {
	*x = Swap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alphabill_transaction_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Swap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Swap) ProtoMessage() {}

func (x *Swap) ProtoReflect() protoreflect.Message {
	mi := &file_alphabill_transaction_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Swap.ProtoReflect.Descriptor instead.
func (*Swap) Descriptor() ([]byte, []int) {
	return file_alphabill_transaction_proto_rawDescGZIP(), []int{5}
}

func (x *Swap) GetOwnerCondition() []byte {
	if x != nil {
		return x.OwnerCondition
	}
	return nil
}

func (x *Swap) GetBillIdentifiers() [][]byte {
	if x != nil {
		return x.BillIdentifiers
	}
	return nil
}

func (x *Swap) GetDcTransfers() []*Transaction {
	if x != nil {
		return x.DcTransfers
	}
	return nil
}

func (x *Swap) GetProofs() [][]byte {
	if x != nil {
		return x.Proofs
	}
	return nil
}

func (x *Swap) GetTargetValue() uint64 {
	if x != nil {
		return x.TargetValue
	}
	return 0
}

var File_alphabill_transaction_proto protoreflect.FileDescriptor

var file_alphabill_transaction_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2e, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61,
	0x62, 0x72, 0x70, 0x63, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xae, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49, 0x64, 0x12, 0x4b, 0x0a, 0x16, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x15,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x6f, 0x66,
	0x22, 0x3f, 0x0a, 0x13, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x6c, 0x0a, 0x0c, 0x42, 0x69, 0x6c, 0x6c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65,
	0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x77, 0x5f, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x6e, 0x65, 0x77, 0x42, 0x65, 0x61, 0x72, 0x65, 0x72,
	0x12, 0x21, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x22,
	0x86, 0x01, 0x0a, 0x0a, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x44, 0x43, 0x12, 0x14,
	0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6e,
	0x6f, 0x6e, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x62,
	0x65, 0x61, 0x72, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x42, 0x65, 0x61, 0x72, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08,
	0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x22, 0x8d, 0x01, 0x0a, 0x09, 0x42, 0x69, 0x6c,
	0x6c, 0x53, 0x70, 0x6c, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23,
	0x0a, 0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x42, 0x65, 0x61,
	0x72, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x72, 0x65,
	0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08,
	0x62, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x6e, 0x6b, 0x22, 0xcc, 0x01, 0x0a, 0x04, 0x53, 0x77, 0x61,
	0x70, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x62, 0x69,
	0x6c, 0x6c, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0c, 0x52, 0x0f, 0x62, 0x69, 0x6c, 0x6c, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x73, 0x12, 0x35, 0x0a, 0x0c, 0x64, 0x63, 0x5f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x62,
	0x72, 0x70, 0x63, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0b, 0x64, 0x63, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x72, 0x6f, 0x6f, 0x66, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x72,
	0x6f, 0x6f, 0x66, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x56, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x46, 0x0a, 0x12, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x2e,
	0x61, 0x62, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x1a, 0x1a, 0x2e, 0x61, 0x62, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x5c, 0x5a, 0x5a, 0x67, 0x69, 0x74, 0x64, 0x63, 0x2e, 0x65, 0x65, 0x2e, 0x67, 0x75, 0x61, 0x72,
	0x64, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62,
	0x69, 0x6c, 0x6c, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x3b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_alphabill_transaction_proto_rawDescOnce sync.Once
	file_alphabill_transaction_proto_rawDescData = file_alphabill_transaction_proto_rawDesc
)

func file_alphabill_transaction_proto_rawDescGZIP() []byte {
	file_alphabill_transaction_proto_rawDescOnce.Do(func() {
		file_alphabill_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_alphabill_transaction_proto_rawDescData)
	})
	return file_alphabill_transaction_proto_rawDescData
}

var file_alphabill_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_alphabill_transaction_proto_goTypes = []interface{}{
	(*Transaction)(nil),         // 0: abrpc.Transaction
	(*TransactionResponse)(nil), // 1: abrpc.TransactionResponse
	(*BillTransfer)(nil),        // 2: abrpc.BillTransfer
	(*TransferDC)(nil),          // 3: abrpc.TransferDC
	(*BillSplit)(nil),           // 4: abrpc.BillSplit
	(*Swap)(nil),                // 5: abrpc.Swap
	(*anypb.Any)(nil),           // 6: google.protobuf.Any
}
var file_alphabill_transaction_proto_depIdxs = []int32{
	6, // 0: abrpc.Transaction.transaction_attributes:type_name -> google.protobuf.Any
	0, // 1: abrpc.Swap.dc_transfers:type_name -> abrpc.Transaction
	0, // 2: abrpc.Transactions.ProcessTransaction:input_type -> abrpc.Transaction
	1, // 3: abrpc.Transactions.ProcessTransaction:output_type -> abrpc.TransactionResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_alphabill_transaction_proto_init() }
func file_alphabill_transaction_proto_init() {
	if File_alphabill_transaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_alphabill_transaction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
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
		file_alphabill_transaction_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransactionResponse); i {
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
		file_alphabill_transaction_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BillTransfer); i {
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
		file_alphabill_transaction_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferDC); i {
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
		file_alphabill_transaction_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BillSplit); i {
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
		file_alphabill_transaction_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Swap); i {
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
			RawDescriptor: file_alphabill_transaction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_alphabill_transaction_proto_goTypes,
		DependencyIndexes: file_alphabill_transaction_proto_depIdxs,
		MessageInfos:      file_alphabill_transaction_proto_msgTypes,
	}.Build()
	File_alphabill_transaction_proto = out.File
	file_alphabill_transaction_proto_rawDesc = nil
	file_alphabill_transaction_proto_goTypes = nil
	file_alphabill_transaction_proto_depIdxs = nil
}
