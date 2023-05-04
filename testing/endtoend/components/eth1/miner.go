package eth1

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/ethereum/go-ethereum/core"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	"github.com/prysmaticlabs/prysm/v4/runtime/interop"
	"github.com/prysmaticlabs/prysm/v4/testing/endtoend/helpers"
	e2e "github.com/prysmaticlabs/prysm/v4/testing/endtoend/params"
	e2etypes "github.com/prysmaticlabs/prysm/v4/testing/endtoend/types"
	log "github.com/sirupsen/logrus"
)

const (
	EthAddress = "0x878705ba3f8bc32fcf7f4caa1a35e72af65cf766"
)

// Miner represents an ETH1 node which mines blocks.
type Miner struct {
	e2etypes.ComponentRunner
	started      chan struct{}
	bootstrapEnr string
	enr          string
	cmd          *exec.Cmd
}

// NewMiner creates and returns an ETH1 node miner.
func NewMiner() *Miner {
	return &Miner{
		started: make(chan struct{}, 1),
	}
}

// ENR returns the miner's enode.
func (m *Miner) ENR() string {
	return m.enr
}

// SetBootstrapENR sets the bootstrap record.
func (m *Miner) SetBootstrapENR(bootstrapEnr string) {
	m.bootstrapEnr = bootstrapEnr
}

func (*Miner) DataDir(sub ...string) string {
	parts := append([]string{e2e.TestParams.TestPath, "eth1data/miner"}, sub...)
	return path.Join(parts...)
}

func (*Miner) Password() string {
	return KeystorePassword
}

func (m *Miner) initDataDir() error {
	eth1Path := m.DataDir()
	// Clear out potentially existing dir to prevent issues.
	if _, err := os.Stat(eth1Path); !os.IsNotExist(err) {
		if err = os.RemoveAll(eth1Path); err != nil {
			return err
		}
	}
	return nil
}

func (m *Miner) writeGenesis(binaryPath string, gen *core.Genesis) error {
	log.Infof("eth1 miner genesis timestamp=%d", e2e.TestParams.Eth1GenesisTime)
	b, err := json.Marshal(gen)
	if err != nil {
		return err
	}
	if err := file.WriteFile(genesisJsonPath(binaryPath), b); err != nil {
		return err
	}

	// write the same thing to the logs dir for debugging (included in outputs.zip)
	gethJsonLogPath := e2e.TestParams.Logfile("genesis.json")
	return file.WriteFile(gethJsonLogPath, b)
}

func genesisJsonPath(binaryPath string) string {
	return path.Join(path.Dir(binaryPath), "genesis.json")
}

func (m *Miner) initAttempt(ctx context.Context, binaryPath string, attempt int) (*os.File, error) {
	if err := m.initDataDir(); err != nil {
		return nil, err
	}
	initCmd := exec.CommandContext(ctx, binaryPath, "init", fmt.Sprintf("--datadir=%s", m.DataDir()), genesisJsonPath(binaryPath)) // #nosec G204 -- Safe

	// redirect stderr to a log file
	initFile, err := helpers.DeleteAndCreatePath(e2e.TestParams.Logfile("eth1-init_miner.log"))
	if err != nil {
		return nil, err
	}
	initCmd.Stderr = initFile

	// run init command and wait until it exits. this will initialize the geth node (required before starting).
	if err = initCmd.Start(); err != nil {
		return nil, err
	}
	if err = initCmd.Wait(); err != nil {
		return nil, err
	}

	pwFile := m.DataDir("keystore", minerPasswordFile)
	args := []string{
		"--nat=none", // disable nat traversal in e2e, it is failure prone and not needed
		fmt.Sprintf("--datadir=%s", m.DataDir()),
		fmt.Sprintf("--http.port=%d", e2e.TestParams.Ports.Eth1RPCPort),
		fmt.Sprintf("--ws.port=%d", e2e.TestParams.Ports.Eth1WSPort),
		fmt.Sprintf("--authrpc.port=%d", e2e.TestParams.Ports.Eth1AuthRPCPort),
		fmt.Sprintf("--bootnodes=%s", m.bootstrapEnr),
		fmt.Sprintf("--port=%d", e2e.TestParams.Ports.Eth1Port),
		fmt.Sprintf("--networkid=%d", NetworkId),
		"--http",
		"--http.api=engine,net,eth",
		"--http.addr=127.0.0.1",
		"--http.corsdomain=\"*\"",
		"--http.vhosts=\"*\"",
		"--rpc.allow-unprotected-txs",
		"--ws",
		"--ws.api=net,eth,engine",
		"--ws.addr=127.0.0.1",
		"--ws.origins=\"*\"",
		"--ipcdisable",
		"--verbosity=4",
		"--mine",
		fmt.Sprintf("--unlock=%s", EthAddress),
		"--allow-insecure-unlock",
		"--syncmode=full",
		fmt.Sprintf("--miner.etherbase=%s", EthAddress),
		fmt.Sprintf("--txpool.locals=%s", EthAddress),
		fmt.Sprintf("--password=%s", pwFile),
	}

	keystorePath, err := e2e.TestParams.Paths.MinerKeyPath()
	if err != nil {
		return nil, err
	}
	if err = file.CopyFile(keystorePath, m.DataDir("keystore", minerFile)); err != nil {
		return nil, errors.Wrapf(err, "error copying %s to %s", keystorePath, m.DataDir("keystore", minerFile))
	}
	err = file.WriteFile(pwFile, []byte(KeystorePassword))
	if err != nil {
		return nil, err
	}

	runCmd := exec.CommandContext(ctx, binaryPath, args...) // #nosec G204 -- Safe
	// redirect miner stderr to a log file
	minerLog, err := helpers.DeleteAndCreatePath(e2e.TestParams.Logfile("eth1_miner.log"))
	if err != nil {
		return nil, err
	}
	runCmd.Stderr = minerLog
	log.Infof("Starting eth1 miner, attempt %d, with flags: %s", attempt, strings.Join(args[2:], " "))
	if err = runCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start eth1 chain: %w", err)
	}
	if err = helpers.WaitForTextInFile(minerLog, "Started P2P networking"); err != nil {
		kerr := runCmd.Process.Kill()
		if kerr != nil {
			log.WithError(kerr).Error("error sending kill to failed miner command process")
		}
		return nil, fmt.Errorf("P2P log not found, this means the eth1 chain had issues starting: %w", err)
	}
	m.cmd = runCmd
	return minerLog, nil
}

