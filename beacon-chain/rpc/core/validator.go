package core

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/epoch/precompute"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/feed"
	opfeed "github.com/prysmaticlabs/prysm/v4/beacon-chain/core/feed/operation"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/helpers"
	coreTime "github.com/prysmaticlabs/prysm/v4/beacon-chain/core/time"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/transition"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/operations/synccommittee"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// AggregateBroadcastFailedError represents an error scenario where
// broadcasting an aggregate selection proof failed.
type AggregateBroadcastFailedError struct {
	err error
}

// NewAggregateBroadcastFailedError creates a new error instance.
func NewAggregateBroadcastFailedError(err error) AggregateBroadcastFailedError {
	return AggregateBroadcastFailedError{
		err: err,
	}
}

// Error returns the underlying error message.
func (e *AggregateBroadcastFailedError) Error() string {
	return fmt.Sprintf("could not broadcast signed aggregated attestation: %s", e.err.Error())
}

// ComputeValidatorPerformance reports the validator's latest balance along with other important metrics on
// rewards and penalties throughout its lifecycle in the beacon chain.
func ComputeValidatorPerformance(
	ctx context.Context,
	req *ethpb.ValidatorPerformanceRequest,
	headFetcher blockchain.HeadFetcher,
	currSlot primitives.Slot,
) (*ethpb.ValidatorPerformanceResponse, *RpcError) {
	headState, err := headFetcher.HeadState(ctx)
	if err != nil {
		return nil, &RpcError{Err: errors.Wrap(err, "could not get head state"), Reason: Internal}
	}
	if currSlot > headState.Slot() {
		headRoot, err := headFetcher.HeadRoot(ctx)
		if err != nil {
			return nil, &RpcError{Err: errors.Wrap(err, "could not get head root"), Reason: Internal}
		}
		headState, err = transition.ProcessSlotsUsingNextSlotCache(ctx, headState, headRoot, currSlot)
		if err != nil {
			return nil, &RpcError{Err: errors.Wrapf(err, "could not process slots up to %d", currSlot), Reason: Internal}
		}
	}
	var validatorSummary []*precompute.Validator
	if headState.Version() == version.Phase0 {
		vp, bp, err := precompute.New(ctx, headState)
		if err != nil {
			return nil, &RpcError{Err: err, Reason: Internal}
		}
		vp, bp, err = precompute.ProcessAttestations(ctx, headState, vp, bp)
		if err != nil {
			return nil, &RpcError{Err: err, Reason: Internal}
		}
		headState, err = precompute.ProcessRewardsAndPenaltiesPrecompute(headState, bp, vp, precompute.AttestationsDelta, precompute.ProposersDelta)
		if err != nil {
			return nil, &RpcError{Err: err, Reason: Internal}
		}
		validatorSummary = vp
	} else if headState.Version() >= version.Altair {
		vp, bp, err := altair.InitializePrecomputeValidators(ctx, headState)
		if err != nil {
			return nil, &RpcError{Err: err, Reason: Internal}
		}
		vp, bp, err = altair.ProcessEpochParticipation(ctx, headState, bp, vp)
		if err != nil {
			return nil, &RpcError{Err: err, Reason: Internal}
		}
		headState, vp, err = altair.ProcessInactivityScores(ctx, headState, vp)
		if err != nil {
			return nil, &RpcError{Err: err, Reason: Internal}
		}
		headState, err = altair.ProcessRewardsAndPenaltiesPrecompute(headState, bp, vp)
		if err != nil {
			return nil, &RpcError{Err: err, Reason: Internal}
		}
		validatorSummary = vp
	} else {
		return nil, &RpcError{Err: errors.Wrapf(err, "head state version %d not supported", headState.Version()), Reason: Internal}
	}

	responseCap := len(req.Indices) + len(req.PublicKeys)
	validatorIndices := make([]primitives.ValidatorIndex, 0, responseCap)
	missingValidators := make([][]byte, 0, responseCap)

	filtered := map[primitives.ValidatorIndex]bool{} // Track filtered validators to prevent duplication in the response.
	// Convert the list of validator public keys to validator indices and add to the indices set.
	for _, pubKey := range req.PublicKeys {
		// Skip empty public key.
		if len(pubKey) == 0 {
			continue
		}
		pubkeyBytes := bytesutil.ToBytes48(pubKey)
		idx, ok := headState.ValidatorIndexByPubkey(pubkeyBytes)
		if !ok {
			// Validator index not found, track as missing.
			missingValidators = append(missingValidators, pubKey)
			continue
		}
		if !filtered[idx] {
			validatorIndices = append(validatorIndices, idx)
			filtered[idx] = true
		}
	}
	// Add provided indices to the indices set.
	for _, idx := range req.Indices {
		if !filtered[idx] {
			validatorIndices = append(validatorIndices, idx)
			filtered[idx] = true
		}
	}
	// Depending on the indices and public keys given, results might not be sorted.
	sort.Slice(validatorIndices, func(i, j int) bool {
		return validatorIndices[i] < validatorIndices[j]
	})

	currentEpoch := coreTime.CurrentEpoch(headState)
	responseCap = len(validatorIndices)
	pubKeys := make([][]byte, 0, responseCap)
	beforeTransitionBalances := make([]uint64, 0, responseCap)
	afterTransitionBalances := make([]uint64, 0, responseCap)
	effectiveBalances := make([]uint64, 0, responseCap)
	correctlyVotedSource := make([]bool, 0, responseCap)
	correctlyVotedTarget := make([]bool, 0, responseCap)
	correctlyVotedHead := make([]bool, 0, responseCap)
	inactivityScores := make([]uint64, 0, responseCap)
	// Append performance summaries.
	// Also track missing validators using public keys.
	for _, idx := range validatorIndices {
		val, err := headState.ValidatorAtIndexReadOnly(idx)
		if err != nil {
			return nil, &RpcError{Err: errors.Wrap(err, "could not get validator"), Reason: Internal}
		}
		pubKey := val.PublicKey()
		if uint64(idx) >= uint64(len(validatorSummary)) {
			// Not listed in validator summary yet; treat it as missing.
			missingValidators = append(missingValidators, pubKey[:])
			continue
		}
		if !helpers.IsActiveValidatorUsingTrie(val, currentEpoch) {
			// Inactive validator; treat it as missing.
			missingValidators = append(missingValidators, pubKey[:])
			continue
		}

		summary := validatorSummary[idx]
		pubKeys = append(pubKeys, pubKey[:])
		effectiveBalances = append(effectiveBalances, summary.CurrentEpochEffectiveBalance)
		beforeTransitionBalances = append(beforeTransitionBalances, summary.BeforeEpochTransitionBalance)
		afterTransitionBalances = append(afterTransitionBalances, summary.AfterEpochTransitionBalance)
		correctlyVotedTarget = append(correctlyVotedTarget, summary.IsPrevEpochTargetAttester)
		correctlyVotedHead = append(correctlyVotedHead, summary.IsPrevEpochHeadAttester)

		if headState.Version() == version.Phase0 {
			correctlyVotedSource = append(correctlyVotedSource, summary.IsPrevEpochAttester)
		} else {
			correctlyVotedSource = append(correctlyVotedSource, summary.IsPrevEpochSourceAttester)
			inactivityScores = append(inactivityScores, summary.InactivityScore)
		}
	}

	return &ethpb.ValidatorPerformanceResponse{
		PublicKeys:                    pubKeys,
		CorrectlyVotedSource:          correctlyVotedSource,
		CorrectlyVotedTarget:          correctlyVotedTarget, // In altair, when this is true then the attestation was definitely included.
		CorrectlyVotedHead:            correctlyVotedHead,
		CurrentEffectiveBalances:      effectiveBalances,
		BalancesBeforeEpochTransition: beforeTransitionBalances,
		BalancesAfterEpochTransition:  afterTransitionBalances,
		MissingValidators:             missingValidators,
		InactivityScores:              inactivityScores, // Only populated in Altair
	}, nil
}

