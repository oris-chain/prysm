package testnet

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/cmd/flags"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/container/trie"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/interop"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	generateGenesisStateFlags = struct {
		DepositJsonFile    string
		ChainConfigFile    string
		ConfigName         string
		NumValidators      uint64
		GenesisTime        uint64
		GenesisTimeDelay        uint64
		OutputSSZ          string
		OutputJSON         string
		OutputYaml         string
		ForkName           string
		OverrideEth1Data   bool
		ExecutionEndpoint  string
		GethGenesisJsonIn  string
		GethGenesisJsonOut string
	}{}
	log           = logrus.WithField("prefix", "genesis")
	outputSSZFlag = &cli.StringFlag{
		Name:        "output-ssz",
		Destination: &generateGenesisStateFlags.OutputSSZ,
		Usage:       "Output filename of the SSZ marshaling of the generated genesis state",
		Value:       "",
	}
	outputYamlFlag = &cli.StringFlag{
		Name:        "output-yaml",
		Destination: &generateGenesisStateFlags.OutputYaml,
		Usage:       "Output filename of the YAML marshaling of the generated genesis state",
		Value:       "",
	}
	outputJsonFlag = &cli.StringFlag{
		Name:        "output-json",
		Destination: &generateGenesisStateFlags.OutputJSON,
		Usage:       "Output filename of the JSON marshaling of the generated genesis state",
		Value:       "",
	}
	generateGenesisStateCmd = &cli.Command{
		Name:  "generate-genesis",
		Usage: "Generate a beacon chain genesis state",
		Action: func(cliCtx *cli.Context) error {
			if err := cliActionGenerateGenesisState(cliCtx); err != nil {
				log.WithError(err).Fatal("Could not generate beacon chain genesis state")
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "chain-config-file",
				Destination: &generateGenesisStateFlags.ChainConfigFile,
				Usage:       "The path to a YAML file with chain config values",
			},
			&cli.StringFlag{
				Name:        "deposit-json-file",
				Destination: &generateGenesisStateFlags.DepositJsonFile,
				Usage:       "Path to deposit_data.json file generated by the staking-deposit-cli tool for optionally specifying validators in genesis state",
			},
			&cli.StringFlag{
				Name:        "config-name",
				Usage:       "Config kind to be used for generating the genesis state. Default: mainnet. Options include mainnet, interop, minimal, prater, sepolia. --chain-config-file will override this flag.",
				Destination: &generateGenesisStateFlags.ConfigName,
				Value:       params.MainnetName,
			},
			&cli.Uint64Flag{
				Name:        "num-validators",
				Usage:       "Number of validators to deterministically generate in the genesis state",
				Destination: &generateGenesisStateFlags.NumValidators,
				Required:    true,
			},
			&cli.Uint64Flag{
				Name:        "genesis-time",
				Destination: &generateGenesisStateFlags.GenesisTime,
				Usage:       "Unix timestamp seconds used as the genesis time in the genesis state. If unset, defaults to now()",
			},
			&cli.Uint64Flag{
				Name:        "genesis-time-delay",
				Destination: &generateGenesisStateFlags.GenesisTimeDelay,
				Usage:       "Delay genesis time by N seconds",
			},
			&cli.BoolFlag{
				Name:        "override-eth1data",
				Destination: &generateGenesisStateFlags.OverrideEth1Data,
				Usage:       "Overrides Eth1Data with values from execution client. If unset, defaults to false",
				Value:       false,
			},
			&cli.StringFlag{
				Name:        "geth-genesis-json-in",
				Destination: &generateGenesisStateFlags.GethGenesisJsonIn,
				Usage:       "Path to a \"genesis.json\" file, containing a json representation of Geth's core.Genesis",
			},
			&cli.StringFlag{
				Name:        "geth-genesis-json-out",
				Destination: &generateGenesisStateFlags.GethGenesisJsonOut,
				Usage:       "Path to write generated \"genesis.json\" file, containing a json representation of Geth's core.Genesis",
			},
			&cli.StringFlag{
				Name:        "execution-endpoint",
				Destination: &generateGenesisStateFlags.ExecutionEndpoint,
				Usage:       "Endpoint to preferred execution client. If unset, defaults to Geth",
				Value:       "http://localhost:8545",
			},
			flags.EnumValue{
				Name:        "fork",
				Usage:       fmt.Sprintf("Name of the BeaconState schema to use in output encoding [%s]", strings.Join(versionNames(), ",")),
				Enum:        versionNames(),
				Value:       versionNames()[0],
				Destination: &generateGenesisStateFlags.ForkName,
			}.GenericFlag(),
			outputSSZFlag,
			outputYamlFlag,
			outputJsonFlag,
		},
	}
)

