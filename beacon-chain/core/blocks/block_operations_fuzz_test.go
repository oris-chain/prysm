package blocks

import (
	"context"
	"testing"

	fuzz "github.com/google/gofuzz"
	types "github.com/prysmaticlabs/eth2-types"
	v "github.com/prysmaticlabs/prysm/beacon-chain/core/validators"
	v1 "github.com/prysmaticlabs/prysm/beacon-chain/state/v1"
	"github.com/prysmaticlabs/prysm/config/params"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/wrapper"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func TestFuzzProcessAttestationNoVerify_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	ctx := context.Background()
	state := &ethpb.BeaconState{}
	att := &eth.Attestation{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(att)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		_, err = ProcessAttestationNoVerifySignature(ctx, s, att)
		_ = err
	}
}

func TestFuzzProcessBlockHeader_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	block := &eth.SignedBeaconBlock{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(block)

		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		_, err = ProcessBlockHeader(context.Background(), s, wrapper.WrappedPhase0SignedBeaconBlock(block))
		_ = err
	}
}

func TestFuzzverifyDepositDataSigningRoot_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	var ba []byte
	pubkey := [48]byte{}
	sig := [96]byte{}
	domain := [4]byte{}
	var p []byte
	var s []byte
	var d []byte
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(&ba)
		fuzzer.Fuzz(&pubkey)
		fuzzer.Fuzz(&sig)
		fuzzer.Fuzz(&domain)
		fuzzer.Fuzz(&p)
		fuzzer.Fuzz(&s)
		fuzzer.Fuzz(&d)
		err := verifySignature(ba, pubkey[:], sig[:], domain[:])
		_ = err
		err = verifySignature(ba, p, s, d)
		_ = err
	}
}

func TestFuzzProcessEth1DataInBlock_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	e := &eth.Eth1Data{}
	state := &v1.BeaconState{}
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(e)
		s, err := ProcessEth1DataInBlock(context.Background(), state, e)
		if err != nil && s != nil {
			t.Fatalf("state should be nil on err. found: %v on error: %v for state: %v and eth1data: %v", s, err, state, e)
		}
	}
}

func TestFuzzareEth1DataEqual_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	eth1data := &eth.Eth1Data{}
	eth1data2 := &eth.Eth1Data{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(eth1data)
		fuzzer.Fuzz(eth1data2)
		AreEth1DataEqual(eth1data, eth1data2)
		AreEth1DataEqual(eth1data, eth1data)
	}
}

func TestFuzzEth1DataHasEnoughSupport_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	eth1data := &eth.Eth1Data{}
	var stateVotes []*eth.Eth1Data
	for i := 0; i < 100000; i++ {
		fuzzer.Fuzz(eth1data)
		fuzzer.Fuzz(&stateVotes)
		s, err := v1.InitializeFromProto(&ethpb.BeaconState{
			Eth1DataVotes: stateVotes,
		})
		require.NoError(t, err)
		_, err = Eth1DataHasEnoughSupport(s, eth1data)
		_ = err
	}

}

func TestFuzzProcessBlockHeaderNoVerify_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	block := &eth.BeaconBlock{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(block)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		_, err = ProcessBlockHeaderNoVerify(context.Background(), s, block.Slot, block.ProposerIndex, block.ParentRoot, []byte{})
		_ = err
	}
}

func TestFuzzProcessRandao_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	b := &eth.SignedBeaconBlock{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(b)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessRandao(context.Background(), s, wrapper.WrappedPhase0SignedBeaconBlock(b))
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and block: %v", r, err, state, b)
		}
	}
}

func TestFuzzProcessRandaoNoVerify_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	blockBody := &eth.BeaconBlockBody{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(blockBody)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessRandaoNoVerify(s, blockBody.RandaoReveal)
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and block: %v", r, err, state, blockBody)
		}
	}
}

func TestFuzzProcessProposerSlashings_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	p := &eth.ProposerSlashing{}
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(p)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessProposerSlashings(ctx, s, []*eth.ProposerSlashing{p}, v.SlashValidator)
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and slashing: %v", r, err, state, p)
		}
	}
}

func TestFuzzVerifyProposerSlashing_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	proposerSlashing := &eth.ProposerSlashing{}
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(proposerSlashing)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		err = VerifyProposerSlashing(s, proposerSlashing)
		_ = err
	}
}

