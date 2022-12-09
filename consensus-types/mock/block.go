package mock

import (
	ssz "github.com/prysmaticlabs/fastssz"
	field_params "github.com/prysmaticlabs/prysm/v3/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/interfaces"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1/validator-client"
	"google.golang.org/protobuf/proto"
)

type SignedBeaconBlock struct {
	BeaconBlock interfaces.BeaconBlock
}

func (SignedBeaconBlock) PbGenericBlock() (*eth.GenericSignedBeaconBlock, error) {
	panic("implement me")
}

func (m SignedBeaconBlock) Block() interfaces.BeaconBlock {
	return m.BeaconBlock
}

func (SignedBeaconBlock) Signature() [field_params.BLSSignatureLength]byte {
	panic("implement me")
}

func (m SignedBeaconBlock) IsNil() bool {
	return m.BeaconBlock == nil || m.Block().IsNil()
}

func (SignedBeaconBlock) Copy() (interfaces.SignedBeaconBlock, error) {
	panic("implement me")
}

func (SignedBeaconBlock) Proto() (proto.Message, error) {
	panic("implement me")
}

func (SignedBeaconBlock) PbPhase0Block() (*eth.SignedBeaconBlock, error) {
	panic("implement me")
}

func (SignedBeaconBlock) PbAltairBlock() (*eth.SignedBeaconBlockAltair, error) {
	panic("implement me")
}

func (SignedBeaconBlock) PbBellatrixBlock() (*eth.SignedBeaconBlockBellatrix, error) {
	panic("implement me")
}

func (SignedBeaconBlock) PbBlindedBellatrixBlock() (*eth.SignedBlindedBeaconBlockBellatrix, error) {
	panic("implement me")
}

func (SignedBeaconBlock) PbCapellaBlock() (*eth.SignedBeaconBlockCapella, error) {
	panic("implement me")
}

func (SignedBeaconBlock) PbBlindedCapellaBlock() (*eth.SignedBlindedBeaconBlockCapella, error) {
	panic("implement me")
}

func (SignedBeaconBlock) MarshalSSZTo(_ []byte) ([]byte, error) {
	panic("implement me")
}

func (SignedBeaconBlock) MarshalSSZ() ([]byte, error) {
	panic("implement me")
}

func (SignedBeaconBlock) SizeSSZ() int {
	panic("implement me")
}

func (SignedBeaconBlock) UnmarshalSSZ(_ []byte) error {
	panic("implement me")
}

func (SignedBeaconBlock) Version() int {
	panic("implement me")
}

func (SignedBeaconBlock) IsBlinded() bool {
	return false
}

func (SignedBeaconBlock) ToBlinded() (interfaces.SignedBeaconBlock, error) {
	panic("implement me")
}

func (SignedBeaconBlock) Header() (*eth.SignedBeaconBlockHeader, error) {
	panic("implement me")
}

func (SignedBeaconBlock) SetBlock(_ interfaces.BeaconBlock) error {
	panic("implement me")
}

func (SignedBeaconBlock) SetSignature(_ []byte) {
	panic("implement me")
}

type BeaconBlock struct {
	Htr             [field_params.RootLength]byte
	HtrErr          error
	BeaconBlockBody interfaces.BeaconBlockBody
	BlockSlot       types.Slot
}

func (BeaconBlock) AsSignRequestObject() (validatorpb.SignRequestObject, error) {
	panic("implement me")
}

func (m BeaconBlock) HashTreeRoot() ([field_params.RootLength]byte, error) {
	return m.Htr, m.HtrErr
}

func (m BeaconBlock) Slot() types.Slot {
	return m.BlockSlot
}

func (BeaconBlock) ProposerIndex() types.ValidatorIndex {
	panic("implement me")
}

func (BeaconBlock) ParentRoot() [field_params.RootLength]byte {
	panic("implement me")
}

func (BeaconBlock) StateRoot() [field_params.RootLength]byte {
	panic("implement me")
}

func (m BeaconBlock) Body() interfaces.BeaconBlockBody {
	return m.BeaconBlockBody
}

func (BeaconBlock) IsNil() bool {
	return false
}

