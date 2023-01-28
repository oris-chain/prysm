package interfaces

import (
	ssz "github.com/prysmaticlabs/fastssz"
	field_params "github.com/prysmaticlabs/prysm/v3/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1/validator-client"
	"google.golang.org/protobuf/proto"
)

// SignedBeaconBlock is an interface describing the method set of
// a signed beacon block.
type SignedBeaconBlock interface {
	Block() BeaconBlock
	Signature() [field_params.BLSSignatureLength]byte
	SetSignature(sig []byte)
	IsNil() bool
	Copy() (SignedBeaconBlock, error)
	Proto() (proto.Message, error)
	PbGenericBlock() (*ethpb.GenericSignedBeaconBlock, error)
	PbPhase0Block() (*ethpb.SignedBeaconBlock, error)
	PbAltairBlock() (*ethpb.SignedBeaconBlockAltair, error)
	ToBlinded() (SignedBeaconBlock, error)
	PbBellatrixBlock() (*ethpb.SignedBeaconBlockBellatrix, error)
	PbBlindedBellatrixBlock() (*ethpb.SignedBlindedBeaconBlockBellatrix, error)
	PbCapellaBlock() (*ethpb.SignedBeaconBlockCapella, error)
	PbDenebBlock() (*ethpb.SignedBeaconBlockDeneb, error)
	PbBlindedCapellaBlock() (*ethpb.SignedBlindedBeaconBlockCapella, error)
	PbBlindedDenebBlock() (*ethpb.SignedBlindedBeaconBlockDeneb, error)
	ssz.Marshaler
	ssz.Unmarshaler
	Version() int
	IsBlinded() bool
	Header() (*ethpb.SignedBeaconBlockHeader, error)
}

// BeaconBlock describes an interface which states the methods
// employed by an object that is a beacon block.
type BeaconBlock interface {
	Slot() primitives.Slot
	SetSlot(slot primitives.Slot)
	ProposerIndex() primitives.ValidatorIndex
	SetProposerIndex(idx primitives.ValidatorIndex)
	ParentRoot() [field_params.RootLength]byte
	SetParentRoot([]byte)
	StateRoot() [field_params.RootLength]byte
	SetStateRoot([]byte)
	Body() BeaconBlockBody
	IsNil() bool
	IsBlinded() bool
	SetBlinded(bool)
	HashTreeRoot() ([field_params.RootLength]byte, error)
	Proto() (proto.Message, error)
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	Version() int
	AsSignRequestObject() (validatorpb.SignRequestObject, error)
	Copy() (BeaconBlock, error)
}

// BeaconBlockBody describes the method set employed by an object
// that is a beacon block body.
type BeaconBlockBody interface {
	RandaoReveal() [field_params.BLSSignatureLength]byte
	SetRandaoReveal([]byte)
	Eth1Data() *ethpb.Eth1Data
	SetEth1Data(*ethpb.Eth1Data)
	Graffiti() [field_params.RootLength]byte
	SetGraffiti([]byte)
	ProposerSlashings() []*ethpb.ProposerSlashing
	SetProposerSlashings([]*ethpb.ProposerSlashing)
	AttesterSlashings() []*ethpb.AttesterSlashing
	SetAttesterSlashings([]*ethpb.AttesterSlashing)
	Attestations() []*ethpb.Attestation
	SetAttestations([]*ethpb.Attestation)
	Deposits() []*ethpb.Deposit
	SetDeposits([]*ethpb.Deposit)
	VoluntaryExits() []*ethpb.SignedVoluntaryExit
	SetVoluntaryExits([]*ethpb.SignedVoluntaryExit)
	SyncAggregate() (*ethpb.SyncAggregate, error)
	SetSyncAggregate(*ethpb.SyncAggregate) error
	IsNil() bool
	HashTreeRoot() ([field_params.RootLength]byte, error)
	Proto() (proto.Message, error)
	Execution() (ExecutionData, error)
	SetExecution(ExecutionData) error
	BLSToExecutionChanges() ([]*ethpb.SignedBLSToExecutionChange, error)
	SetBLSToExecutionChanges([]*ethpb.SignedBLSToExecutionChange) error
	BlobKzgCommitments() ([][]byte, error)
	SetBlobKzgCommitments(c [][]byte) error
}

// ExecutionData represents execution layer information that is contained
// within post-Bellatrix beacon block bodies.
type ExecutionData interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	Proto() proto.Message
	ParentHash() []byte
	FeeRecipient() []byte
	StateRoot() []byte
	ReceiptsRoot() []byte
	LogsBloom() []byte
	PrevRandao() []byte
	BlockNumber() uint64
	GasLimit() uint64
	GasUsed() uint64
	Timestamp() uint64
	ExtraData() []byte
	BaseFeePerGas() []byte
	ExcessiveDataGas() ([]byte, error)
	BlockHash() []byte
	Transactions() ([][]byte, error)
	TransactionsRoot() ([]byte, error)
	Withdrawals() ([]*enginev1.Withdrawal, error)
	WithdrawalsRoot() ([]byte, error)
	PbCapella() (*enginev1.ExecutionPayloadCapella, error)
	PbBellatrix() (*enginev1.ExecutionPayload, error)
}
