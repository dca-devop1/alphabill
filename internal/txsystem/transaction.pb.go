// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: transaction.proto

package txsystem

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

// Transaction is a generic transaction order, same for all transaction systems.
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SystemId              []byte     `protobuf:"bytes,1,opt,name=system_id,json=systemId,proto3" json:"system_id,omitempty"`
	UnitId                []byte     `protobuf:"bytes,2,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	TransactionAttributes *anypb.Any `protobuf:"bytes,3,opt,name=transaction_attributes,json=transactionAttributes,proto3" json:"transaction_attributes,omitempty"`
	OwnerProof            []byte     `protobuf:"bytes,5,opt,name=owner_proof,json=ownerProof,proto3" json:"owner_proof,omitempty"`
	// optional fee authorization proof (omitted when the main owner proof also satisfies the fee owner condition)
	FeeProof []byte `protobuf:"bytes,6,opt,name=fee_proof,json=feeProof,proto3" json:"fee_proof,omitempty"`
	// transaction metadata added by user
	ClientMetadata *ClientMetadata `protobuf:"bytes,7,opt,name=client_metadata,json=clientMetadata,proto3" json:"client_metadata,omitempty"`
	// transaction metadata added by node
	ServerMetadata *ServerMetadata `protobuf:"bytes,8,opt,name=server_metadata,json=serverMetadata,proto3" json:"server_metadata,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[0]
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
	return file_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetSystemId() []byte {
	if x != nil {
		return x.SystemId
	}
	return nil
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

func (x *Transaction) GetOwnerProof() []byte {
	if x != nil {
		return x.OwnerProof
	}
	return nil
}

func (x *Transaction) GetFeeProof() []byte {
	if x != nil {
		return x.FeeProof
	}
	return nil
}

func (x *Transaction) GetClientMetadata() *ClientMetadata {
	if x != nil {
		return x.ClientMetadata
	}
	return nil
}

func (x *Transaction) GetServerMetadata() *ServerMetadata {
	if x != nil {
		return x.ServerMetadata
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
		mi := &file_transaction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionResponse) ProtoMessage() {}

func (x *TransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[1]
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
	return file_transaction_proto_rawDescGZIP(), []int{1}
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

type ClientMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// transaction timeout block height
	Timeout uint64 `protobuf:"varint,1,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// maximum fee the user is willing to pay for the execution of this transaction
	MaxFee uint64 `protobuf:"varint,2,opt,name=max_fee,json=maxFee,proto3" json:"max_fee,omitempty"`
	// fee credit record identifier (should be empty for fee credit transactions)
	FeeCreditRecordId []byte `protobuf:"bytes,3,opt,name=fee_credit_record_id,json=feeCreditRecordId,proto3" json:"fee_credit_record_id,omitempty"`
}

func (x *ClientMetadata) Reset() {
	*x = ClientMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientMetadata) ProtoMessage() {}

func (x *ClientMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientMetadata.ProtoReflect.Descriptor instead.
func (*ClientMetadata) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{2}
}

func (x *ClientMetadata) GetTimeout() uint64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *ClientMetadata) GetMaxFee() uint64 {
	if x != nil {
		return x.MaxFee
	}
	return 0
}

func (x *ClientMetadata) GetFeeCreditRecordId() []byte {
	if x != nil {
		return x.FeeCreditRecordId
	}
	return nil
}

type ServerMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// actual transaction fee charged
	Fee uint64 `protobuf:"varint,1,opt,name=fee,proto3" json:"fee,omitempty"`
}

func (x *ServerMetadata) Reset() {
	*x = ServerMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMetadata) ProtoMessage() {}

func (x *ServerMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMetadata.ProtoReflect.Descriptor instead.
func (*ServerMetadata) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{3}
}

func (x *ServerMetadata) GetFee() uint64 {
	if x != nil {
		return x.Fee
	}
	return 0
}

type Transactions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transactions []*Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
}

func (x *Transactions) Reset() {
	*x = Transactions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transactions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transactions) ProtoMessage() {}

func (x *Transactions) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transactions.ProtoReflect.Descriptor instead.
func (*Transactions) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{4}
}

func (x *Transactions) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

var File_transaction_proto protoreflect.FileDescriptor

var file_transaction_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc2,
	0x02, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b,
	0x0a, 0x09, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x08, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x75, 0x6e,
	0x69, 0x74, 0x49, 0x64, 0x12, 0x4b, 0x0a, 0x16, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x15, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x50, 0x72, 0x6f,
	0x6f, 0x66, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x66, 0x65, 0x65, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12,
	0x38, 0x0a, 0x0f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x38, 0x0a, 0x0f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x3f, 0x0a, 0x13, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x74, 0x0a, 0x0e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x6d, 0x61, 0x78, 0x5f, 0x66, 0x65, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x6d, 0x61, 0x78, 0x46, 0x65, 0x65, 0x12, 0x2f, 0x0a, 0x14, 0x66, 0x65, 0x65,
	0x5f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x66, 0x65, 0x65, 0x43, 0x72, 0x65, 0x64,
	0x69, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x22, 0x22, 0x0a, 0x0e, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03,
	0x66, 0x65, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x66, 0x65, 0x65, 0x22, 0x40,
	0x0a, 0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x30,
	0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x62, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x74, 0x78, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x3b, 0x74, 0x78, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transaction_proto_rawDescOnce sync.Once
	file_transaction_proto_rawDescData = file_transaction_proto_rawDesc
)

func file_transaction_proto_rawDescGZIP() []byte {
	file_transaction_proto_rawDescOnce.Do(func() {
		file_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_transaction_proto_rawDescData)
	})
	return file_transaction_proto_rawDescData
}

var file_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_transaction_proto_goTypes = []interface{}{
	(*Transaction)(nil),         // 0: Transaction
	(*TransactionResponse)(nil), // 1: TransactionResponse
	(*ClientMetadata)(nil),      // 2: ClientMetadata
	(*ServerMetadata)(nil),      // 3: ServerMetadata
	(*Transactions)(nil),        // 4: Transactions
	(*anypb.Any)(nil),           // 5: google.protobuf.Any
}
var file_transaction_proto_depIdxs = []int32{
	5, // 0: Transaction.transaction_attributes:type_name -> google.protobuf.Any
	2, // 1: Transaction.client_metadata:type_name -> ClientMetadata
	3, // 2: Transaction.server_metadata:type_name -> ServerMetadata
	0, // 3: Transactions.transactions:type_name -> Transaction
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_transaction_proto_init() }
func file_transaction_proto_init() {
	if File_transaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_transaction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_transaction_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_transaction_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMetadata); i {
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
		file_transaction_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMetadata); i {
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
		file_transaction_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transactions); i {
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
			RawDescriptor: file_transaction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_transaction_proto_goTypes,
		DependencyIndexes: file_transaction_proto_depIdxs,
		MessageInfos:      file_transaction_proto_msgTypes,
	}.Build()
	File_transaction_proto = out.File
	file_transaction_proto_rawDesc = nil
	file_transaction_proto_goTypes = nil
	file_transaction_proto_depIdxs = nil
}
