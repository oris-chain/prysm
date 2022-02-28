package checkpoint

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/api/client/openapi"
	"github.com/prysmaticlabs/prysm/proto/sniff"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var saveFlags = struct {
	BeaconNodeHost string
	Timeout        string
}{}

var saveCmd = &cli.Command{
	Name:   "save",
	Usage:  "query for the current weak subjectivity period epoch, then download the corresponding state and block. To be used for checkpoint sync.",
	Action: cliActionSave,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "beacon-node-host",
			Usage:       "host:port for beacon node connection",
			Destination: &saveFlags.BeaconNodeHost,
			Value:       "localhost:3500",
		},
		&cli.StringFlag{
			Name:        "http-timeout",
			Usage:       "timeout for http requests made to beacon-node-url (uses duration format, ex: 2m31s). default: 2m",
			Destination: &saveFlags.Timeout,
			Value:       "4m",
		},
	},
}

func cliActionSave(_ *cli.Context) error {
	f := saveFlags
	opts := make([]openapi.ClientOpt, 0)
	log.Printf("--beacon-node-url=%s", f.BeaconNodeHost)
	timeout, err := time.ParseDuration(f.Timeout)
	if err != nil {
		return err
	}
	opts = append(opts, openapi.WithTimeout(timeout))
	client, err := openapi.NewClient(saveFlags.BeaconNodeHost, opts...)
	if err != nil {
		return err
	}

	return saveCheckpoint(client)
}

