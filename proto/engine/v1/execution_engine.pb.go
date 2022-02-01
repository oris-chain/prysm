// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.15.8
// source: proto/engine/v1/execution_engine.proto

package enginev1

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/prysmaticlabs/prysm/proto/eth/ext"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PayloadStatus_Status int32

const (
	PayloadStatus_VALID                  PayloadStatus_Status = 0
	PayloadStatus_INVALID                PayloadStatus_Status = 1
	PayloadStatus_SYNCING                PayloadStatus_Status = 2
	PayloadStatus_ACCEPTED               PayloadStatus_Status = 3
	PayloadStatus_INVALID_BLOCK_HASH     PayloadStatus_Status = 4
	PayloadStatus_INVALID_TERMINAL_BLOCK PayloadStatus_Status = 5
)

// Enum value maps for PayloadStatus_Status.
var (
	PayloadStatus_Status_name = map[int32]string{
		0: "VALID",
		1: "INVALID",
		2: "SYNCING",
		3: "ACCEPTED",
		4: "INVALID_BLOCK_HASH",
		5: "INVALID_TERMINAL_BLOCK",
	}
	PayloadStatus_Status_value = map[string]int32{
		"VALID":                  0,
		"INVALID":                1,
		"SYNCING":                2,
		"ACCEPTED":               3,
		"INVALID_BLOCK_HASH":     4,
		"INVALID_TERMINAL_BLOCK": 5,
	}
)

func (x PayloadStatus_Status) Enum() *PayloadStatus_Status {
	p := new(PayloadStatus_Status)
	*p = x
	return p
}

func (x PayloadStatus_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PayloadStatus_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_engine_v1_execution_engine_proto_enumTypes[0].Descriptor()
}

func (PayloadStatus_Status) Type() protoreflect.EnumType {
	return &file_proto_engine_v1_execution_engine_proto_enumTypes[0]
}

func (x PayloadStatus_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PayloadStatus_Status.Descriptor instead.
func (PayloadStatus_Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_engine_v1_execution_engine_proto_rawDescGZIP(), []int{3, 0}
}

type ExecutionPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParentHash     []byte   `protobuf:"bytes,1,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty" ssz-size:"32"`
	FeeRecipient   []byte   `protobuf:"bytes,2,opt,name=fee_recipient,json=feeRecipient,proto3" json:"fee_recipient,omitempty" ssz-size:"20"`
	StateRoot      []byte   `protobuf:"bytes,3,opt,name=state_root,json=stateRoot,proto3" json:"state_root,omitempty" ssz-size:"32"`
	RecipientsRoot []byte   `protobuf:"bytes,4,opt,name=recipients_root,json=recipientsRoot,proto3" json:"recipients_root,omitempty" ssz-size:"32"`
	LogsBloom      []byte   `protobuf:"bytes,5,opt,name=logs_bloom,json=logsBloom,proto3" json:"logs_bloom,omitempty" ssz-size:"256"`
	Random         []byte   `protobuf:"bytes,6,opt,name=random,proto3" json:"random,omitempty" ssz-size:"32"`
	BlockNumber    uint64   `protobuf:"varint,7,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	GasLimit       uint64   `protobuf:"varint,8,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasUsed        uint64   `protobuf:"varint,9,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Timestamp      uint64   `protobuf:"varint,10,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ExtraData      []byte   `protobuf:"bytes,11,opt,name=extra_data,json=extraData,proto3" json:"extra_data,omitempty" ssz-max:"32"`
	BaseFeePerGas  []byte   `protobuf:"bytes,12,opt,name=base_fee_per_gas,json=baseFeePerGas,proto3" json:"base_fee_per_gas,omitempty" ssz-size:"32"`
	BlockHash      []byte   `protobuf:"bytes,13,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty" ssz-size:"32"`
	Transactions   [][]byte `protobuf:"bytes,14,rep,name=transactions,proto3" json:"transactions,omitempty" ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}

func (x *ExecutionPayload) Reset() {
	*x = ExecutionPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutionPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutionPayload) ProtoMessage() {}

func (x *ExecutionPayload) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutionPayload.ProtoReflect.Descriptor instead.
func (*ExecutionPayload) Descriptor() ([]byte, []int) {
	return file_proto_engine_v1_execution_engine_proto_rawDescGZIP(), []int{0}
}

