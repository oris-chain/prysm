package flags

import (
	"github.com/prysmaticlabs/prysm/v3/cmd"
	"github.com/urfave/cli/v2"
)

// GlobalFlags specifies all the global flags for the
// beacon node.
type GlobalFlags struct {
	SubscribeToAllSubnets      bool
	MinimumSyncPeers           int
	MinimumPeersPerSubnet      int
	BlockBatchLimit            int
	BlockBatchLimitBurstFactor int
	BlobsTransferRate          int
	BlobsTransferRateThresh    int
}

var globalConfig *GlobalFlags

// Get retrieves the global config.
func Get() *GlobalFlags {
	if globalConfig == nil {
		return &GlobalFlags{}
	}
	return globalConfig
}

// Init sets the global config equal to the config that is passed in.
func Init(c *GlobalFlags) {
	globalConfig = c
}

// ConfigureGlobalFlags initializes the global config.
// based on the provided cli context.
func ConfigureGlobalFlags(ctx *cli.Context) {
	cfg := &GlobalFlags{}
	if ctx.Bool(SubscribeToAllSubnets.Name) {
		log.Warn("Subscribing to All Attestation Subnets")
		cfg.SubscribeToAllSubnets = true
	}
	// TODO(EIP-4844): assert BlockBatchLimit < MAX_REQUEST_BLOBS_SIDECARS (128) as initial-sync uses this same config for sidecar fetch
	cfg.BlockBatchLimit = ctx.Int(BlockBatchLimit.Name)
	cfg.BlockBatchLimitBurstFactor = ctx.Int(BlockBatchLimitBurstFactor.Name)
	cfg.BlobsTransferRate = ctx.Int(BlobsTransferRate.Name)
	cfg.BlobsTransferRateThresh = ctx.Int(BlobsTransferRateThresh.Name)
	cfg.MinimumPeersPerSubnet = ctx.Int(MinPeersPerSubnet.Name)
	configureMinimumPeers(ctx, cfg)

	Init(cfg)
}

func configureMinimumPeers(ctx *cli.Context, cfg *GlobalFlags) {
	cfg.MinimumSyncPeers = ctx.Int(MinSyncPeers.Name)
	maxPeers := ctx.Int(cmd.P2PMaxPeers.Name)
	if cfg.MinimumSyncPeers > maxPeers {
		log.Warnf("Changing Minimum Sync Peers to %d", maxPeers)
		cfg.MinimumSyncPeers = maxPeers
	}
}
