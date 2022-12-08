package validator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	fastssz "github.com/prysmaticlabs/fastssz"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v3/config/features"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v3/crypto/hash"
	"github.com/prysmaticlabs/prysm/v3/crypto/rand"
	"github.com/prysmaticlabs/prysm/v3/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v3/time/slots"
)

// eth1DataMajorityVote determines the appropriate eth1data for a block proposal using
// an algorithm called Voting with the Majority. The algorithm works as follows:
//   - Determine the timestamp for the start slot for the eth1 voting period.
//   - Determine the earliest and latest timestamps that a valid block can have.
//   - Determine the first block not before the earliest timestamp. This block is the lower bound.
//   - Determine the last block not after the latest timestamp. This block is the upper bound.
//   - If the last block is too early, use current eth1data from the beacon state.
//   - Filter out votes on unknown blocks and blocks which are outside of the range determined by the lower and upper bounds.
//   - If no blocks are left after filtering votes, use eth1data from the latest valid block.
//   - Otherwise:
//   - Determine the vote with the highest count. Prefer the vote with the highest eth1 block height in the event of a tie.
//   - This vote's block is the eth1 block to use for the block proposal.
func (vs *Server) eth1DataMajorityVote(ctx context.Context, beaconState state.BeaconState) (*ethpb.Eth1Data, error) {
	ctx, cancel := context.WithTimeout(ctx, eth1dataTimeout)
	defer cancel()

	slot := beaconState.Slot()
	votingPeriodStartTime := vs.slotStartTime(slot)

	if vs.MockEth1Votes {
		return vs.mockETH1DataVote(ctx, slot)
	}
	if !vs.Eth1InfoFetcher.ExecutionClientConnected() {
		return vs.randomETH1DataVote(ctx, "ExecutionClientConnected == false")
	}
	eth1DataNotification = false

	genesisTime, _ := vs.Eth1InfoFetcher.GenesisExecutionChainInfo()
	followDistanceSeconds := params.BeaconConfig().Eth1FollowDistance * params.BeaconConfig().SecondsPerETH1Block
	latestValidTime := votingPeriodStartTime - followDistanceSeconds
	earliestValidTime := votingPeriodStartTime - 2*followDistanceSeconds

	// Special case for starting from a pre-mined genesis: the eth1 vote should be genesis until the chain has advanced
	// by ETH1_FOLLOW_DISTANCE. The head state should maintain the same ETH1Data until this condition has passed, so
	// trust the existing head for the right eth1 vote until we can get a meaningful value from the deposit contract.
	log.WithField("votingPeriodStartTime", votingPeriodStartTime-genesisTime).WithField("latestValidTime", latestValidTime-genesisTime).WithField("earliestValidTime", earliestValidTime-genesisTime).Info("eth1DataMajorityVote")
	if latestValidTime < genesisTime+followDistanceSeconds {
		log.WithField("genesisTime", genesisTime).WithField("latestValidTime", latestValidTime).Info("voting period before genesis + follow distance, return head eth1data")
		return vs.HeadFetcher.HeadETH1Data(), nil
	}

	lastBlockByLatestValidTime, err := vs.Eth1BlockFetcher.BlockByTimestamp(ctx, latestValidTime)
	if err != nil {
		log.WithError(err).Error("Could not get last block by latest valid time")
		return vs.randomETH1DataVote(ctx, "error from BlockByTimestamp")
	}
	if lastBlockByLatestValidTime.Time < earliestValidTime {
		ed := vs.HeadFetcher.HeadETH1Data()
		log.WithField("latestBlock.Time", lastBlockByLatestValidTime.Time).WithField("latestBlock.Height", lastBlockByLatestValidTime.Number).WithFields(logFields(ed)).WithField("condition", "lastBlockByLatestValidTime.Time < earliestValidTime").Info("eth1DataMajorityVote")
		return ed, nil
	}

	lastBlockDepositCount, lastBlockDepositRoot := vs.DepositFetcher.DepositsNumberAndRootAtHeight(ctx, lastBlockByLatestValidTime.Number)
	if lastBlockDepositCount == 0 {
		log.WithFields(logFields(vs.ChainStartFetcher.ChainStartEth1Data())).WithField("condition", "ChainStartEth1Data").Info("eth1DataMajorityVote")
		return vs.ChainStartFetcher.ChainStartEth1Data(), nil
	}

	if lastBlockDepositCount >= vs.HeadFetcher.HeadETH1Data().DepositCount {
		h, err := vs.Eth1BlockFetcher.BlockHashByHeight(ctx, lastBlockByLatestValidTime.Number)
		if err != nil {
			log.WithError(err).Error("Could not get hash of last block by latest valid time")
			return vs.randomETH1DataVote(ctx, "deposit count too high")
		}
		ed := &ethpb.Eth1Data{
			BlockHash:    h.Bytes(),
			DepositCount: lastBlockDepositCount,
			DepositRoot:  lastBlockDepositRoot[:],
		}
		log.WithFields(logFields(ed)).WithField("condition", "lastBlockDepositCount >= vs.HeadFetcher.HeadETH1Data().DepositCount").Info("eth1DataMajorityVote")
		return ed, nil
	}
	return vs.HeadFetcher.HeadETH1Data(), nil
}

