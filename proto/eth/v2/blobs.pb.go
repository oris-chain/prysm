// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.15.8
// source: proto/eth/v2/blobs.proto

package eth

import (
	reflect "reflect"
	sync "sync"

	github_com_prysmaticlabs_prysm_v4_consensus_types_primitives "github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	_ "github.com/prysmaticlabs/prysm/v4/proto/eth/ext"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BlobSidecars struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sidecars []*BlobSidecar `protobuf:"bytes,1,rep,name=sidecars,proto3" json:"sidecars,omitempty"`
}

func (x *BlobSidecars) Reset() {
	*x = BlobSidecars{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eth_v2_blobs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlobSidecars) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlobSidecars) ProtoMessage() {}

func (x *BlobSidecars) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eth_v2_blobs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlobSidecars.ProtoReflect.Descriptor instead.
func (*BlobSidecars) Descriptor() ([]byte, []int) {
	return file_proto_eth_v2_blobs_proto_rawDescGZIP(), []int{0}
}

func (x *BlobSidecars) GetSidecars() []*BlobSidecar {
	if x != nil {
		return x.Sidecars
	}
	return nil
}

type BlobSidecar struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockRoot       []byte                                                                      `protobuf:"bytes,1,opt,name=block_root,json=blockRoot,proto3" json:"block_root,omitempty" ssz-size:"32"`
	Index           uint64                                                                      `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Slot            github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.Slot           `protobuf:"varint,3,opt,name=slot,proto3" json:"slot,omitempty" cast-type:"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives.Slot"`
	BlockParentRoot []byte                                                                      `protobuf:"bytes,4,opt,name=block_parent_root,json=blockParentRoot,proto3" json:"block_parent_root,omitempty" ssz-size:"32"`
	ProposerIndex   github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.ValidatorIndex `protobuf:"varint,5,opt,name=proposer_index,json=proposerIndex,proto3" json:"proposer_index,omitempty" cast-type:"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives.ValidatorIndex"`
	Blob            []byte                                                                      `protobuf:"bytes,6,opt,name=blob,proto3" json:"blob,omitempty" ssz-size:"131072"`
	KzgCommitment   []byte                                                                      `protobuf:"bytes,7,opt,name=kzg_commitment,json=kzgCommitment,proto3" json:"kzg_commitment,omitempty" ssz-size:"48"`
	KzgProof        []byte                                                                      `protobuf:"bytes,8,opt,name=kzg_proof,json=kzgProof,proto3" json:"kzg_proof,omitempty" ssz-size:"48"`
}

func (x *BlobSidecar) Reset() {
	*x = BlobSidecar{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eth_v2_blobs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlobSidecar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlobSidecar) ProtoMessage() {}

func (x *BlobSidecar) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eth_v2_blobs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlobSidecar.ProtoReflect.Descriptor instead.
func (*BlobSidecar) Descriptor() ([]byte, []int) {
	return file_proto_eth_v2_blobs_proto_rawDescGZIP(), []int{1}
}

func (x *BlobSidecar) GetBlockRoot() []byte {
	if x != nil {
		return x.BlockRoot
	}
	return nil
}

func (x *BlobSidecar) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *BlobSidecar) GetSlot() github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.Slot {
	if x != nil {
		return x.Slot
	}
	return github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.Slot(0)
}

func (x *BlobSidecar) GetBlockParentRoot() []byte {
	if x != nil {
		return x.BlockParentRoot
	}
	return nil
}

func (x *BlobSidecar) GetProposerIndex() github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.ValidatorIndex {
	if x != nil {
		return x.ProposerIndex
	}
	return github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.ValidatorIndex(0)
}

func (x *BlobSidecar) GetBlob() []byte {
	if x != nil {
		return x.Blob
	}
	return nil
}

func (x *BlobSidecar) GetKzgCommitment() []byte {
	if x != nil {
		return x.KzgCommitment
	}
	return nil
}

func (x *BlobSidecar) GetKzgProof() []byte {
	if x != nil {
		return x.KzgProof
	}
	return nil
}

type SignedBlobSidecar struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message   *BlobSidecar `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Signature []byte       `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty" ssz-size:"96"`
}

