package validator

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prysmaticlabs/prysm/v4/api/client/builder"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v4/encoding/ssz"
	"github.com/prysmaticlabs/prysm/v4/monitoring/tracing"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// builderGetPayloadMissCount tracks the number of misses when validator tries to get a payload from builder
var builderGetPayloadMissCount = promauto.NewCounter(prometheus.CounterOpts{
	Name: "builder_get_payload_miss_count",
	Help: "The number of get payload misses for validator requests to builder",
})

var gweiPerEth = big.NewInt(int64(params.BeaconConfig().GweiPerEth))

// emptyTransactionsRoot represents the returned value of ssz.TransactionsRoot([][]byte{}) and
// can be used as a constant to avoid recomputing this value in every call.
var emptyTransactionsRoot = [32]byte{127, 254, 36, 30, 166, 1, 135, 253, 176, 24, 123, 250, 34, 222, 53, 209, 249, 190, 215, 171, 6, 29, 148, 1, 253, 71, 227, 74, 84, 251, 237, 225}

// blockBuilderTimeout is the maximum amount of time allowed for a block builder to respond to a
// block request. This value is known as `BUILDER_PROPOSAL_DELAY_TOLERANCE` in builder spec.
const blockBuilderTimeout = 1 * time.Second

// Sets the execution data for the block. Execution data can come from local EL client or remote builder depends on validator registration and circuit breaker conditions.
func (vs *Server) setExecutionData(ctx context.Context, blk interfaces.SignedBeaconBlock, headState state.BeaconState) ([]*enginev1.Blob, error) {
	ctx, span := trace.StartSpan(ctx, "ProposerServer.setExecutionData")
	defer span.End()

	idx := blk.Block().ProposerIndex()
	slot := blk.Block().Slot()
	if slots.ToEpoch(slot) < params.BeaconConfig().BellatrixForkEpoch {
		return nil, nil
	}

	canUseBuilder, err := vs.canUseBuilder(ctx, slot, idx)
	span.AddAttributes(trace.BoolAttribute("canUseBuilder", canUseBuilder))
	if err != nil {
		log.WithError(err).Warn("Proposer: failed to check if builder can be used")
	} else if canUseBuilder {
		builderPayload, err := vs.getPayloadHeaderFromBuilder(ctx, slot, idx)
		if err != nil {
			builderGetPayloadMissCount.Inc()
			log.WithError(err).Warn("Proposer: failed to get payload header from builder")
		} else {
			switch {
			case blk.Version() >= version.Capella:
				localPayload, _, err := vs.getExecutionPayload(ctx, slot, idx, blk.Block().ParentRoot(), headState)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get execution payload")
				}
				// Compare payload values between local and builder. Default to the local value if it is higher.
				v, err := localPayload.Value()
				if err != nil {
					return nil, errors.Wrap(err, "failed to get local payload value")
				}
				v.Div(v, gweiPerEth)
				localValue := v.Uint64()
				v, err = builderPayload.Value()
				if err != nil {
					log.WithError(err).Warn("Proposer: failed to get builder payload value") // Default to local if can't get builder value.
					v = big.NewInt(0)                                                        // Default to local if can't get builder value.
				}
				v.Div(v, gweiPerEth)
				builderValue := v.Uint64()

				withdrawalsMatched, err := matchingWithdrawalsRoot(localPayload, builderPayload)
				if err != nil {
					tracing.AnnotateError(span, err)
					return nil, errors.Wrap(err, "failed to match withdrawals root")
				}

				// Use builder payload if the following in true:
				// builder_bid_value * 100 > local_block_value * (local-block-value-boost + 100)
				boost := params.BeaconConfig().LocalBlockValueBoost
				higherValueBuilder := builderValue*100 > localValue*(100+boost)

				// If we can't get the builder value, just use local block.
				if higherValueBuilder && withdrawalsMatched { // Builder value is higher and withdrawals match.
					blk.SetBlinded(true)
					if err := blk.SetExecution(builderPayload); err != nil {
						log.WithError(err).Warn("Proposer: failed to set builder payload")
						blk.SetBlinded(false)
					} else {
						return nil, nil
					}
				}
				if !higherValueBuilder {
					log.WithFields(logrus.Fields{
						"localGweiValue":       localValue,
						"localBoostPercentage": boost,
						"builderGweiValue":     builderValue,
					}).Warn("Proposer: using local execution payload because higher value")
				}
				span.AddAttributes(
					trace.BoolAttribute("higherValueBuilder", higherValueBuilder),
					trace.Int64Attribute("localGweiValue", int64(localValue)),     // lint:ignore uintcast -- This is OK for tracing.
					trace.Int64Attribute("localBoostPercentage", int64(boost)),    // lint:ignore uintcast -- This is OK for tracing.
					trace.Int64Attribute("builderGweiValue", int64(builderValue)), // lint:ignore uintcast -- This is OK for tracing.
				)

				return nil, blk.SetExecution(localPayload)
			default: // Bellatrix case.
				blk.SetBlinded(true)
				if err := blk.SetExecution(builderPayload); err != nil {
					log.WithError(err).Warn("Proposer: failed to set builder payload")
					blk.SetBlinded(false)
				} else {
					return nil, nil
				}
			}
		}
	}

	executionData, blobsBundle, err := vs.getExecutionPayload(ctx, slot, idx, blk.Block().ParentRoot(), headState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get execution payload")
	}
	if slots.ToEpoch(slot) >= params.BeaconConfig().DenebForkEpoch && len(blobsBundle.KzgCommitments) > 0 {
		// TODO: check block hash matches blob bundle hash
		if err := blk.SetBlobKzgCommitments(blobsBundle.KzgCommitments); err != nil {
			return nil, errors.Wrap(err, "could not set blob kzg commitments")
		}
		return blobsBundle.Blobs, blk.SetExecution(executionData)
	}

	return nil, blk.SetExecution(executionData)
}

