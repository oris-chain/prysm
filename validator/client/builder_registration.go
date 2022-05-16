package client

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/config/params"
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/time/slots"
	"go.opencensus.io/trace"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ProposeBuilderValidatorRegistration performs a build exit on a validator.
// The exit is signed by the validator before being sent to the beacon node for broadcasting.
func ProposeBuilderValidatorRegistration(
	ctx context.Context,
	validatorClient ethpb.BeaconNodeValidatorClient,
	nodeClient ethpb.NodeClient,
	signer signingFunc,
	reg *ethpb.ValidatorRegistrationV1,
) error {
	ctx, span := trace.StartSpan(ctx, "validator.ProposeBuilderValidatorRegistration")
	defer span.End()

	genesisResponse, err := nodeClient.GetGenesis(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.Wrap(err, "gRPC call to get genesis time failed")
	}
	ts := time.Unix(int64(reg.Timestamp), 0)
	secs := int64(ts.Second()) - genesisResponse.GenesisTime.Seconds
	currentSlot := types.Slot(uint64(secs) / params.BeaconConfig().SecondsPerSlot)

	_, err = signBuilderValidatorRegistration(ctx, currentSlot, validatorClient, signer, reg)
	if err != nil {
		return errors.Wrap(err, "failed to sign builder validator registration obj")
	}

	return nil
}

// Sings validator registration obj with proposer domain and private key.
func signBuilderValidatorRegistration(
	ctx context.Context,
	slot types.Slot,
	validatorClient ethpb.BeaconNodeValidatorClient,
	signer signingFunc,
	reg *ethpb.ValidatorRegistrationV1,
) ([]byte, error) {
	req := &ethpb.DomainRequest{
		Epoch:  slots.ToEpoch(slot),
		Domain: []byte{},
	}

	domain, err := validatorClient.DomainData(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, domainDataErr)
	}
	if domain == nil {
		return nil, errors.New(domainDataErr)
	}

	exitRoot, err := signing.ComputeSigningRoot(reg, domain.SignatureDomain)
	if err != nil {
		return nil, errors.Wrap(err, signingRootErr)
	}

	sig, err := signer(ctx, &validatorpb.SignRequest{
		PublicKey:       reg.Pubkey,
		SigningRoot:     exitRoot[:],
		SignatureDomain: domain.SignatureDomain,
		Object:          &validatorpb.SignRequest_Registration{Registration: reg},
	})
	if err != nil {
		return nil, errors.Wrap(err, signExitErr)
	}
	return sig.Marshal(), nil
}