func (x *SignedBlobSidecar) Reset() {
	*x = SignedBlobSidecar{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eth_v2_blobs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedBlobSidecar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedBlobSidecar) ProtoMessage() {}

func (x *SignedBlobSidecar) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eth_v2_blobs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedBlobSidecar.ProtoReflect.Descriptor instead.
func (*SignedBlobSidecar) Descriptor() ([]byte, []int) {
	return file_proto_eth_v2_blobs_proto_rawDescGZIP(), []int{2}
}

func (x *SignedBlobSidecar) GetMessage() *BlobSidecar {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *SignedBlobSidecar) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type BlobIdentifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockRoot []byte `protobuf:"bytes,1,opt,name=block_root,json=blockRoot,proto3" json:"block_root,omitempty" ssz-size:"32"`
	Index     uint64 `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *BlobIdentifier) Reset() {
	*x = BlobIdentifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eth_v2_blobs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlobIdentifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlobIdentifier) ProtoMessage() {}

func (x *BlobIdentifier) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eth_v2_blobs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlobIdentifier.ProtoReflect.Descriptor instead.
func (*BlobIdentifier) Descriptor() ([]byte, []int) {
	return file_proto_eth_v2_blobs_proto_rawDescGZIP(), []int{3}
}

func (x *BlobIdentifier) GetBlockRoot() []byte {
	if x != nil {
		return x.BlockRoot
	}
	return nil
}

func (x *BlobIdentifier) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

type SignedBlindedBlobSidecar struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message   *BlindedBlobSidecar `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Signature []byte              `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty" ssz-size:"96"`
}