// This function retrieves the payload header given the slot number and the validator index.
// It's a no-op if the latest head block is not versioned bellatrix.
func (vs *Server) getPayloadHeaderFromBuilder(ctx context.Context, slot primitives.Slot, idx primitives.ValidatorIndex) (interfaces.ExecutionData, error) {
	ctx, span := trace.StartSpan(ctx, "ProposerServer.getPayloadHeaderFromBuilder")
	defer span.End()

	if slots.ToEpoch(slot) < params.BeaconConfig().BellatrixForkEpoch {
		return nil, errors.New("can't get payload header from builder before bellatrix epoch")
	}

	b, err := vs.HeadFetcher.HeadBlock(ctx)
	if err != nil {
		return nil, err
	}

	h, err := b.Block().Body().Execution()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get execution header")
	}
	pk, err := vs.HeadFetcher.HeadValidatorIndexToPublicKey(ctx, idx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, blockBuilderTimeout)
	defer cancel()

	signedBid, err := vs.BlockBuilder.GetHeader(ctx, slot, bytesutil.ToBytes32(h.BlockHash()), pk)
	if err != nil {
		return nil, err
	}
	if signedBid.IsNil() {
		return nil, errors.New("builder returned nil bid")
	}
	if signedBid.Version() != b.Version() {
		return nil, fmt.Errorf("builder bid response version: %d is different from head block version: %d", signedBid.Version(), b.Version())
	}
	bid, err := signedBid.Message()
	if err != nil {
		return nil, errors.Wrap(err, "could not get bid")
	}
	if bid.IsNil() {
		return nil, errors.New("builder returned nil bid")
	}

	v := bytesutil.LittleEndianBytesToBigInt(bid.Value())
	if v.String() == "0" {
		return nil, errors.New("builder returned header with 0 bid amount")
	}

	header, err := bid.Header()
	if err != nil {
		return nil, errors.Wrap(err, "could not get bid header")
	}
	txRoot, err := header.TransactionsRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not get transaction root")
	}
	if bytesutil.ToBytes32(txRoot) == emptyTransactionsRoot {
		return nil, errors.New("builder returned header with an empty tx root")
	}

	if !bytes.Equal(header.ParentHash(), h.BlockHash()) {
		return nil, fmt.Errorf("incorrect parent hash %#x != %#x", header.ParentHash(), h.BlockHash())
	}

	t, err := slots.ToTime(uint64(vs.TimeFetcher.GenesisTime().Unix()), slot)
	if err != nil {
		return nil, err
	}
	if header.Timestamp() != uint64(t.Unix()) {
		return nil, fmt.Errorf("incorrect timestamp %d != %d", header.Timestamp(), uint64(t.Unix()))
	}

	if err := validateBuilderSignature(signedBid); err != nil {
		return nil, errors.Wrap(err, "could not validate builder signature")
	}

	log.WithFields(logrus.Fields{
		"value":         v.String(),
		"builderPubKey": fmt.Sprintf("%#x", bid.Pubkey()),
		"blockHash":     fmt.Sprintf("%#x", header.BlockHash()),
	}).Info("Received header with bid")

	span.AddAttributes(
		trace.StringAttribute("value", v.String()),
		trace.StringAttribute("builderPubKey", fmt.Sprintf("%#x", bid.Pubkey())),
		trace.StringAttribute("blockHash", fmt.Sprintf("%#x", header.BlockHash())),
	)

	return header, nil
}

// Validates builder signature and returns an error if the signature is invalid.
func validateBuilderSignature(signedBid builder.SignedBid) error {
	d, err := signing.ComputeDomain(params.BeaconConfig().DomainApplicationBuilder,
		nil, /* fork version */
		nil /* genesis val root */)
	if err != nil {
		return err
	}
	if signedBid.IsNil() {
		return errors.New("nil builder bid")
	}
	bid, err := signedBid.Message()
	if err != nil {
		return errors.Wrap(err, "could not get bid")
	}
	if bid.IsNil() {
		return errors.New("builder returned nil bid")
	}
	return signing.VerifySigningRoot(bid, bid.Pubkey(), signedBid.Signature(), d)
}

func matchingWithdrawalsRoot(local, builder interfaces.ExecutionData) (bool, error) {
	wds, err := local.Withdrawals()
	if err != nil {
		return false, errors.Wrap(err, "could not get local withdrawals")
	}
	br, err := builder.WithdrawalsRoot()
	if err != nil {
		return false, errors.Wrap(err, "could not get builder withdrawals root")
	}
	wr, err := ssz.WithdrawalSliceRoot(wds, fieldparams.MaxWithdrawalsPerPayload)
	if err != nil {
		return false, errors.Wrap(err, "could not compute local withdrawals root")
	}

	if !bytes.Equal(br, wr[:]) {
		log.WithFields(logrus.Fields{
			"local":   fmt.Sprintf("%#x", wr),
			"builder": fmt.Sprintf("%#x", br),
		}).Warn("Proposer: withdrawal roots don't match, using local block")
		return false, nil
	}
	return true, nil
}