func (BeaconBlock) IsBlinded() bool {
	return false
}

func (BeaconBlock) Proto() (proto.Message, error) {
	panic("implement me")
}

func (BeaconBlock) MarshalSSZTo(_ []byte) ([]byte, error) {
	panic("implement me")
}

func (BeaconBlock) MarshalSSZ() ([]byte, error) {
	panic("implement me")
}

func (BeaconBlock) SizeSSZ() int {
	panic("implement me")
}

func (BeaconBlock) UnmarshalSSZ(_ []byte) error {
	panic("implement me")
}

func (BeaconBlock) HashTreeRootWith(_ *ssz.Hasher) error {
	panic("implement me")
}

func (BeaconBlock) Version() int {
	panic("implement me")
}

func (BeaconBlock) ToBlinded() (interfaces.BeaconBlock, error) {
	panic("implement me")
}

func (BeaconBlock) SetSlot(_ types.Slot) {
	panic("implement me")
}

func (BeaconBlock) SetProposerIndex(_ types.ValidatorIndex) {
	panic("implement me")
}

func (BeaconBlock) SetParentRoot(_ []byte) {
	panic("implement me")
}

func (BeaconBlock) SetBlinded(_ bool) {
	panic("implement me")
}

func (BeaconBlock) Copy() (interfaces.BeaconBlock, error) {
	panic("implement me")
}

type BeaconBlockBody struct{}

func (BeaconBlockBody) RandaoReveal() [field_params.BLSSignatureLength]byte {
	panic("implement me")
}

func (BeaconBlockBody) Eth1Data() *eth.Eth1Data {
	panic("implement me")
}

func (BeaconBlockBody) Graffiti() [field_params.RootLength]byte {
	panic("implement me")
}

func (BeaconBlockBody) ProposerSlashings() []*eth.ProposerSlashing {
	panic("implement me")
}

func (BeaconBlockBody) AttesterSlashings() []*eth.AttesterSlashing {
	panic("implement me")
}

func (BeaconBlockBody) Attestations() []*eth.Attestation {
	panic("implement me")
}

func (BeaconBlockBody) Deposits() []*eth.Deposit {
	panic("implement me")
}

func (BeaconBlockBody) VoluntaryExits() []*eth.SignedVoluntaryExit {
	panic("implement me")
}

func (BeaconBlockBody) SyncAggregate() (*eth.SyncAggregate, error) {
	panic("implement me")
}

func (BeaconBlockBody) IsNil() bool {
	return false
}

func (BeaconBlockBody) HashTreeRoot() ([field_params.RootLength]byte, error) {
	panic("implement me")
}

func (BeaconBlockBody) Proto() (proto.Message, error) {
	panic("implement me")
}

func (BeaconBlockBody) Execution() (interfaces.ExecutionData, error) {
	panic("implement me")
}

func (BeaconBlockBody) BLSToExecutionChanges() ([]*eth.SignedBLSToExecutionChange, error) {
	panic("implement me")
}

func (b *BeaconBlock) SetStateRoot(root []byte) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetRandaoReveal([]byte) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetEth1Data(*eth.Eth1Data) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetGraffiti([]byte) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetProposerSlashings([]*eth.ProposerSlashing) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetAttesterSlashings([]*eth.AttesterSlashing) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetAttestations([]*eth.Attestation) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetDeposits([]*eth.Deposit) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetVoluntaryExits([]*eth.SignedVoluntaryExit) {
	panic("implement me")
}

func (b *BeaconBlockBody) SetSyncAggregate(*eth.SyncAggregate) error {
	panic("implement me")
}

func (b *BeaconBlockBody) SetExecution(interfaces.ExecutionData) error {
	panic("implement me")
}

func (b *BeaconBlockBody) SetBLSToExecutionChanges([]*eth.SignedBLSToExecutionChange) error {
	panic("implement me")
}

var _ interfaces.SignedBeaconBlock = &SignedBeaconBlock{}
var _ interfaces.BeaconBlock = &BeaconBlock{}
var _ interfaces.BeaconBlockBody = &BeaconBlockBody{}