func versionNames() []string {
	enum := version.All()
	names := make([]string, len(enum))
	for i := range enum {
		names[i] = version.String(enum[i])
	}
	return names
}

// Represents a json object of hex string and uint64 values for
// validators on Ethereum. This file can be generated using the official staking-deposit-cli.
type depositDataJSON struct {
	PubKey                string `json:"pubkey"`
	Amount                uint64 `json:"amount"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	DepositDataRoot       string `json:"deposit_data_root"`
	Signature             string `json:"signature"`
}

func cliActionGenerateGenesisState(cliCtx *cli.Context) error {
	outputJson := generateGenesisStateFlags.OutputJSON
	outputYaml := generateGenesisStateFlags.OutputYaml
	outputSSZ := generateGenesisStateFlags.OutputSSZ
	noOutputFlag := outputSSZ == "" && outputJson == "" && outputYaml == ""
	if noOutputFlag {
		return fmt.Errorf(
			"no %s, %s, %s flag(s) specified. At least one is required",
			outputJsonFlag.Name,
			outputYamlFlag.Name,
			outputSSZFlag.Name,
		)
	}
	if err := setGlobalParams(); err != nil {
		return fmt.Errorf("could not set config params: %v", err)
	}
	st, err := generateGenesis(cliCtx.Context)
	if err != nil {
		return fmt.Errorf("could not generate genesis state: %v", err)
	}

	if outputJson != "" {
		if err := writeToOutputFile(outputJson, st, json.Marshal); err != nil {
			return err
		}
	}
	if outputYaml != "" {
		if err := writeToOutputFile(outputYaml, st, yaml.Marshal); err != nil {
			return err
		}
	}
	if outputSSZ != "" {
		type MinimumSSZMarshal interface {
			MarshalSSZ() ([]byte, error)
		}
		marshalFn := func(o interface{}) ([]byte, error) {
			marshaler, ok := o.(MinimumSSZMarshal)
			if !ok {
				return nil, errors.New("not a marshaler")
			}
			return marshaler.MarshalSSZ()
		}
		if err := writeToOutputFile(outputSSZ, st, marshalFn); err != nil {
			return err
		}
	}
	log.Info("Command completed")
	return nil
}

func setGlobalParams() error {
	chainConfigFile := generateGenesisStateFlags.ChainConfigFile
	if chainConfigFile != "" {
		log.Infof("Specified a chain config file: %s", chainConfigFile)
		return params.LoadChainConfigFile(chainConfigFile, nil)
	}
	cfg, err := params.ByName(generateGenesisStateFlags.ConfigName)
	if err != nil {
		return fmt.Errorf("unable to find config using name %s: %v", generateGenesisStateFlags.ConfigName, err)
	}
	return params.SetActive(cfg.Copy())
}

func generateGenesis(ctx context.Context) (state.BeaconState, error) {
	f := &generateGenesisStateFlags
	if f.GenesisTime == 0 {
		f.GenesisTime = uint64(time.Now().Unix())
		log.Info("No genesis time specified, defaulting to now()")
	}
	log.Infof("Delaying genesis %v by %v seconds", f.GenesisTime, f.GenesisTimeDelay)
	f.GenesisTime += f.GenesisTimeDelay
	log.Infof("Genesis is now %v", f.GenesisTime)

	v, err := version.FromString(f.ForkName)
	if err != nil {
		return nil, err
	}
	opts := make([]interop.PremineGenesisOpt, 0)
	nv := f.NumValidators
	if f.DepositJsonFile != "" {
		expanded, err := file.ExpandPath(f.DepositJsonFile)
		if err != nil {
			return nil, err
		}
		log.Printf("reading deposits from JSON at %s", expanded)
		b, err := os.ReadFile(expanded) // #nosec G304
		if err != nil {
			return nil, err
		}
		roots, dds, err := depositEntriesFromJSON(b)
		if err != nil {
			return nil, err
		}
		opts = append(opts, interop.WithDepositData(dds, roots))
	} else if nv == 0 {
		return nil, fmt.Errorf(
			"expected --num-validators > 0 or --deposit-json-file to have been provided",
		)
	}

	gen := &core.Genesis{}
	if f.GethGenesisJsonIn != "" {
		gbytes, err := os.ReadFile(f.GethGenesisJsonIn) // #nosec G304
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read %s", f.GethGenesisJsonIn)
		}
		if err := json.Unmarshal(gbytes, gen); err != nil {
			return nil, err
		}
		// set timestamps for genesis and shanghai fork
		gen.Timestamp = f.GenesisTime
		gen.Config.ShanghaiTime = interop.GethShanghaiTime(f.GenesisTime, params.BeaconConfig())
		//gen.Config.CancunTime = interop.GethCancunTime(f.GenesisTime, params.BeaconConfig())
		gen.Config.CancunTime = interop.GethCancunTime(f.GenesisTime, params.BeaconConfig())
		log.
			WithField("shanghai", fmt.Sprintf("%d", *gen.Config.ShanghaiTime)).
			WithField("cancun", fmt.Sprintf("%d", *gen.Config.CancunTime)).
			Info("setting fork geth times")
		if v > version.Altair {
			// set ttd to zero so EL goes post-merge immediately
			gen.Config.TerminalTotalDifficulty = big.NewInt(0)
			gen.Config.TerminalTotalDifficultyPassed = true
		}
	} else {
		gen = interop.GethTestnetGenesis(f.GenesisTime, params.BeaconConfig())
	}

	if f.GethGenesisJsonOut != "" {
		gbytes, err := json.MarshalIndent(gen, "", "\t")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(f.GethGenesisJsonOut, gbytes, os.ModePerm); err != nil {
			return nil, errors.Wrapf(err, "failed to write %s", f.GethGenesisJsonOut)
		}
	}

	gb := gen.ToBlock()

	// TODO: expose the PregenesisCreds option with a cli flag - for now defaulting to no withdrawal credentials at genesis
	log.Infof("Writing premined genesis with timestamp %d", f.GenesisTime)
	genesisState, err := interop.NewPreminedGenesis(ctx, f.GenesisTime, nv, 0, v, gb, opts...)
	if err != nil {
		return nil, err
	}
	log.Infof("beacon state genesis time %d", genesisState.GenesisTime())

	if f.OverrideEth1Data {
		log.Print("Overriding Eth1Data with data from execution client")
		conn, err := rpc.Dial(generateGenesisStateFlags.ExecutionEndpoint)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"could not dial %s please make sure you are running your execution client",
				generateGenesisStateFlags.ExecutionEndpoint)
		}
		client := ethclient.NewClient(conn)
		header, err := client.HeaderByNumber(ctx, big.NewInt(0))
		if err != nil {
			return nil, errors.Wrap(err, "could not get header by number")
		}
		t, err := trie.NewTrie(params.BeaconConfig().DepositContractTreeDepth)
		if err != nil {
			return nil, errors.Wrap(err, "could not create deposit tree")
		}
		depositRoot, err := t.HashTreeRoot()
		if err != nil {
			return nil, errors.Wrap(err, "could not get hash tree root")
		}
		e1d := &ethpb.Eth1Data{
			DepositRoot:  depositRoot[:],
			DepositCount: 0,
			BlockHash:    header.Hash().Bytes(),
		}
		if err := genesisState.SetEth1Data(e1d); err != nil {
			return nil, err
		}
		if err := genesisState.SetEth1DepositIndex(0); err != nil {
			return nil, err
		}
	}

	return genesisState, err
}

func depositEntriesFromJSON(enc []byte) ([][]byte, []*ethpb.Deposit_Data, error) {
	var depositJSON []*depositDataJSON
	if err := json.Unmarshal(enc, &depositJSON); err != nil {
		return nil, nil, err
	}
	dds := make([]*ethpb.Deposit_Data, len(depositJSON))
	roots := make([][]byte, len(depositJSON))
	for i, val := range depositJSON {
		root, data, err := depositJSONToDepositData(val)
		if err != nil {
			return nil, nil, err
		}
		dds[i] = data
		roots[i] = root
	}
	return roots, dds, nil
}

func depositJSONToDepositData(input *depositDataJSON) ([]byte, *ethpb.Deposit_Data, error) {
	root, err := hex.DecodeString(strings.TrimPrefix(input.DepositDataRoot, "0x"))
	if err != nil {
		return nil, nil, err
	}
	pk, err := hex.DecodeString(strings.TrimPrefix(input.PubKey, "0x"))
	if err != nil {
		return nil, nil, err
	}
	creds, err := hex.DecodeString(strings.TrimPrefix(input.WithdrawalCredentials, "0x"))
	if err != nil {
		return nil, nil, err
	}
	sig, err := hex.DecodeString(strings.TrimPrefix(input.Signature, "0x"))
	if err != nil {
		return nil, nil, err
	}
	return root, &ethpb.Deposit_Data{
		PublicKey:             pk,
		WithdrawalCredentials: creds,
		Amount:                input.Amount,
		Signature:             sig,
	}, nil
}

func writeToOutputFile(
	fPath string,
	data interface{},
	marshalFn func(o interface{}) ([]byte, error),
) error {
	encoded, err := marshalFn(data)
	if err != nil {
		return err
	}
	if err := file.WriteFile(fPath, encoded); err != nil {
		return err
	}
	log.Printf("Done writing genesis state to %s", fPath)
	return nil
}
