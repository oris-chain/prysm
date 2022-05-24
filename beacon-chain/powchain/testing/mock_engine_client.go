package testing

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/consensus-types/wrapper"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	pb "github.com/prysmaticlabs/prysm/proto/engine/v1"
)

// EngineClient --
type EngineClient struct {
	NewPayloadResp              []byte
	PayloadIDBytes              *pb.PayloadIDBytes
	ForkChoiceUpdatedResp       []byte
	ExecutionPayload            *pb.ExecutionPayload
	ExecutionBlock              *pb.ExecutionBlock
	Err                         error
	ErrLatestExecBlock          error
	ErrExecBlockByHash          error
	ErrForkchoiceUpdated        error
	ErrNewPayload               error
	ExecutionPayloadByBlockHash map[[32]byte]*pb.ExecutionPayload
	BlockByHashMap              map[[32]byte]*pb.ExecutionBlock
	BlockWithTxsByHashMap       map[[32]byte]*pb.ExecutionBlockWithTxs
	NumReconstructedPayloads    uint64
	TerminalBlockHash           []byte
	TerminalBlockHashExists     bool
	OverrideValidHash           [32]byte
}

// NewPayload --
func (e *EngineClient) NewPayload(_ context.Context, _ *pb.ExecutionPayload) ([]byte, error) {
	return e.NewPayloadResp, e.ErrNewPayload
}

// ForkchoiceUpdated --
func (e *EngineClient) ForkchoiceUpdated(
	_ context.Context, fcs *pb.ForkchoiceState, _ *pb.PayloadAttributes,
) (*pb.PayloadIDBytes, []byte, error) {
	if e.OverrideValidHash != [32]byte{} && bytesutil.ToBytes32(fcs.HeadBlockHash) == e.OverrideValidHash {
		return e.PayloadIDBytes, e.ForkChoiceUpdatedResp, nil
	}
	return e.PayloadIDBytes, e.ForkChoiceUpdatedResp, e.ErrForkchoiceUpdated
}

// GetPayload --
func (e *EngineClient) GetPayload(_ context.Context, _ [8]byte) (*pb.ExecutionPayload, error) {
	return e.ExecutionPayload, nil
}

// ExchangeTransitionConfiguration --
func (e *EngineClient) ExchangeTransitionConfiguration(_ context.Context, _ *pb.TransitionConfiguration) error {
	return e.Err
}

// LatestExecutionBlock --
func (e *EngineClient) LatestExecutionBlock(_ context.Context) (*pb.ExecutionBlock, error) {
	return e.ExecutionBlock, e.ErrLatestExecBlock
}

// ExecutionBlockByHash --
func (e *EngineClient) ExecutionBlockByHash(_ context.Context, h common.Hash) (*pb.ExecutionBlock, error) {
	b, ok := e.BlockByHashMap[h]
	if !ok {
		return nil, errors.New("block not found")
	}
	return b, e.ErrExecBlockByHash
}

func (e *EngineClient) ReconstructFullBellatrixBlock(
	ctx context.Context, blindedBlock interfaces.SignedBeaconBlock,
) (interfaces.SignedBeaconBlock, error) {
	if !blindedBlock.Block().IsBlinded() {
		return nil, errors.New("block must be blinded")
	}
	header, err := blindedBlock.Block().Body().ExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}
	payload, ok := e.ExecutionPayloadByBlockHash[bytesutil.ToBytes32(header.BlockHash)]
	if !ok {
		return nil, errors.New("block not found")
	}
	e.NumReconstructedPayloads++
	return wrapper.BuildSignedBeaconBlockFromExecutionPayload(blindedBlock, payload)
}

// ExecutionBlockByHashWithTxs --
func (e *EngineClient) ExecutionBlockByHashWithTxs(_ context.Context, h common.Hash) (*pb.ExecutionBlockWithTxs, error) {
	b, ok := e.BlockWithTxsByHashMap[h]
	if !ok {
		return nil, errors.New("block not found")
	}
	return b, e.ErrExecBlockByHash
}

// GetTerminalBlockHash --
func (e *EngineClient) GetTerminalBlockHash(ctx context.Context) ([]byte, bool, error) {
	ttd := new(big.Int)
	ttd.SetString(params.BeaconConfig().TerminalTotalDifficulty, 10)
	terminalTotalDifficulty, overflows := uint256.FromBig(ttd)
	if overflows {
		return nil, false, errors.New("could not convert terminal total difficulty to uint256")
	}
	blk, err := e.LatestExecutionBlock(ctx)
	if err != nil {
		return nil, false, errors.Wrap(err, "could not get latest execution block")
	}
	if blk == nil {
		return nil, false, errors.New("latest execution block is nil")
	}

	for {
		b, err := hexutil.DecodeBig(blk.TotalDifficulty)
		if err != nil {
			return nil, false, errors.Wrap(err, "could not convert total difficulty to uint256")
		}
		currentTotalDifficulty, _ := uint256.FromBig(b)
		blockReachedTTD := currentTotalDifficulty.Cmp(terminalTotalDifficulty) >= 0

		parentHash := bytesutil.ToBytes32(blk.ParentHash)
		if len(blk.ParentHash) == 0 || parentHash == params.BeaconConfig().ZeroHash {
			return nil, false, nil
		}
		parentBlk, err := e.ExecutionBlockByHash(ctx, parentHash)
		if err != nil {
			return nil, false, errors.Wrap(err, "could not get parent execution block")
		}
		if blockReachedTTD {
			b, err := hexutil.DecodeBig(parentBlk.TotalDifficulty)
			if err != nil {
				return nil, false, errors.Wrap(err, "could not convert total difficulty to uint256")
			}
			parentTotalDifficulty, _ := uint256.FromBig(b)
			parentReachedTTD := parentTotalDifficulty.Cmp(terminalTotalDifficulty) >= 0
			if !parentReachedTTD {
				return blk.Hash, true, nil
			}
		} else {
			return nil, false, nil
		}
		blk = parentBlk
	}
}