func (x *ExecutionPayload) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *ExecutionPayload) GetFeeRecipient() []byte {
	if x != nil {
		return x.FeeRecipient
	}
	return nil
}

func (x *ExecutionPayload) GetStateRoot() []byte {
	if x != nil {
		return x.StateRoot
	}
	return nil
}

func (x *ExecutionPayload) GetRecipientsRoot() []byte {
	if x != nil {
		return x.RecipientsRoot
	}
	return nil
}

func (x *ExecutionPayload) GetLogsBloom() []byte {
	if x != nil {
		return x.LogsBloom
	}
	return nil
}

func (x *ExecutionPayload) GetRandom() []byte {
	if x != nil {
		return x.Random
	}
	return nil
}

func (x *ExecutionPayload) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *ExecutionPayload) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *ExecutionPayload) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *ExecutionPayload) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ExecutionPayload) GetExtraData() []byte {
	if x != nil {
		return x.ExtraData
	}
	return nil
}

func (x *ExecutionPayload) GetBaseFeePerGas() []byte {
	if x != nil {
		return x.BaseFeePerGas
	}
	return nil
}

func (x *ExecutionPayload) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *ExecutionPayload) GetTransactions() [][]byte {
	if x != nil {
		return x.Transactions
	}
	return nil
}

type ExecutionPayloadHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParentHash       []byte `protobuf:"bytes,1,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty" ssz-size:"32"`
	FeeRecipient     []byte `protobuf:"bytes,2,opt,name=fee_recipient,json=feeRecipient,proto3" json:"fee_recipient,omitempty" ssz-size:"20"`
	StateRoot        []byte `protobuf:"bytes,3,opt,name=state_root,json=stateRoot,proto3" json:"state_root,omitempty" ssz-size:"32"`
	RecipientsRoot   []byte `protobuf:"bytes,4,opt,name=recipients_root,json=recipientsRoot,proto3" json:"recipients_root,omitempty" ssz-size:"32"`
	LogsBloom        []byte `protobuf:"bytes,5,opt,name=logs_bloom,json=logsBloom,proto3" json:"logs_bloom,omitempty" ssz-size:"256"`
	Random           []byte `protobuf:"bytes,6,opt,name=random,proto3" json:"random,omitempty" ssz-size:"32"`
	BlockNumber      uint64 `protobuf:"varint,7,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	GasLimit         uint64 `protobuf:"varint,8,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasUsed          uint64 `protobuf:"varint,9,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Timestamp        uint64 `protobuf:"varint,10,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ExtraData        []byte `protobuf:"bytes,11,opt,name=extra_data,json=extraData,proto3" json:"extra_data,omitempty" ssz-max:"32"`
	BaseFeePerGas    []byte `protobuf:"bytes,12,opt,name=base_fee_per_gas,json=baseFeePerGas,proto3" json:"base_fee_per_gas,omitempty" ssz-size:"32"`
	BlockHash        []byte `protobuf:"bytes,13,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty" ssz-size:"32"`
	TransactionsRoot []byte `protobuf:"bytes,14,opt,name=transactions_root,json=transactionsRoot,proto3" json:"transactions_root,omitempty" ssz-size:"32"`
}

func (x *ExecutionPayloadHeader) Reset() {
	*x = ExecutionPayloadHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutionPayloadHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutionPayloadHeader) ProtoMessage() {}

func (x *ExecutionPayloadHeader) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutionPayloadHeader.ProtoReflect.Descriptor instead.
func (*ExecutionPayloadHeader) Descriptor() ([]byte, []int) {
	return file_proto_engine_v1_execution_engine_proto_rawDescGZIP(), []int{1}
}

func (x *ExecutionPayloadHeader) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetFeeRecipient() []byte {
	if x != nil {
		return x.FeeRecipient
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetStateRoot() []byte {
	if x != nil {
		return x.StateRoot
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetRecipientsRoot() []byte {
	if x != nil {
		return x.RecipientsRoot
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetLogsBloom() []byte {
	if x != nil {
		return x.LogsBloom
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetRandom() []byte {
	if x != nil {
		return x.Random
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *ExecutionPayloadHeader) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *ExecutionPayloadHeader) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *ExecutionPayloadHeader) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ExecutionPayloadHeader) GetExtraData() []byte {
	if x != nil {
		return x.ExtraData
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetBaseFeePerGas() []byte {
	if x != nil {
		return x.BaseFeePerGas
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *ExecutionPayloadHeader) GetTransactionsRoot() []byte {
	if x != nil {
		return x.TransactionsRoot
	}
	return nil
}

type PayloadAttributes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp             uint64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Random                []byte `protobuf:"bytes,2,opt,name=random,proto3" json:"random,omitempty" ssz-size:"32"`
	SuggestedFeeRecipient []byte `protobuf:"bytes,3,opt,name=suggested_fee_recipient,json=suggestedFeeRecipient,proto3" json:"suggested_fee_recipient,omitempty" ssz-size:"20"`
}

func (x *PayloadAttributes) Reset() {
	*x = PayloadAttributes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayloadAttributes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayloadAttributes) ProtoMessage() {}

func (x *PayloadAttributes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayloadAttributes.ProtoReflect.Descriptor instead.
func (*PayloadAttributes) Descriptor() ([]byte, []int) {
	return file_proto_engine_v1_execution_engine_proto_rawDescGZIP(), []int{2}
}

func (x *PayloadAttributes) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *PayloadAttributes) GetRandom() []byte {
	if x != nil {
		return x.Random
	}
	return nil
}

func (x *PayloadAttributes) GetSuggestedFeeRecipient() []byte {
	if x != nil {
		return x.SuggestedFeeRecipient
	}
	return nil
}

type PayloadStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status          PayloadStatus_Status `protobuf:"varint,1,opt,name=status,proto3,enum=ethereum.engine.v1.PayloadStatus_Status" json:"status,omitempty"`
	LatestValidHash []byte               `protobuf:"bytes,2,opt,name=latest_valid_hash,json=latestValidHash,proto3" json:"latest_valid_hash,omitempty" ssz-size:"32"`
	ValidationError string               `protobuf:"bytes,3,opt,name=validation_error,json=validationError,proto3" json:"validation_error,omitempty"`
}

func (x *PayloadStatus) Reset() {
	*x = PayloadStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayloadStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayloadStatus) ProtoMessage() {}

func (x *PayloadStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayloadStatus.ProtoReflect.Descriptor instead.
func (*PayloadStatus) Descriptor() ([]byte, []int) {
	return file_proto_engine_v1_execution_engine_proto_rawDescGZIP(), []int{3}
}

func (x *PayloadStatus) GetStatus() PayloadStatus_Status {
	if x != nil {
		return x.Status
	}
	return PayloadStatus_VALID
}

func (x *PayloadStatus) GetLatestValidHash() []byte {
	if x != nil {
		return x.LatestValidHash
	}
	return nil
}

func (x *PayloadStatus) GetValidationError() string {
	if x != nil {
		return x.ValidationError
	}
	return ""
}

type ForkchoiceState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HeadBlockHash      []byte `protobuf:"bytes,1,opt,name=head_block_hash,json=headBlockHash,proto3" json:"head_block_hash,omitempty" ssz-size:"32"`
	SafeBlockHash      []byte `protobuf:"bytes,2,opt,name=safe_block_hash,json=safeBlockHash,proto3" json:"safe_block_hash,omitempty" ssz-size:"32"`
	FinalizedBlockHash []byte `protobuf:"bytes,3,opt,name=finalized_block_hash,json=finalizedBlockHash,proto3" json:"finalized_block_hash,omitempty" ssz-size:"32"`
}

func (x *ForkchoiceState) Reset() {
	*x = ForkchoiceState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForkchoiceState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForkchoiceState) ProtoMessage() {}

func (x *ForkchoiceState) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_v1_execution_engine_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForkchoiceState.ProtoReflect.Descriptor instead.
func (*ForkchoiceState) Descriptor() ([]byte, []int) {
	return file_proto_engine_v1_execution_engine_proto_rawDescGZIP(), []int{4}
}

func (x *ForkchoiceState) GetHeadBlockHash() []byte {
	if x != nil {
		return x.HeadBlockHash
	}
	return nil
}

func (x *ForkchoiceState) GetSafeBlockHash() []byte {
	if x != nil {
		return x.SafeBlockHash
	}
	return nil
}

func (x *ForkchoiceState) GetFinalizedBlockHash() []byte {
	if x != nil {
		return x.FinalizedBlockHash
	}
	return nil
}

var File_proto_engine_v1_execution_engine_proto protoreflect.FileDescriptor

var file_proto_engine_v1_execution_engine_proto_rawDesc = []byte{
	0x0a, 0x26, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65,
	0x75, 0x6d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x65, 0x78, 0x74, 0x2f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc3, 0x04, 0x0a, 0x10, 0x45, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x27,
	0x0a, 0x0b, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0a, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x2b, 0x0a, 0x0d, 0x66, 0x65, 0x65, 0x5f, 0x72,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06,
	0x8a, 0xb5, 0x18, 0x02, 0x32, 0x30, 0x52, 0x0c, 0x66, 0x65, 0x65, 0x52, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x6f,
	0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32,
	0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x2f, 0x0a, 0x0f, 0x72,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0e, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x26, 0x0a, 0x0a,
	0x6c, 0x6f, 0x67, 0x73, 0x5f, 0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x42, 0x07, 0x8a, 0xb5, 0x18, 0x03, 0x32, 0x35, 0x36, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x73, 0x42,
	0x6c, 0x6f, 0x6f, 0x6d, 0x12, 0x1e, 0x0a, 0x06, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x06, 0x72, 0x61,
	0x6e, 0x64, 0x6f, 0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x67, 0x61, 0x73, 0x4c,
	0x69, 0x6d, 0x69, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x61, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x64,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x25, 0x0a,
	0x0a, 0x65, 0x78, 0x74, 0x72, 0x61, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0c, 0x42, 0x06, 0x92, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x2f, 0x0a, 0x10, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x66, 0x65, 0x65,
	0x5f, 0x70, 0x65, 0x72, 0x5f, 0x67, 0x61, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06,
	0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0d, 0x62, 0x61, 0x73, 0x65, 0x46, 0x65, 0x65, 0x50,
	0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x25, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33,
	0x32, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x41, 0x0a, 0x0c,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0e, 0x20, 0x03,
	0x28, 0x0c, 0x42, 0x1d, 0x8a, 0xb5, 0x18, 0x03, 0x3f, 0x2c, 0x3f, 0x92, 0xb5, 0x18, 0x12, 0x31,
	0x30, 0x34, 0x38, 0x35, 0x37, 0x36, 0x2c, 0x31, 0x30, 0x37, 0x33, 0x37, 0x34, 0x31, 0x38, 0x32,
	0x34, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0xbb, 0x04, 0x0a, 0x16, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0b, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x42,
	0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x2b, 0x0a, 0x0d, 0x66, 0x65, 0x65, 0x5f, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02,
	0x32, 0x30, 0x52, 0x0c, 0x66, 0x65, 0x65, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x12, 0x25, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x2f, 0x0a, 0x0f, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x26, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x73,
	0x5f, 0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07, 0x8a, 0xb5,
	0x18, 0x03, 0x32, 0x35, 0x36, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x73, 0x42, 0x6c, 0x6f, 0x6f, 0x6d,
	0x12, 0x1e, 0x0a, 0x06, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c,
	0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x06, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x67, 0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x67, 0x61, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x25, 0x0a, 0x0a, 0x65, 0x78, 0x74,
	0x72, 0x61, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x92,
	0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x2f, 0x0a, 0x10, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x72,
	0x5f, 0x67, 0x61, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02,
	0x33, 0x32, 0x52, 0x0d, 0x62, 0x61, 0x73, 0x65, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61,
	0x73, 0x12, 0x25, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x09, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x33, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x10, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x22, 0x91, 0x01,
	0x0a, 0x11, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x1e, 0x0a, 0x06, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x06, 0x72, 0x61, 0x6e, 0x64, 0x6f,
	0x6d, 0x12, 0x3e, 0x0a, 0x17, 0x73, 0x75, 0x67, 0x67, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x66,
	0x65, 0x65, 0x5f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x32, 0x30, 0x52, 0x15, 0x73, 0x75, 0x67, 0x67,
	0x65, 0x73, 0x74, 0x65, 0x64, 0x46, 0x65, 0x65, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x22, 0xa1, 0x02, 0x0a, 0x0d, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x40, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x32, 0x0a, 0x11, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x5f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0f, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x48, 0x61, 0x73, 0x68, 0x12, 0x29, 0x0a, 0x10, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x6f, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x09,
	0x0a, 0x05, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x59, 0x4e, 0x43, 0x49, 0x4e,
	0x47, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x43, 0x43, 0x45, 0x50, 0x54, 0x45, 0x44, 0x10,
	0x03, 0x12, 0x16, 0x0a, 0x12, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x42, 0x4c, 0x4f,
	0x43, 0x4b, 0x5f, 0x48, 0x41, 0x53, 0x48, 0x10, 0x04, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x5f, 0x54, 0x45, 0x52, 0x4d, 0x49, 0x4e, 0x41, 0x4c, 0x5f, 0x42, 0x4c,
	0x4f, 0x43, 0x4b, 0x10, 0x05, 0x22, 0xab, 0x01, 0x0a, 0x0f, 0x46, 0x6f, 0x72, 0x6b, 0x63, 0x68,
	0x6f, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x2e, 0x0a, 0x0f, 0x68, 0x65, 0x61,
	0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0d, 0x68, 0x65, 0x61, 0x64,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x2e, 0x0a, 0x0f, 0x73, 0x61, 0x66,
	0x65, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0d, 0x73, 0x61, 0x66, 0x65,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x38, 0x0a, 0x14, 0x66, 0x69, 0x6e,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52,
	0x12, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x61, 0x73, 0x68, 0x42, 0x93, 0x01, 0x0a, 0x16, 0x6f, 0x72, 0x67, 0x2e, 0x65, 0x74, 0x68, 0x65,
	0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x14,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6c, 0x61, 0x62, 0x73,
	0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x76, 0x31, 0xaa,
	0x02, 0x12, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x45, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x12, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x5c,
	0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5c, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_engine_v1_execution_engine_proto_rawDescOnce sync.Once
	file_proto_engine_v1_execution_engine_proto_rawDescData = file_proto_engine_v1_execution_engine_proto_rawDesc
)

func file_proto_engine_v1_execution_engine_proto_rawDescGZIP() []byte {
	file_proto_engine_v1_execution_engine_proto_rawDescOnce.Do(func() {
		file_proto_engine_v1_execution_engine_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_engine_v1_execution_engine_proto_rawDescData)
	})
	return file_proto_engine_v1_execution_engine_proto_rawDescData
}

var file_proto_engine_v1_execution_engine_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_engine_v1_execution_engine_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_engine_v1_execution_engine_proto_goTypes = []interface{}{
	(PayloadStatus_Status)(0),      // 0: ethereum.engine.v1.PayloadStatus.Status
	(*ExecutionPayload)(nil),       // 1: ethereum.engine.v1.ExecutionPayload
	(*ExecutionPayloadHeader)(nil), // 2: ethereum.engine.v1.ExecutionPayloadHeader
	(*PayloadAttributes)(nil),      // 3: ethereum.engine.v1.PayloadAttributes
	(*PayloadStatus)(nil),          // 4: ethereum.engine.v1.PayloadStatus
	(*ForkchoiceState)(nil),        // 5: ethereum.engine.v1.ForkchoiceState
}
var file_proto_engine_v1_execution_engine_proto_depIdxs = []int32{
	0, // 0: ethereum.engine.v1.PayloadStatus.status:type_name -> ethereum.engine.v1.PayloadStatus.Status
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_engine_v1_execution_engine_proto_init() }
func file_proto_engine_v1_execution_engine_proto_init() {
	if File_proto_engine_v1_execution_engine_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_engine_v1_execution_engine_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecutionPayload); i {
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
		file_proto_engine_v1_execution_engine_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecutionPayloadHeader); i {
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
		file_proto_engine_v1_execution_engine_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayloadAttributes); i {
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
		file_proto_engine_v1_execution_engine_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayloadStatus); i {
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
		file_proto_engine_v1_execution_engine_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForkchoiceState); i {
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
			RawDescriptor: file_proto_engine_v1_execution_engine_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_engine_v1_execution_engine_proto_goTypes,
		DependencyIndexes: file_proto_engine_v1_execution_engine_proto_depIdxs,
		EnumInfos:         file_proto_engine_v1_execution_engine_proto_enumTypes,
		MessageInfos:      file_proto_engine_v1_execution_engine_proto_msgTypes,
	}.Build()
	File_proto_engine_v1_execution_engine_proto = out.File
	file_proto_engine_v1_execution_engine_proto_rawDesc = nil
	file_proto_engine_v1_execution_engine_proto_goTypes = nil
	file_proto_engine_v1_execution_engine_proto_depIdxs = nil
}