func (x *SignedBlindedBlobSidecar) Reset() {
	*x = SignedBlindedBlobSidecar{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eth_v2_blobs_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedBlindedBlobSidecar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedBlindedBlobSidecar) ProtoMessage() {}

func (x *SignedBlindedBlobSidecar) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eth_v2_blobs_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedBlindedBlobSidecar.ProtoReflect.Descriptor instead.
func (*SignedBlindedBlobSidecar) Descriptor() ([]byte, []int) {
	return file_proto_eth_v2_blobs_proto_rawDescGZIP(), []int{4}
}

func (x *SignedBlindedBlobSidecar) GetMessage() *BlindedBlobSidecar {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *SignedBlindedBlobSidecar) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type BlindedBlobSidecar struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockRoot       []byte                                                                      `protobuf:"bytes,1,opt,name=block_root,json=blockRoot,proto3" json:"block_root,omitempty" ssz-size:"32"`
	Index           uint64                                                                      `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Slot            github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.Slot           `protobuf:"varint,3,opt,name=slot,proto3" json:"slot,omitempty" cast-type:"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives.Slot"`
	BlockParentRoot []byte                                                                      `protobuf:"bytes,4,opt,name=block_parent_root,json=blockParentRoot,proto3" json:"block_parent_root,omitempty" ssz-size:"32"`
	ProposerIndex   github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.ValidatorIndex `protobuf:"varint,5,opt,name=proposer_index,json=proposerIndex,proto3" json:"proposer_index,omitempty" cast-type:"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives.ValidatorIndex"`
	BlobRoot        []byte                                                                      `protobuf:"bytes,6,opt,name=blob_root,json=blobRoot,proto3" json:"blob_root,omitempty" ssz-size:"32"`
	KzgCommitment   []byte                                                                      `protobuf:"bytes,7,opt,name=kzg_commitment,json=kzgCommitment,proto3" json:"kzg_commitment,omitempty" ssz-size:"48"`
	KzgProof        []byte                                                                      `protobuf:"bytes,8,opt,name=kzg_proof,json=kzgProof,proto3" json:"kzg_proof,omitempty" ssz-size:"48"`
}

func (x *BlindedBlobSidecar) Reset() {
	*x = BlindedBlobSidecar{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eth_v2_blobs_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlindedBlobSidecar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlindedBlobSidecar) ProtoMessage() {}

func (x *BlindedBlobSidecar) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eth_v2_blobs_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlindedBlobSidecar.ProtoReflect.Descriptor instead.
func (*BlindedBlobSidecar) Descriptor() ([]byte, []int) {
	return file_proto_eth_v2_blobs_proto_rawDescGZIP(), []int{5}
}

func (x *BlindedBlobSidecar) GetBlockRoot() []byte {
	if x != nil {
		return x.BlockRoot
	}
	return nil
}

func (x *BlindedBlobSidecar) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *BlindedBlobSidecar) GetSlot() github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.Slot {
	if x != nil {
		return x.Slot
	}
	return github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.Slot(0)
}

func (x *BlindedBlobSidecar) GetBlockParentRoot() []byte {
	if x != nil {
		return x.BlockParentRoot
	}
	return nil
}

func (x *BlindedBlobSidecar) GetProposerIndex() github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.ValidatorIndex {
	if x != nil {
		return x.ProposerIndex
	}
	return github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.ValidatorIndex(0)
}

func (x *BlindedBlobSidecar) GetBlobRoot() []byte {
	if x != nil {
		return x.BlobRoot
	}
	return nil
}

func (x *BlindedBlobSidecar) GetKzgCommitment() []byte {
	if x != nil {
		return x.KzgCommitment
	}
	return nil
}

func (x *BlindedBlobSidecar) GetKzgProof() []byte {
	if x != nil {
		return x.KzgProof
	}
	return nil
}

var File_proto_eth_v2_blobs_proto protoreflect.FileDescriptor

var file_proto_eth_v2_blobs_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x32, 0x2f, 0x62,
	0x6c, 0x6f, 0x62, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x65, 0x74, 0x68, 0x65,
	0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x32, 0x1a, 0x1b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x65, 0x78, 0x74, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a, 0x0c, 0x42, 0x6c, 0x6f, 0x62,
	0x53, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x73, 0x12, 0x38, 0x0a, 0x08, 0x73, 0x69, 0x64, 0x65,
	0x63, 0x61, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x65, 0x74, 0x68,
	0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x32, 0x2e, 0x42, 0x6c, 0x6f,
	0x62, 0x53, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x52, 0x08, 0x73, 0x69, 0x64, 0x65, 0x63, 0x61,
	0x72, 0x73, 0x22, 0xc5, 0x03, 0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x62, 0x53, 0x69, 0x64, 0x65, 0x63,
	0x61, 0x72, 0x12, 0x25, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x72, 0x6f, 0x6f, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x09,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x59, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x45, 0x82,
	0xb5, 0x18, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72,
	0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x79, 0x73,
	0x6d, 0x2f, 0x76, 0x34, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2d, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x73, 0x2e,
	0x53, 0x6c, 0x6f, 0x74, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x12, 0x32, 0x0a, 0x11, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x0f, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x76,
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x42, 0x4f, 0x82, 0xb5, 0x18, 0x4b, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63,
	0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x76, 0x34, 0x2f, 0x63, 0x6f,
	0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2d, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x72,
	0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1e, 0x0a, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0c, 0x42, 0x0a, 0x8a, 0xb5, 0x18, 0x06, 0x31, 0x33, 0x31, 0x30, 0x37, 0x32,
	0x52, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x12, 0x2d, 0x0a, 0x0e, 0x6b, 0x7a, 0x67, 0x5f, 0x63, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06,
	0x8a, 0xb5, 0x18, 0x02, 0x34, 0x38, 0x52, 0x0d, 0x6b, 0x7a, 0x67, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x09, 0x6b, 0x7a, 0x67, 0x5f, 0x70, 0x72, 0x6f,
	0x6f, 0x66, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x34, 0x38,
	0x52, 0x08, 0x6b, 0x7a, 0x67, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x22, 0x71, 0x0a, 0x11, 0x53, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x62, 0x53, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x12,
	0x36, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e,
	0x76, 0x32, 0x2e, 0x42, 0x6c, 0x6f, 0x62, 0x53, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02,
	0x39, 0x36, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x4d, 0x0a,
	0x0e, 0x42, 0x6c, 0x6f, 0x62, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12,
	0x25, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52, 0x09, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x7f, 0x0a, 0x18,
	0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x42, 0x6c, 0x69, 0x6e, 0x64, 0x65, 0x64, 0x42, 0x6c, 0x6f,
	0x62, 0x53, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x12, 0x3d, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x65, 0x74, 0x68, 0x65,
	0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x32, 0x2e, 0x42, 0x6c, 0x69, 0x6e,
	0x64, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x62, 0x53, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02,
	0x39, 0x36, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0xd1, 0x03,
	0x0a, 0x12, 0x42, 0x6c, 0x69, 0x6e, 0x64, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x62, 0x53, 0x69, 0x64,
	0x65, 0x63, 0x61, 0x72, 0x12, 0x25, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x72, 0x6f,
	0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32,
	0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x12, 0x59, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42,
	0x45, 0x82, 0xb5, 0x18, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72,
	0x79, 0x73, 0x6d, 0x2f, 0x76, 0x34, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73,
	0x2d, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65,
	0x73, 0x2e, 0x53, 0x6c, 0x6f, 0x74, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x12, 0x32, 0x0a, 0x11,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x6f, 0x6f,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x33, 0x32, 0x52,
	0x0f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6f, 0x74,
	0x12, 0x76, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x42, 0x4f, 0x82, 0xb5, 0x18, 0x4b, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74,
	0x69, 0x63, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x76, 0x34, 0x2f,
	0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x75, 0x73, 0x2d, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f,
	0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x23, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x62,
	0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18,
	0x02, 0x33, 0x32, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x62, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x2d, 0x0a,
	0x0e, 0x6b, 0x7a, 0x67, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x06, 0x8a, 0xb5, 0x18, 0x02, 0x34, 0x38, 0x52, 0x0d, 0x6b,
	0x7a, 0x67, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x09,
	0x6b, 0x7a, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x42,
	0x06, 0x8a, 0xb5, 0x18, 0x02, 0x34, 0x38, 0x52, 0x08, 0x6b, 0x7a, 0x67, 0x50, 0x72, 0x6f, 0x6f,
	0x66, 0x42, 0x7b, 0x0a, 0x13, 0x6f, 0x72, 0x67, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75,
	0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x32, 0x42, 0x0a, 0x42, 0x6c, 0x6f, 0x62, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6c, 0x61, 0x62, 0x73,
	0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x76, 0x34, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x65, 0x74, 0x68, 0x2f, 0x76, 0x32, 0x3b, 0x65, 0x74, 0x68, 0xaa, 0x02, 0x0f, 0x45, 0x74, 0x68,
	0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x45, 0x74, 0x68, 0x2e, 0x56, 0x32, 0xca, 0x02, 0x0f, 0x45,
	0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x5c, 0x45, 0x74, 0x68, 0x5c, 0x76, 0x32, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_eth_v2_blobs_proto_rawDescOnce sync.Once
	file_proto_eth_v2_blobs_proto_rawDescData = file_proto_eth_v2_blobs_proto_rawDesc
)

func file_proto_eth_v2_blobs_proto_rawDescGZIP() []byte {
	file_proto_eth_v2_blobs_proto_rawDescOnce.Do(func() {
		file_proto_eth_v2_blobs_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_eth_v2_blobs_proto_rawDescData)
	})
	return file_proto_eth_v2_blobs_proto_rawDescData
}

var file_proto_eth_v2_blobs_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_eth_v2_blobs_proto_goTypes = []interface{}{
	(*BlobSidecars)(nil),             // 0: ethereum.eth.v2.BlobSidecars
	(*BlobSidecar)(nil),              // 1: ethereum.eth.v2.BlobSidecar
	(*SignedBlobSidecar)(nil),        // 2: ethereum.eth.v2.SignedBlobSidecar
	(*BlobIdentifier)(nil),           // 3: ethereum.eth.v2.BlobIdentifier
	(*SignedBlindedBlobSidecar)(nil), // 4: ethereum.eth.v2.SignedBlindedBlobSidecar
	(*BlindedBlobSidecar)(nil),       // 5: ethereum.eth.v2.BlindedBlobSidecar
}
var file_proto_eth_v2_blobs_proto_depIdxs = []int32{
	1, // 0: ethereum.eth.v2.BlobSidecars.sidecars:type_name -> ethereum.eth.v2.BlobSidecar
	1, // 1: ethereum.eth.v2.SignedBlobSidecar.message:type_name -> ethereum.eth.v2.BlobSidecar
	5, // 2: ethereum.eth.v2.SignedBlindedBlobSidecar.message:type_name -> ethereum.eth.v2.BlindedBlobSidecar
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_eth_v2_blobs_proto_init() }
func file_proto_eth_v2_blobs_proto_init() {
	if File_proto_eth_v2_blobs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_eth_v2_blobs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlobSidecars); i {
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
		file_proto_eth_v2_blobs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlobSidecar); i {
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
		file_proto_eth_v2_blobs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedBlobSidecar); i {
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
		file_proto_eth_v2_blobs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlobIdentifier); i {
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
		file_proto_eth_v2_blobs_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedBlindedBlobSidecar); i {
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
		file_proto_eth_v2_blobs_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlindedBlobSidecar); i {
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
			RawDescriptor: file_proto_eth_v2_blobs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_eth_v2_blobs_proto_goTypes,
		DependencyIndexes: file_proto_eth_v2_blobs_proto_depIdxs,
		MessageInfos:      file_proto_eth_v2_blobs_proto_msgTypes,
	}.Build()
	File_proto_eth_v2_blobs_proto = out.File
	file_proto_eth_v2_blobs_proto_rawDesc = nil
	file_proto_eth_v2_blobs_proto_goTypes = nil
	file_proto_eth_v2_blobs_proto_depIdxs = nil
}
