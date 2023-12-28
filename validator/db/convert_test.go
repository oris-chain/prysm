package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	validator_service_config "github.com/prysmaticlabs/prysm/v4/config/validator/service"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/validator/db/filesystem"
	"github.com/prysmaticlabs/prysm/v4/validator/db/iface"

	"github.com/prysmaticlabs/prysm/v4/validator/db/kv"
)

func getPubkeyFromString(t *testing.T, pubkeyString string) [fieldparams.BLSPubkeyLength]byte {
	var pubkey [fieldparams.BLSPubkeyLength]byte
	pubkeyBytes, err := hexutil.Decode(pubkeyString)
	require.NoError(t, err, "hexutil.Decode should not return an error")
	copy(pubkey[:], pubkeyBytes)
	return pubkey
}

func getFeeRecipientFromString(t *testing.T, feeRecipientString string) [fieldparams.FeeRecipientLength]byte {
	var feeRecipient [fieldparams.FeeRecipientLength]byte
	feeRecipientBytes, err := hexutil.Decode(feeRecipientString)
	require.NoError(t, err, "hexutil.Decode should not return an error")
	copy(feeRecipient[:], feeRecipientBytes)
	return feeRecipient
}

func TestDB_IsCompleteDatabaseExisting(t *testing.T) {
	ctx := context.Background()

	// Create a directory path.
	datadir := t.TempDir()

	// Check if the complete database exists.
	completeDatabaseExists, err := IsCompleteDatabaseExisting(datadir)
	require.NoError(t, err, "could not check if complete database exists")
	require.Equal(t, false, completeDatabaseExists, "complete database should not exist")

	// Create the complete database.
	completeStore, err := kv.NewKVStore(ctx, datadir, nil)
	require.NoError(t, err, "could not create complete database")

	// Check if the complete database exists.
	completeDatabaseExists, err = IsCompleteDatabaseExisting(datadir)
	require.NoError(t, err, "could not check if complete database exists")
	require.Equal(t, true, completeDatabaseExists, "complete database should not exist")

	// Close the complete database.
	err = completeStore.Close()
	require.NoError(t, err, "could not close complete database")
}

func TestDB_IsMinimalDatabaseExisting(t *testing.T) {
	ctx := context.Background()

	// Create a directory path.
	datadir := t.TempDir()

	// Check if the minimal database exists.
	minimalDatabaseExists, err := IsMinimalDatabaseExisting(datadir)
	require.NoError(t, err, "could not check if minimal database exists")
	require.Equal(t, false, minimalDatabaseExists, "minimal database should not exist")

	// Create the minimal database.
	minimalStore, err := filesystem.NewStore(datadir, nil)
	require.NoError(t, err, "could not create minimal database")

	// Save the genesis validator root.
	err = minimalStore.SaveGenesisValidatorsRoot(ctx, []byte("genesis-validator-root"))
	require.NoError(t, err, "could not save genesis validator root")

	// Check if the minimal database exists.
	minimalDatabaseExists, err = IsMinimalDatabaseExisting(datadir)
	require.NoError(t, err, "could not check if minimal database exists")
	require.Equal(t, true, minimalDatabaseExists, "minimal database should not exist")

	// Close the complete database.
	err = minimalStore.Close()
	require.NoError(t, err, "could not close minimal database")
}