func TestFuzzProcessAttesterSlashings_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	a := &eth.AttesterSlashing{}
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(a)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessAttesterSlashings(ctx, s, []*eth.AttesterSlashing{a}, v.SlashValidator)
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and slashing: %v", r, err, state, a)
		}
	}
}

func TestFuzzVerifyAttesterSlashing_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	attesterSlashing := &eth.AttesterSlashing{}
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(attesterSlashing)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		err = VerifyAttesterSlashing(ctx, s, attesterSlashing)
		_ = err
	}
}

func TestFuzzIsSlashableAttestationData_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	attestationData := &eth.AttestationData{}
	attestationData2 := &eth.AttestationData{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(attestationData)
		fuzzer.Fuzz(attestationData2)
		IsSlashableAttestationData(attestationData, attestationData2)
	}
}

func TestFuzzslashableAttesterIndices_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	attesterSlashing := &eth.AttesterSlashing{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(attesterSlashing)
		slashableAttesterIndices(attesterSlashing)
	}
}

func TestFuzzProcessAttestationsNoVerify_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	b := &eth.SignedBeaconBlock{}
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(b)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessAttestationsNoVerifySignature(ctx, s, wrapper.WrappedPhase0SignedBeaconBlock(b))
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and block: %v", r, err, state, b)
		}
	}
}

func TestFuzzVerifyIndexedAttestationn_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	idxAttestation := &eth.IndexedAttestation{}
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(idxAttestation)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		err = VerifyIndexedAttestation(ctx, s, idxAttestation)
		_ = err
	}
}

func TestFuzzVerifyAttestation_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	attestation := &eth.Attestation{}
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(attestation)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		err = VerifyAttestationSignature(ctx, s, attestation)
		_ = err
	}
}

func TestFuzzProcessDeposits_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	deposits := make([]*eth.Deposit, 100)
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		for i := range deposits {
			fuzzer.Fuzz(deposits[i])
		}
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessDeposits(ctx, s, deposits)
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and block: %v", r, err, state, deposits)
		}
	}
}

func TestFuzzProcessPreGenesisDeposit_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	deposit := &eth.Deposit{}
	ctx := context.Background()

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(deposit)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessPreGenesisDeposits(ctx, s, []*eth.Deposit{deposit})
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and block: %v", r, err, state, deposit)
		}
	}
}

func TestFuzzProcessDeposit_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	deposit := &eth.Deposit{}

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(deposit)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, _, err := ProcessDeposit(s, deposit, true)
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and block: %v", r, err, state, deposit)
		}
	}
}

func TestFuzzverifyDeposit_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	deposit := &eth.Deposit{}
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(deposit)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		err = verifyDeposit(s, deposit)
		_ = err
	}
}

func TestFuzzProcessVoluntaryExits_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	e := &eth.SignedVoluntaryExit{}
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(e)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessVoluntaryExits(ctx, s, []*eth.SignedVoluntaryExit{e})
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and exit: %v", r, err, state, e)
		}
	}
}

func TestFuzzProcessVoluntaryExitsNoVerify_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	state := &ethpb.BeaconState{}
	e := &eth.SignedVoluntaryExit{}
	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(state)
		fuzzer.Fuzz(e)
		s, err := v1.InitializeFromProtoUnsafe(state)
		require.NoError(t, err)
		r, err := ProcessVoluntaryExits(context.Background(), s, []*eth.SignedVoluntaryExit{e})
		if err != nil && r != nil {
			t.Fatalf("return value should be nil on err. found: %v on error: %v for state: %v and block: %v", r, err, state, e)
		}
	}
}

func TestFuzzVerifyExit_10000(t *testing.T) {
	fuzzer := fuzz.NewWithSeed(0)
	ve := &eth.SignedVoluntaryExit{}
	rawVal := &ethpb.Validator{}
	fork := &ethpb.Fork{}
	var slot types.Slot

	for i := 0; i < 10000; i++ {
		fuzzer.Fuzz(ve)
		fuzzer.Fuzz(rawVal)
		fuzzer.Fuzz(fork)
		fuzzer.Fuzz(&slot)
		val, err := v1.NewValidator(&ethpb.Validator{})
		_ = err
		err = VerifyExitAndSignature(val, slot, fork, ve, params.BeaconConfig().ZeroHash[:])
		_ = err
	}
}