// Start runs a mining ETH1 node.
// The miner is responsible for moving the ETH1 chain forward and for deploying the deposit contract.
func (m *Miner) Start(ctx context.Context) error {

	// find geth so we can run it.
	binaryPath, found := bazel.FindBinary("cmd/geth", "geth")
	if !found {
		return errors.New("go-ethereum binary not found")
	}

	gen := interop.GethTestnetGenesis(e2e.TestParams.Eth1GenesisTime, params.BeaconConfig())
	if err := m.writeGenesis(binaryPath, gen); err != nil {
		return errors.Wrapf(err, "unable to write genesis.json file to geth binary path %s", binaryPath)
	}

	// give the miner start a couple of tries, since the p2p networking check is flaky
	var retryErr error
	var minerLog *os.File
	for attempt := 0; attempt < 3; attempt++ {
		minerLog, retryErr = m.initAttempt(ctx, binaryPath, attempt)
		if retryErr == nil {
			log.Infof("miner started after %d retries", attempt)
			break
		}
	}
	if retryErr != nil {
		return retryErr
	}

	enode, err := enodeFromLogFile(minerLog.Name())
	if err != nil {
		return err
	}
	enode = "enode://" + enode + "@127.0.0.1:" + fmt.Sprintf("%d", e2e.TestParams.Ports.Eth1Port)
	m.enr = enode
	log.Infof("Communicated enode. Enode is %s", enode)

	block := gen.ToBlock()
	eth1BlockHash := block.Hash()
	e2e.TestParams.Eth1GenesisBlock = block
	log.Infof("genesis block timestamp=%d, hash=%#x", block.Time(), eth1BlockHash)

	// Mark node as ready.
	close(m.started)
	return m.cmd.Wait()
}

// Started checks whether ETH1 node is started and ready to be queried.
func (m *Miner) Started() <-chan struct{} {
	return m.started
}

// Pause pauses the component and its underlying process.
func (m *Miner) Pause() error {
	return m.cmd.Process.Signal(syscall.SIGSTOP)
}

// Resume resumes the component and its underlying process.
func (m *Miner) Resume() error {
	return m.cmd.Process.Signal(syscall.SIGCONT)
}

// Stop kills the component and its underlying process.
func (m *Miner) Stop() error {
	return m.cmd.Process.Kill()
}

func enodeFromLogFile(name string) (string, error) {
	byteContent, err := os.ReadFile(name) // #nosec G304
	if err != nil {
		return "", err
	}
	contents := string(byteContent)

	searchText := "self=enode://"
	startIdx := strings.Index(contents, searchText)
	if startIdx == -1 {
		return "", fmt.Errorf("did not find ENR text in %s", contents)
	}
	startIdx += len(searchText)
	endIdx := strings.Index(contents[startIdx:], "@")
	if endIdx == -1 {
		return "", fmt.Errorf("did not find ENR text in %s", contents)
	}
	enode := contents[startIdx : startIdx+endIdx]
	return strings.TrimPrefix(enode, "-"), nil
}