func saveCheckpoint(client *openapi.Client) error {
	ctx := context.Background()

	od, err := openapi.DownloadOriginData(ctx, client)
	if err != nil {
		log.Fatalf(err.Error())
	}

	blockPath := fname("block", od.ConfigFork, od.Block.Block().Slot(), od.WeakSubjectivity.BlockRoot)
	log.Printf("saving ssz-encoded block to to %s", blockPath)
	err = os.WriteFile(blockPath, od.BlockBytes, 0600)
	if err != nil {
		return err
	}

	stateRoot, err := od.State.HashTreeRoot(ctx)
	if err != nil {
		return errors.Wrap(err, "Could not compute HTR of state downloaded from remote beacon node")
	}
	statePath := fname("state", od.ConfigFork, od.State.Slot(), stateRoot)
	log.Printf("saving ssz-encoded state to to %s", statePath)
	err = os.WriteFile(statePath, od.StateBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func fname(prefix string, cf *sniff.ConfigFork, slot types.Slot, root [32]byte) string {
	return fmt.Sprintf("%s_%s_%s_%d-%#x.ssz", prefix, cf.ConfigName.String(), cf.Fork.String(), slot, root)
}

/*
func _saveCheckpoint(client *openapi.Client) error {
	ctx := context.Background()

	headReader, err := client.GetStateById(openapi.StateIdHead)
	headBytes, err := io.ReadAll(headReader)
	if err != nil {
		return errors.Wrap(err, "failed to read response body for get head state api call")
	}
	log.Printf("state response byte len=%d", len(headBytes))
	headState, err := sniff.BeaconState(headBytes)
	if err != nil {
		return errors.Wrap(err, "error unmarshaling state to correct version")
	}
	cf, err := sniff.ConfigForkForState(headBytes)
	if err != nil {
		return errors.Wrap(err, "error detecting chain config for beacon state")
	}
	log.Printf("detected supported config for state & block version detection, name=%s, fork=%s", cf.ConfigName.String(), cf.Fork)

	// LatestWeakSubjectivityEpoch uses package-level vars from the params package, so we need to override it
	params.OverrideBeaconConfig(cf.Config)
	epoch, err := helpers.LatestWeakSubjectivityEpoch(ctx, headState)
	if err != nil {
		return errors.Wrap(err, "error computing the weak subjectivity epoch from head state")
	}

	// use first slot of the epoch for the block slot
	bSlot, err := slots.EpochStart(epoch)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error computing first slot of epoch=%d", epoch))
	}
	// assigning this variable to make it extra obvious that the state slot is different
	sSlot := bSlot + 1
	// using the state at (slot % 32 = 1) instead of the epoch boundary ensures the
	// next block applied to the state will have the block at the weak subjectivity checkpoint
	// as its parent, satisfying prysm's sync code current verification that the parent block is present in the db
	log.Printf("weak subjectivity epoch computed as %d. download block@slot=%d, state@slot=%d", epoch, bSlot, sSlot)

	stateReader, err := client.GetStateBySlot(uint64(sSlot))
	if err != nil {
		return err
	}
	stateBytes, err := io.ReadAll(stateReader)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read response body for state at slot %d from api", sSlot))
	}
	log.Printf("state response byte len=%d", len(stateBytes))
	state, err := sniff.BeaconStateForConfigFork(stateBytes, cf)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unmarshaling state using auto-detected schema failed for state at slot %d", sSlot))
	}
	stateRoot, err := state.HashTreeRoot(ctx)
	if err != nil {
		return err
	}
	log.Printf("retrieved state for checkpoint at slot %d, w/ root=%#x", sSlot, stateRoot)
	latestBlockRoot, err := state.LatestBlockHeader().HashTreeRoot()
	if err != nil {
		return err
	}

	// Fork schedule is used to query for chain config metadata that is used to unmarshal values with the correct type
	fs, err := client.GetForkSchedule()
	if err != nil {
		return err
	}
	// get the version and find the config fork
	version, err := fs.VersionForEpoch(epoch)
	if err != nil {
		return err
	}
	cf, err = sniff.FindConfigFork(version)
	if err != nil {
		return errors.Wrap(err, "beacon node provided an unrecognized fork schedule")
	}

	blockReader, err := client.GetBlockBySlot(uint64(bSlot))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to retrieve block bytes for slot %d from api", bSlot))
	}
	blockBytes, err := io.ReadAll(blockReader)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read response body for block at slot %d from api", bSlot))
	}
	block, err := sniff.BlockForConfigFork(blockBytes, cf)
	if err != nil {
		return errors.Wrap(err, "failed to detect config/version from block inspection")
	}
	blockRoot, err := block.Block().HashTreeRoot()
	if err != nil {
		return err
	}
	log.Printf("retrieved block at slot %d with root=%#x", bSlot, blockRoot)
	blockStateRoot := block.Block().StateRoot()
	log.Printf("retrieved block has state root=%#x", blockStateRoot)

	// we only want to provide checkpoints+state pairs where the state integrates the checkpoint block as its latest root
	// this ensures that when syncing begins from the provided state, the next block in the chain can find the
	// latest block in the db.
	if blockRoot == latestBlockRoot {
		log.Printf("State latest_block_header root matches block root=%#x", latestBlockRoot)
	} else {
		return fmt.Errorf("fatal error, state latest_block_header root=%#x, does not match block root=%#x", latestBlockRoot, blockRoot)
	}

	bb, err := block.MarshalSSZ()
	if err != nil {
		return err
	}
	blockPath := fname("block", cf, bSlot, blockRoot)
	log.Printf("saving ssz-encoded block to to %s", blockPath)
	err = os.WriteFile(blockPath, bb, 0600)
	if err != nil {
		return err
	}

	sb, err := state.MarshalSSZ()
	if err != nil {
		return err
	}
	statePath := fname("state", cf, sSlot, stateRoot)
	log.Printf("saving ssz-encoded state to to %s", statePath)
	err = os.WriteFile(statePath, sb, 0600)
	if err != nil {
		return err
	}

	fmt.Println("To validate that your client is using this checkpoint, specify the following flag when starting prysm:")
	fmt.Printf("--weak-subjectivity-checkpoint=%#x:%d\n\n", blockRoot, epoch)
	fmt.Println("To sync a new beacon node starting from the checkpoint state, you may specify the following flags (assuming the files are in your current working directory)")
	fmt.Printf("--checkpoint-block=%s --checkpoint-state=%s\n", blockPath, statePath)
	return nil
}
*/
