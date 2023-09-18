// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.3
// source: proto/eth/service/validator_service.proto

package service

import (
	context "context"
	reflect "reflect"

	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/golang/protobuf/ptypes/empty"
	v1 "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	v2 "github.com/prysmaticlabs/prysm/v4/proto/eth/v2"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_proto_eth_service_validator_service_proto protoreflect.FileDescriptor

var file_proto_eth_service_validator_service_proto_rawDesc = []byte{
	0x0a, 0x29, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x65, 0x74, 0x68,
	0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x32, 0x2f, 0x73, 0x73, 0x7a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f,
	0x76, 0x32, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x32, 0x8d, 0x06, 0x0a, 0x0f, 0x42, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x91, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x56, 0x32, 0x12, 0x24, 0x2e, 0x65, 0x74, 0x68, 0x65,
	0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x27, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76,
	0x32, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x32, 0x22, 0x30, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2a,
	0x12, 0x28, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x74, 0x68, 0x2f,
	0x76, 0x32, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x73, 0x2f, 0x7b, 0x73, 0x6c, 0x6f, 0x74, 0x7d, 0x12, 0x8e, 0x01, 0x0a, 0x11, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x56, 0x32, 0x53, 0x53, 0x5a,
	0x12, 0x24, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75,
	0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x32, 0x2e, 0x53, 0x53, 0x5a, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x22, 0x34, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2e, 0x12, 0x2c, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x32, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73,
	0x2f, 0x7b, 0x73, 0x6c, 0x6f, 0x74, 0x7d, 0x2f, 0x73, 0x73, 0x7a, 0x12, 0xa3, 0x01, 0x0a, 0x13,
	0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x42, 0x6c, 0x69, 0x6e, 0x64, 0x65, 0x64, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x12, 0x24, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65,
	0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x65, 0x74, 0x68, 0x65,
	0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x32, 0x2e, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x65, 0x42, 0x6c, 0x69, 0x6e, 0x64, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x38, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x32, 0x12,
	0x30, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76,
	0x31, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x62, 0x6c, 0x69, 0x6e,
	0x64, 0x65, 0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x2f, 0x7b, 0x73, 0x6c, 0x6f, 0x74,
	0x7d, 0x12, 0x9b, 0x01, 0x0a, 0x16, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x42, 0x6c, 0x69,
	0x6e, 0x64, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x53, 0x5a, 0x12, 0x24, 0x2e, 0x65,
	0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74,
	0x68, 0x2e, 0x76, 0x32, 0x2e, 0x53, 0x53, 0x5a, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x22, 0x3c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x36, 0x12, 0x34, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x62, 0x6c, 0x69, 0x6e, 0x64, 0x65, 0x64, 0x5f, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x73, 0x2f, 0x7b, 0x73, 0x6c, 0x6f, 0x74, 0x7d, 0x2f, 0x73, 0x73, 0x7a, 0x12,
	0x90, 0x01, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x73, 0x73, 0x12,
	0x23, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76,
	0x32, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e,
	0x65, 0x74, 0x68, 0x2e, 0x76, 0x32, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x36, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x30, 0x3a, 0x01, 0x2a, 0x22, 0x2b, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x73, 0x73, 0x2f, 0x7b, 0x65, 0x70, 0x6f, 0x63,
	0x68, 0x7d, 0x42, 0x96, 0x01, 0x0a, 0x18, 0x6f, 0x72, 0x67, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72,
	0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42,
	0x15, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6c, 0x61,
	0x62, 0x73, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x76, 0x34, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0xaa, 0x02, 0x14,
	0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x45, 0x74, 0x68, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0xca, 0x02, 0x14, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x5c,
	0x45, 0x74, 0x68, 0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var file_proto_eth_service_validator_service_proto_goTypes = []interface{}{
	(*v1.ProduceBlockRequest)(nil),         // 0: ethereum.eth.v1.ProduceBlockRequest
	(*v2.GetLivenessRequest)(nil),          // 1: ethereum.eth.v2.GetLivenessRequest
	(*v2.ProduceBlockResponseV2)(nil),      // 2: ethereum.eth.v2.ProduceBlockResponseV2
	(*v2.SSZContainer)(nil),                // 3: ethereum.eth.v2.SSZContainer
	(*v2.ProduceBlindedBlockResponse)(nil), // 4: ethereum.eth.v2.ProduceBlindedBlockResponse
	(*v2.GetLivenessResponse)(nil),         // 5: ethereum.eth.v2.GetLivenessResponse
}
var file_proto_eth_service_validator_service_proto_depIdxs = []int32{
	0, // 0: ethereum.eth.service.BeaconValidator.ProduceBlockV2:input_type -> ethereum.eth.v1.ProduceBlockRequest
	0, // 1: ethereum.eth.service.BeaconValidator.ProduceBlockV2SSZ:input_type -> ethereum.eth.v1.ProduceBlockRequest
	0, // 2: ethereum.eth.service.BeaconValidator.ProduceBlindedBlock:input_type -> ethereum.eth.v1.ProduceBlockRequest
	0, // 3: ethereum.eth.service.BeaconValidator.ProduceBlindedBlockSSZ:input_type -> ethereum.eth.v1.ProduceBlockRequest
	1, // 4: ethereum.eth.service.BeaconValidator.GetLiveness:input_type -> ethereum.eth.v2.GetLivenessRequest
	2, // 5: ethereum.eth.service.BeaconValidator.ProduceBlockV2:output_type -> ethereum.eth.v2.ProduceBlockResponseV2
	3, // 6: ethereum.eth.service.BeaconValidator.ProduceBlockV2SSZ:output_type -> ethereum.eth.v2.SSZContainer
	4, // 7: ethereum.eth.service.BeaconValidator.ProduceBlindedBlock:output_type -> ethereum.eth.v2.ProduceBlindedBlockResponse
	3, // 8: ethereum.eth.service.BeaconValidator.ProduceBlindedBlockSSZ:output_type -> ethereum.eth.v2.SSZContainer
	5, // 9: ethereum.eth.service.BeaconValidator.GetLiveness:output_type -> ethereum.eth.v2.GetLivenessResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_eth_service_validator_service_proto_init() }
func file_proto_eth_service_validator_service_proto_init() {
	if File_proto_eth_service_validator_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_eth_service_validator_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_eth_service_validator_service_proto_goTypes,
		DependencyIndexes: file_proto_eth_service_validator_service_proto_depIdxs,
	}.Build()
	File_proto_eth_service_validator_service_proto = out.File
	file_proto_eth_service_validator_service_proto_rawDesc = nil
	file_proto_eth_service_validator_service_proto_goTypes = nil
	file_proto_eth_service_validator_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BeaconValidatorClient is the client API for BeaconValidator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BeaconValidatorClient interface {
	ProduceBlockV2(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.ProduceBlockResponseV2, error)
	ProduceBlockV2SSZ(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.SSZContainer, error)
	ProduceBlindedBlock(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.ProduceBlindedBlockResponse, error)
	ProduceBlindedBlockSSZ(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.SSZContainer, error)
	GetLiveness(ctx context.Context, in *v2.GetLivenessRequest, opts ...grpc.CallOption) (*v2.GetLivenessResponse, error)
}

type beaconValidatorClient struct {
	cc grpc.ClientConnInterface
}

func NewBeaconValidatorClient(cc grpc.ClientConnInterface) BeaconValidatorClient {
	return &beaconValidatorClient{cc}
}

func (c *beaconValidatorClient) ProduceBlockV2(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.ProduceBlockResponseV2, error) {
	out := new(v2.ProduceBlockResponseV2)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconValidator/ProduceBlockV2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconValidatorClient) ProduceBlockV2SSZ(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.SSZContainer, error) {
	out := new(v2.SSZContainer)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconValidator/ProduceBlockV2SSZ", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconValidatorClient) ProduceBlindedBlock(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.ProduceBlindedBlockResponse, error) {
	out := new(v2.ProduceBlindedBlockResponse)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconValidator/ProduceBlindedBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconValidatorClient) ProduceBlindedBlockSSZ(ctx context.Context, in *v1.ProduceBlockRequest, opts ...grpc.CallOption) (*v2.SSZContainer, error) {
	out := new(v2.SSZContainer)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconValidator/ProduceBlindedBlockSSZ", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconValidatorClient) GetLiveness(ctx context.Context, in *v2.GetLivenessRequest, opts ...grpc.CallOption) (*v2.GetLivenessResponse, error) {
	out := new(v2.GetLivenessResponse)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconValidator/GetLiveness", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BeaconValidatorServer is the server API for BeaconValidator service.
type BeaconValidatorServer interface {
	ProduceBlockV2(context.Context, *v1.ProduceBlockRequest) (*v2.ProduceBlockResponseV2, error)
	ProduceBlockV2SSZ(context.Context, *v1.ProduceBlockRequest) (*v2.SSZContainer, error)
	ProduceBlindedBlock(context.Context, *v1.ProduceBlockRequest) (*v2.ProduceBlindedBlockResponse, error)
	ProduceBlindedBlockSSZ(context.Context, *v1.ProduceBlockRequest) (*v2.SSZContainer, error)
	GetLiveness(context.Context, *v2.GetLivenessRequest) (*v2.GetLivenessResponse, error)
}

// UnimplementedBeaconValidatorServer can be embedded to have forward compatible implementations.
type UnimplementedBeaconValidatorServer struct {
}

func (*UnimplementedBeaconValidatorServer) ProduceBlockV2(context.Context, *v1.ProduceBlockRequest) (*v2.ProduceBlockResponseV2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceBlockV2 not implemented")
}
func (*UnimplementedBeaconValidatorServer) ProduceBlockV2SSZ(context.Context, *v1.ProduceBlockRequest) (*v2.SSZContainer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceBlockV2SSZ not implemented")
}
func (*UnimplementedBeaconValidatorServer) ProduceBlindedBlock(context.Context, *v1.ProduceBlockRequest) (*v2.ProduceBlindedBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceBlindedBlock not implemented")
}
func (*UnimplementedBeaconValidatorServer) ProduceBlindedBlockSSZ(context.Context, *v1.ProduceBlockRequest) (*v2.SSZContainer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceBlindedBlockSSZ not implemented")
}
func (*UnimplementedBeaconValidatorServer) GetLiveness(context.Context, *v2.GetLivenessRequest) (*v2.GetLivenessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLiveness not implemented")
}

func RegisterBeaconValidatorServer(s *grpc.Server, srv BeaconValidatorServer) {
	s.RegisterService(&_BeaconValidator_serviceDesc, srv)
}

func _BeaconValidator_ProduceBlockV2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ProduceBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconValidatorServer).ProduceBlockV2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconValidator/ProduceBlockV2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconValidatorServer).ProduceBlockV2(ctx, req.(*v1.ProduceBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeaconValidator_ProduceBlockV2SSZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ProduceBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconValidatorServer).ProduceBlockV2SSZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconValidator/ProduceBlockV2SSZ",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconValidatorServer).ProduceBlockV2SSZ(ctx, req.(*v1.ProduceBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeaconValidator_ProduceBlindedBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ProduceBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconValidatorServer).ProduceBlindedBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconValidator/ProduceBlindedBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconValidatorServer).ProduceBlindedBlock(ctx, req.(*v1.ProduceBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeaconValidator_ProduceBlindedBlockSSZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ProduceBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconValidatorServer).ProduceBlindedBlockSSZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconValidator/ProduceBlindedBlockSSZ",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconValidatorServer).ProduceBlindedBlockSSZ(ctx, req.(*v1.ProduceBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeaconValidator_GetLiveness_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v2.GetLivenessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconValidatorServer).GetLiveness(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconValidator/GetLiveness",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconValidatorServer).GetLiveness(ctx, req.(*v2.GetLivenessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BeaconValidator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.eth.service.BeaconValidator",
	HandlerType: (*BeaconValidatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProduceBlockV2",
			Handler:    _BeaconValidator_ProduceBlockV2_Handler,
		},
		{
			MethodName: "ProduceBlockV2SSZ",
			Handler:    _BeaconValidator_ProduceBlockV2SSZ_Handler,
		},
		{
			MethodName: "ProduceBlindedBlock",
			Handler:    _BeaconValidator_ProduceBlindedBlock_Handler,
		},
		{
			MethodName: "ProduceBlindedBlockSSZ",
			Handler:    _BeaconValidator_ProduceBlindedBlockSSZ_Handler,
		},
		{
			MethodName: "GetLiveness",
			Handler:    _BeaconValidator_GetLiveness_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/eth/service/validator_service.proto",
}