func TestDB_ConvertDatabase(t *testing.T) {
	ctx := context.Background()

	pubKeyString1 := "0x80000060606fa05c7339dd7bcd0d3e4d8b573fa30dea2fdb4997031a703e3300326e3c054be682f92d9c367cd647bbea"
	pubKeyString2 := "0x81000060606fa05c7339dd7bcd0d3e4d8b573fa30dea2fdb4997031a703e3300326e3c054be682f92d9c367cd647bbea"
	defaultFeeRecipientString := "0xe688b84b23f322a994A53dbF8E15FA82CDB71127"
	customFeeRecipientString := "0xeD33259a056F4fb449FFB7B7E2eCB43a9B5685Bf"

	pubkey1 := getPubkeyFromString(t, pubKeyString1)
	pubkey2 := getPubkeyFromString(t, pubKeyString2)
	defaultFeeRecipient := getFeeRecipientFromString(t, defaultFeeRecipientString)
	customFeeRecipient := getFeeRecipientFromString(t, customFeeRecipientString)

	for _, minimalToComplete := range []bool{false, true} {
		t.Run(fmt.Sprintf("minimalToComplete=%v", minimalToComplete), func(t *testing.T) {
			// Create signing root
			signingRoot := [fieldparams.RootLength]byte{}
			var signingRootBytes []byte
			if minimalToComplete {
				signingRootBytes = signingRoot[:]
			}

			// Create database directoriy path.
			datadir := t.TempDir()

			// Run source DB preparation.
			// --------------------------
			// Create the source database.
			var (
				sourceDatabase, targetDatabase iface.ValidatorDB
				err                            error
			)
			if minimalToComplete {
				sourceDatabase, err = filesystem.NewStore(datadir, &filesystem.Config{
					PubKeys: [][fieldparams.BLSPubkeyLength]byte{pubkey1, pubkey2},
				})
			} else {
				sourceDatabase, err = kv.NewKVStore(ctx, datadir, &kv.Config{
					PubKeys: [][fieldparams.BLSPubkeyLength]byte{pubkey1, pubkey2},
				})
			}

			require.NoError(t, err, "could not create source database")

			// Save the genesis validator root.
			expectedGenesisValidatorRoot := []byte("genesis-validator-root")
			err = sourceDatabase.SaveGenesisValidatorsRoot(ctx, expectedGenesisValidatorRoot)
			require.NoError(t, err, "could not save genesis validator root")

			// Save the graffiti file hash.
			// (Getting the graffiti ordered index will set the graffiti file hash)
			expectedGraffitiFileHash := [32]byte{1}
			_, err = sourceDatabase.GraffitiOrderedIndex(ctx, expectedGraffitiFileHash)
			require.NoError(t, err, "could not get graffiti ordered index")

			// Save the graffiti ordered index.
			expectedGraffitiOrderedIndex := uint64(1)
			err = sourceDatabase.SaveGraffitiOrderedIndex(ctx, expectedGraffitiOrderedIndex)
			require.NoError(t, err, "could not save graffiti ordered index")

			// Save the proposer settings.
			relays := []string{}
			if minimalToComplete {
				relays = nil
			}

			expectedProposerSettings := &validator_service_config.ProposerSettings{
				ProposeConfig: map[[fieldparams.BLSPubkeyLength]byte]*validator_service_config.ProposerOption{
					pubkey1: {
						FeeRecipientConfig: &validator_service_config.FeeRecipientConfig{
							FeeRecipient: customFeeRecipient,
						},
						BuilderConfig: &validator_service_config.BuilderConfig{
							Enabled:  true,
							GasLimit: 42,
							Relays:   relays,
						},
					},
				},
				DefaultConfig: &validator_service_config.ProposerOption{
					FeeRecipientConfig: &validator_service_config.FeeRecipientConfig{
						FeeRecipient: defaultFeeRecipient,
					},
					BuilderConfig: &validator_service_config.BuilderConfig{
						Enabled:  false,
						GasLimit: 43,
						Relays:   relays,
					},
				},
			}

			err = sourceDatabase.SaveProposerSettings(ctx, expectedProposerSettings)
			require.NoError(t, err, "could not save proposer settings")

			// Save some attestations.
			completeAttestations := []*ethpb.IndexedAttestation{
				{
					Data: &ethpb.AttestationData{
						Source: &ethpb.Checkpoint{
							Epoch: 1,
						},
						Target: &ethpb.Checkpoint{
							Epoch: 2,
						},
					},
				},
				{
					Data: &ethpb.AttestationData{
						Source: &ethpb.Checkpoint{
							Epoch: 2,
						},
						Target: &ethpb.Checkpoint{
							Epoch: 3,
						},
					},
				},
			}

			expectedAttestationRecords1 := []*kv.AttestationRecord{
				{
					PubKey:      pubkey1,
					Source:      primitives.Epoch(2),
					Target:      primitives.Epoch(3),
					SigningRoot: signingRootBytes,
				},
			}

			expectedAttestationRecords2 := []*kv.AttestationRecord{
				{
					PubKey:      pubkey2,
					Source:      primitives.Epoch(2),
					Target:      primitives.Epoch(3),
					SigningRoot: signingRootBytes,
				},
			}

			err = sourceDatabase.SaveAttestationsForPubKey(ctx, pubkey1, [][]byte{{1}, {2}}, completeAttestations)
			require.NoError(t, err, "could not save attestations")

			err = sourceDatabase.SaveAttestationsForPubKey(ctx, pubkey2, [][]byte{{1}, {2}}, completeAttestations)
			require.NoError(t, err, "could not save attestations")

			// Save some block proposals.
			err = sourceDatabase.SaveProposalHistoryForSlot(ctx, pubkey1, 42, []byte{})
			require.NoError(t, err, "could not save block proposal")

			err = sourceDatabase.SaveProposalHistoryForSlot(ctx, pubkey1, 43, []byte{})
			require.NoError(t, err, "could not save block proposal")

			expectedProposals := []*kv.Proposal{
				{
					Slot:        43,
					SigningRoot: signingRootBytes,
				},
			}

			// Close the source database.
			err = sourceDatabase.Close()
			require.NoError(t, err, "could not close source database")

			// Source to target DB conversion.
			// ----------------------------------------
			err = ConvertDatabase(ctx, datadir, datadir, minimalToComplete)
			require.NoError(t, err, "could not convert source to target database")

			// Check the target database.
			// --------------------------
			if minimalToComplete {
				targetDatabase, err = kv.NewKVStore(ctx, datadir, nil)
			} else {
				targetDatabase, err = filesystem.NewStore(datadir, nil)
			}
			require.NoError(t, err, "could not get minimal database")

			// Check the genesis validator root.
			actualGenesisValidatoRoot, err := targetDatabase.GenesisValidatorsRoot(ctx)
			require.NoError(t, err, "could not get genesis validator root from target database")
			require.DeepSSZEqual(t, expectedGenesisValidatorRoot, actualGenesisValidatoRoot, "genesis validator root should match")

			// Check the graffiti file hash.
			actualGraffitiFileHash, exists, err := targetDatabase.GraffitiFileHash()
			require.NoError(t, err, "could not get graffiti file hash from target database")
			require.Equal(t, true, exists, "graffiti file hash should exist")
			require.Equal(t, expectedGraffitiFileHash, actualGraffitiFileHash, "graffiti file hash should match")

			// Check the graffiti ordered index.
			actualGraffitiOrderedIndex, err := targetDatabase.GraffitiOrderedIndex(ctx, expectedGraffitiFileHash)
			require.NoError(t, err, "could not get graffiti ordered index from target database")
			require.Equal(t, expectedGraffitiOrderedIndex, actualGraffitiOrderedIndex, "graffiti ordered index should match")

			// Check the proposer settings.
			actualProposerSettings, err := targetDatabase.ProposerSettings(ctx)
			require.NoError(t, err, "could not get proposer settings from target database")
			require.DeepEqual(t, expectedProposerSettings, actualProposerSettings, "proposer settings should match")

			// Check the attestations.
			actualAttestationRecords, err := targetDatabase.AttestationHistoryForPubKey(ctx, pubkey1)
			require.NoError(t, err, "could not get attestations from target database")
			require.DeepEqual(t, expectedAttestationRecords1, actualAttestationRecords, "attestations should match")

			actualAttestationRecords, err = targetDatabase.AttestationHistoryForPubKey(ctx, pubkey2)
			require.NoError(t, err, "could not get attestations from target database")
			require.DeepEqual(t, expectedAttestationRecords2, actualAttestationRecords, "attestations should match")

			// Check the block proposals.
			actualProposals, err := targetDatabase.ProposalHistoryForPubKey(ctx, pubkey1)
			require.NoError(t, err, "could not get block proposals from target database")
			require.DeepEqual(t, expectedProposals, actualProposals, "block proposals should match")

			// Close the target database.
			err = targetDatabase.Close()
			require.NoError(t, err, "could not close target database")

			// Check the source database does not exist anymore.
			var existing bool

			if minimalToComplete {
				existing, err = IsMinimalDatabaseExisting(datadir)
			} else {
				existing, err = IsCompleteDatabaseExisting(datadir)
			}

			require.NoError(t, err, "could not check if source database exists")
			require.Equal(t, false, existing, "source database should not exist")
		})
	}
}
