package validator

import (
	"context"
	"math/big"

	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/execution"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/time"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/transition"
	"github.com/prysmaticlabs/prysm/config/params"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/time/slots"
)

func (vs *Server) getExecutionPayload(ctx context.Context, slot types.Slot) (*ethpb.ExecutionPayload, error) {
	// TODO: Reuse the same head state as in building phase0 block attestation.
	st, err := vs.HeadFetcher.HeadState(ctx)
	if err != nil {
		return nil, err
	}
	st, err = transition.ProcessSlots(ctx, st, slot)
	if err != nil {
		return nil, err
	}

	var parentHash []byte
	var hasTerminalBlock bool
	complete, err := execution.IsMergeComplete(st)
	if err != nil {
		return nil, err
	}
	if !complete {
		parentHash, hasTerminalBlock, err = vs.getPreMergeParentHash(ctx, slot)
		if err != nil {
			return nil, err
		}
		if !hasTerminalBlock {
			return execution.EmptyPayload(), nil
		}
	} else {
		header, err := st.LatestExecutionPayloadHeader()
		if err != nil {
			return nil, err
		}
		parentHash = header.ParentHash
	}

	t, err := slots.ToTime(st.GenesisTime(), slot)
	if err != nil {
		return nil, err
	}
	random, err := helpers.RandaoMix(st, time.CurrentEpoch(st))
	if err != nil {
		return nil, err
	}
	id, err := vs.ExecutionEngineCaller.PreparePayload(ctx, parentHash, uint64(t.Unix()), random, params.BeaconConfig().FeeRecipient.Bytes())
	if err != nil {
		return nil, err
	}
	return vs.ExecutionEngineCaller.GetPayload(ctx, id)
}

func (vs *Server) getPreMergeParentHash(ctx context.Context, slot types.Slot) ([]byte, bool, error) {
	terminalBlockHash := params.BeaconConfig().TerminalBlockHash
	if params.BeaconConfig().TerminalBlockHash != params.BeaconConfig().ZeroHash {
		e, _, err := vs.Eth1BlockFetcher.BlockExists(ctx, terminalBlockHash)
		if err != nil {
			return nil, false, err
		}
		if !e {
			return nil, false, nil
		}
		return terminalBlockHash.Bytes(), true, nil
	}

	return vs.getPowBlockHashAtTerminalTotalDifficulty(ctx)
}

func (vs *Server) getPowBlockHashAtTerminalTotalDifficulty(ctx context.Context) ([]byte, bool, error) {
	b, err := vs.BlockFetcher.BlockByNumber(ctx, nil /* latest block */)
	if err != nil {
		return nil, false, err
	}
	terminalTotalDifficulty := new(big.Int)
	terminalTotalDifficulty.SetBytes(params.BeaconConfig().TerminalTotalDifficulty)
	var terminalBlockHash []byte

	// TODO: Add pow block cache.
	for {
		if b.TotalDifficulty().Cmp(terminalTotalDifficulty) >= 0 {
			terminalBlockHash = b.Hash().Bytes()
			if b.ParentHash() == b.Hash() {
				return nil, false, errors.New("invalid block")
			}
			b, err = vs.BlockFetcher.BlockByHash(ctx, b.ParentHash())
			if err != nil {
				return nil, false, err
			}
		} else {
			return terminalBlockHash, true, err
		}
	}

	return []byte{}, false, nil
}