// SubmitSignedContributionAndProof is called by a sync committee aggregator
// to submit signed contribution and proof object.
func SubmitSignedContributionAndProof(
	ctx context.Context,
	s *ethpb.SignedContributionAndProof,
	broadcaster p2p.Broadcaster,
	pool synccommittee.Pool,
	notifier opfeed.Notifier,
) *RpcError {
	errs, ctx := errgroup.WithContext(ctx)

	// Broadcasting and saving contribution into the pool in parallel. As one fail should not affect another.
	errs.Go(func() error {
		return broadcaster.Broadcast(ctx, s)
	})

	if err := pool.SaveSyncCommitteeContribution(s.Message.Contribution); err != nil {
		return &RpcError{Err: err, Reason: Internal}
	}

	// Wait for p2p broadcast to complete and return the first error (if any)
	err := errs.Wait()
	if err != nil {
		return &RpcError{Err: err, Reason: Internal}
	}

	notifier.OperationFeed().Send(&feed.Event{
		Type: opfeed.SyncCommitteeContributionReceived,
		Data: &opfeed.SyncCommitteeContributionReceivedData{
			Contribution: s,
		},
	})

	return nil
}

// SubmitSignedAggregateSelectionProof verifies given aggregate and proofs and publishes them on appropriate gossipsub topic.
func SubmitSignedAggregateSelectionProof(
	ctx context.Context,
	req *ethpb.SignedAggregateSubmitRequest,
	genesisTime time.Time,
	broadcaster p2p.Broadcaster,
) *RpcError {
	if req.SignedAggregateAndProof == nil || req.SignedAggregateAndProof.Message == nil ||
		req.SignedAggregateAndProof.Message.Aggregate == nil || req.SignedAggregateAndProof.Message.Aggregate.Data == nil {
		return &RpcError{Err: errors.New("signed aggregate request can't be nil"), Reason: BadRequest}
	}
	emptySig := make([]byte, fieldparams.BLSSignatureLength)
	if bytes.Equal(req.SignedAggregateAndProof.Signature, emptySig) ||
		bytes.Equal(req.SignedAggregateAndProof.Message.SelectionProof, emptySig) {
		return &RpcError{Err: errors.New("signed signatures can't be zero hashes"), Reason: BadRequest}
	}

	// As a preventive measure, a beacon node shouldn't broadcast an attestation whose slot is out of range.
	if err := helpers.ValidateAttestationTime(req.SignedAggregateAndProof.Message.Aggregate.Data.Slot,
		genesisTime, params.BeaconNetworkConfig().MaximumGossipClockDisparity); err != nil {
		return &RpcError{Err: errors.New("attestation slot is no longer valid from current time"), Reason: BadRequest}
	}

	if err := broadcaster.Broadcast(ctx, req.SignedAggregateAndProof); err != nil {
		return &RpcError{Err: &AggregateBroadcastFailedError{err: err}, Reason: Internal}
	}

	log.WithFields(logrus.Fields{
		"slot":            req.SignedAggregateAndProof.Message.Aggregate.Data.Slot,
		"committeeIndex":  req.SignedAggregateAndProof.Message.Aggregate.Data.CommitteeIndex,
		"validatorIndex":  req.SignedAggregateAndProof.Message.AggregatorIndex,
		"aggregatedCount": req.SignedAggregateAndProof.Message.Aggregate.AggregationBits.Count(),
	}).Debug("Broadcasting aggregated attestation and proof")

	return nil
}