func logFields(ed *ethpb.Eth1Data) map[string]interface{} {
	fields := make(map[string]interface{})
	fields["deposit_root"] = fmt.Sprintf("%#x", ed.DepositRoot)
	fields["deposit_count"] = ed.DepositCount
	fields["block_hash"] = fmt.Sprintf("%#x", ed.BlockHash)
	return fields
}

func (vs *Server) slotStartTime(slot types.Slot) uint64 {
	startTime, _ := vs.Eth1InfoFetcher.GenesisExecutionChainInfo()
	return slots.VotingPeriodStartTime(startTime, slot)
}

// canonicalEth1Data determines the canonical eth1data and eth1 block height to use for determining deposits.
func (vs *Server) canonicalEth1Data(
	ctx context.Context,
	beaconState state.BeaconState,
	currentVote *ethpb.Eth1Data) (*ethpb.Eth1Data, *big.Int, error) {

	var eth1BlockHash [32]byte

	// Add in current vote, to get accurate vote tally
	if err := beaconState.AppendEth1DataVotes(currentVote); err != nil {
		return nil, nil, errors.Wrap(err, "could not append eth1 data votes to state")
	}
	hasSupport, err := blocks.Eth1DataHasEnoughSupport(beaconState, currentVote)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not determine if current eth1data vote has enough support")
	}
	var canonicalEth1Data *ethpb.Eth1Data
	if hasSupport {
		canonicalEth1Data = currentVote
		eth1BlockHash = bytesutil.ToBytes32(currentVote.BlockHash)
	} else {
		canonicalEth1Data = beaconState.Eth1Data()
		eth1BlockHash = bytesutil.ToBytes32(beaconState.Eth1Data().BlockHash)
	}
	if features.Get().DisableStakinContractCheck && eth1BlockHash == [32]byte{} {
		return canonicalEth1Data, new(big.Int).SetInt64(0), nil
	}
	_, canonicalEth1DataHeight, err := vs.Eth1BlockFetcher.BlockExists(ctx, eth1BlockHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not fetch eth1data height")
	}
	return canonicalEth1Data, canonicalEth1DataHeight, nil
}

func (vs *Server) mockETH1DataVote(ctx context.Context, slot types.Slot) (*ethpb.Eth1Data, error) {
	if !eth1DataNotification {
		log.Warn("Beacon Node is no longer connected to an ETH1 chain, so ETH1 data votes are now mocked.")
		eth1DataNotification = true
	}
	// If a mock eth1 data votes is specified, we use the following for the
	// eth1data we provide to every proposer based on https://github.com/ethereum/eth2.0-pm/issues/62:
	//
	// slot_in_voting_period = current_slot % SLOTS_PER_ETH1_VOTING_PERIOD
	// Eth1Data(
	//   DepositRoot = hash(current_epoch + slot_in_voting_period),
	//   DepositCount = state.eth1_deposit_index,
	//   BlockHash = hash(hash(current_epoch + slot_in_voting_period)),
	// )
	slotInVotingPeriod := slot.ModSlot(params.BeaconConfig().SlotsPerEpoch.Mul(uint64(params.BeaconConfig().EpochsPerEth1VotingPeriod)))
	headState, err := vs.HeadFetcher.HeadState(ctx)
	if err != nil {
		return nil, err
	}
	var enc []byte
	enc = fastssz.MarshalUint64(enc, uint64(slots.ToEpoch(slot))+uint64(slotInVotingPeriod))
	depRoot := hash.Hash(enc)
	blockHash := hash.Hash(depRoot[:])
	return &ethpb.Eth1Data{
		DepositRoot:  depRoot[:],
		DepositCount: headState.Eth1DepositIndex(),
		BlockHash:    blockHash[:],
	}, nil
}

func (vs *Server) randomETH1DataVote(ctx context.Context, reason string) (*ethpb.Eth1Data, error) {
	log.WithField("reason", reason).Warn("randomETH1DataVote")
	if !eth1DataNotification {
		log.Warn("Beacon Node is no longer connected to an ETH1 chain, so ETH1 data votes are now random.")
		eth1DataNotification = true
	}
	headState, err := vs.HeadFetcher.HeadState(ctx)
	if err != nil {
		return nil, err
	}

	// set random roots and block hashes to prevent a majority from being
	// built if the eth1 node is offline
	randGen := rand.NewGenerator()
	depRoot := hash.Hash(bytesutil.Bytes32(randGen.Uint64()))
	blockHash := hash.Hash(bytesutil.Bytes32(randGen.Uint64()))
	return &ethpb.Eth1Data{
		DepositRoot:  depRoot[:],
		DepositCount: headState.Eth1DepositIndex(),
		BlockHash:    blockHash[:],
	}, nil
}
