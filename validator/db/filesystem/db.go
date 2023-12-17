package filesystem

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	validatorpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	slashingProtectionDirName = "slashing-protection"
	configurationFileName     = "configuration.yaml"

	DatabaseDirName = "validatorClientData"
)

type (
	// Store is a filesystem implementation of the validator client database.
	Store struct {
		configurationMu    sync.RWMutex
		pkToSlashingMu     map[[fieldparams.BLSPubkeyLength]byte]*sync.RWMutex
		slashingMuMapMu    sync.Mutex
		databaseParentPath string
		databasePath       string
	}

	// Graffiti contains the graffiti information.
	Graffiti struct {
		// In BoltDB implementation, calling GraffitiOrderedIndex with
		// the filehash stored in DB, but without an OrderedIndex already
		// stored in DB returns 0.
		// ==> Using the default value of uint64 is OK.
		OrderedIndex uint64
		FileHash     *string
	}

	// Configuration contains the genesis information, the proposer settings and the graffiti.
	Configuration struct {
		GenesisValidatorsRoot *string                              `yaml:"genesisValidatorsRoot,omitempty"`
		ProposerSettings      *validatorpb.ProposerSettingsPayload `yaml:"proposerSettings,omitempty"`
		Graffiti              *Graffiti                            `yaml:"graffiti,omitempty"`
	}

	// ValidatorSlashingProtection contains the latest signed block slot, the last signed attestation.
	// It is used to protect against validator slashing, implementing the EIP-3076 minimal slashing protection database.
	// https://eips.ethereum.org/EIPS/eip-3076
	ValidatorSlashingProtection struct {
		LatestSignedBlockSlot            *uint64 `yaml:"latestSignedBlockSlot,omitempty"`
		LastSignedAttestationSourceEpoch uint64  `yaml:"lastSignedAttestationSourceEpoch"`
		LastSignedAttestationTargetEpoch *uint64 `yaml:"lastSignedAttestationTargetEpoch,omitempty"`
	}

	// Config represents store's config object.
	Config struct {
		PubKeys [][fieldparams.BLSPubkeyLength]byte
	}
)

var log = logrus.WithField("prefix", "db")

// NewStore creates a new filesystem store.
func NewStore(databaseParentPath string, config *Config) (*Store, error) {
	s := &Store{
		databaseParentPath: databaseParentPath,
		databasePath:       path.Join(databaseParentPath, DatabaseDirName),
		pkToSlashingMu:     make(map[[fieldparams.BLSPubkeyLength]byte]*sync.RWMutex),
	}

	// Initialize the required public keys into the DB to ensure they're not empty.
	if config != nil {
		if err := s.UpdatePublicKeysBuckets(config.PubKeys); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Close only exists to satisfy the interface.
func (*Store) Close() error {
	return nil
}

// DatabasePath returns the path at which this database writes files.
func (s *Store) DatabasePath() string {
	// The returned path is actually the parent path, to be consistent with the BoltDB implementation.
	return s.databaseParentPath
}

// ClearDB removes any previously stored data at the configured data directory.
func (s *Store) ClearDB() error {
	if err := os.RemoveAll(s.databasePath); err != nil {
		return errors.Wrapf(err, "cannot remove database at path %s", s.databasePath)
	}

	return nil
}

// UpdatePublicKeysBuckets creates a file for each public key in the database directory if needed.
func (s *Store) UpdatePublicKeysBuckets(pubKeys [][fieldparams.BLSPubkeyLength]byte) error {
	validatorSlashingProtection := ValidatorSlashingProtection{}

	// Marshal the ValidatorSlashingProtection struct.
	yfile, err := yaml.Marshal(validatorSlashingProtection)
	if err != nil {
		return errors.Wrap(err, "could not marshal validator slashing protection")
	}

	// Create the directory if needed.
	slashingProtectionDirPath := s.slashingProtectionDirPath()
	if err := file.MkdirAll(slashingProtectionDirPath); err != nil {
		return errors.Wrapf(err, "could not create directory %s", s.databasePath)
	}

	for _, pubKey := range pubKeys {
		// Get the file path for the public key.
		path := s.pubkeySlashingProtectionFilePath(pubKey)

		// Check if the public key has a file in the database.
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		// Write the ValidatorSlashingProtection struct to the file.
		if err := file.WriteFile(path, yfile); err != nil {
			return errors.Wrapf(err, "could not write into %s.yaml", path)
		}
	}

	return nil
}

// slashingProtectionDirPath returns the path of the slashing protection directory.
func (s *Store) slashingProtectionDirPath() string {
	return path.Join(s.databasePath, slashingProtectionDirName)
}

// pubkeySlashingProtectionFilePath returns the path of the slashing protection file for a public key.
func (s *Store) pubkeySlashingProtectionFilePath(pubKey [fieldparams.BLSPubkeyLength]byte) string {
	slashingProtectionDirPath := s.slashingProtectionDirPath()
	pubkeyFileName := fmt.Sprintf("%s.yaml", hexutil.Encode(pubKey[:]))

	return path.Join(slashingProtectionDirPath, pubkeyFileName)
}

// configurationFilePath returns the path of the configuration file.
func (s *Store) configurationFilePath() string {
	return path.Join(s.databasePath, configurationFileName)
}

// configuration returns the configuration.
func (s *Store) configuration() (*Configuration, error) {
	config := &Configuration{}

	// Get the path of config file.
	configFilePath := s.configurationFilePath()
	cleanedConfigFilePath := filepath.Clean(configFilePath)

	// Read lock the mutex.
	s.configurationMu.RLock()
	defer s.configurationMu.RUnlock()

	// Check if config file exists.
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil
	}

	// Read the config file.
	yfile, err := os.ReadFile(cleanedConfigFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read %s", cleanedConfigFilePath)
	}

	// Unmarshal the config file into Config struct.
	if err := yaml.Unmarshal(yfile, &config); err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal %s", cleanedConfigFilePath)
	}

	return config, nil
}

// saveConfiguration saves the configuration.
func (s *Store) saveConfiguration(config *Configuration) error {
	// If config is nil, return
	if config == nil {
		return nil
	}

	// Create the directory if needed.
	if err := file.MkdirAll(s.databasePath); err != nil {
		return errors.Wrapf(err, "could not create directory %s", s.databasePath)
	}

	// Get the path of config file.
	configFilePath := s.configurationFilePath()

	// Marshal config into yaml.
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "could not marshal config.yaml")
	}

	// Write lock the mutex.
	s.configurationMu.Lock()
	defer s.configurationMu.Unlock()

	// Write the data to config.yaml.
	if err := file.WriteFile(configFilePath, data); err != nil {
		return errors.Wrap(err, "could not write genesis info into config.yaml")
	}

	return nil
}

// validatorSlashingProtection returns the slashing protection for a public key.
func (s *Store) validatorSlashingProtection(publicKey [fieldparams.BLSPubkeyLength]byte) (*ValidatorSlashingProtection, error) {
	var mu *sync.RWMutex
	validatorSlashingProtection := &ValidatorSlashingProtection{}

	// Get the slashing protection file path.
	path := s.pubkeySlashingProtectionFilePath(publicKey)
	cleanedPath := filepath.Clean(path)

	// Check if the public key has a file in the database.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	// Lock the mutex protecting the map of public keys to slashing protection mutexes.
	s.slashingMuMapMu.Lock()

	// Get / create the mutex for the public key.
	mu, ok := s.pkToSlashingMu[publicKey]
	if !ok {
		mu = &sync.RWMutex{}
		s.pkToSlashingMu[publicKey] = mu
	}

	// Release the mutex protecting the map of public keys to slashing protection mutexes.
	s.slashingMuMapMu.Unlock()

	// Read lock the mutex for the public key.
	mu.RLock()
	defer mu.RUnlock()

	// Read the file and unmarshal it into ValidatorSlashingProtection struct.
	yfile, err := os.ReadFile(cleanedPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read %s", cleanedPath)
	}

	if err := yaml.Unmarshal(yfile, validatorSlashingProtection); err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal %s", cleanedPath)
	}

	return validatorSlashingProtection, nil
}

// saveValidatorSlashingProtection saves the slashing protection for a public key.
func (s *Store) saveValidatorSlashingProtection(
	publicKey [fieldparams.BLSPubkeyLength]byte,
	validatorSlashingProtection *ValidatorSlashingProtection,
) error {
	// If the ValidatorSlashingProtection struct is nil, return.
	if validatorSlashingProtection == nil {
		return nil
	}

	// Create the directory if needed.
	slashingProtectionDirPath := s.slashingProtectionDirPath()
	if err := file.MkdirAll(slashingProtectionDirPath); err != nil {
		return errors.Wrapf(err, "could not create directory %s", s.databasePath)
	}

	// Get the file path for the public key.
	path := s.pubkeySlashingProtectionFilePath(publicKey)

	// Lock the mutex protecting the map of public keys to slashing protection mutexes.
	s.slashingMuMapMu.Lock()

	// Get / create the mutex for the public key.
	mu, ok := s.pkToSlashingMu[publicKey]
	if !ok {
		mu = &sync.RWMutex{}
		s.pkToSlashingMu[publicKey] = mu
	}

	// Release the mutex protecting the map of public keys to slashing protection mutexes.
	s.slashingMuMapMu.Unlock()

	// Write lock the mutex.
	mu.Lock()
	defer mu.Unlock()

	// Marshal the ValidatorSlashingProtection struct.
	yfile, err := yaml.Marshal(validatorSlashingProtection)
	if err != nil {
		return errors.Wrap(err, "could not marshal validator slashing protection")
	}

	// Write the ValidatorSlashingProtection struct to the file.
	if err := file.WriteFile(path, yfile); err != nil {
		return errors.Wrapf(err, "could not write into %s.yaml", path)
	}

	return nil
}

// publicKeys returns the public keys existing in the database directory.
func (s *Store) publicKeys() ([][fieldparams.BLSPubkeyLength]byte, error) {
	// Get the slashing protection directory path.
	slashingProtectionDirPath := s.slashingProtectionDirPath()

	// If the slashing protection directory does not exist, return an empty slice.
	if _, err := os.Stat(slashingProtectionDirPath); os.IsNotExist(err) {
		return nil, nil
	}

	// Get all entries in the slashing protection directory.
	entries, err := os.ReadDir(slashingProtectionDirPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read database directory")
	}

	// Collect public keys.
	publicKeys := make([][fieldparams.BLSPubkeyLength]byte, 0, len(entries))
	for _, entry := range entries {
		if !(entry.Type().IsRegular() && strings.HasPrefix(entry.Name(), "0x")) {
			continue
		}

		// Convert the file name to a public key.
		publicKeyHex := strings.TrimSuffix(entry.Name(), ".yaml")
		publicKeyBytes, err := hexutil.Decode(publicKeyHex)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode %s", publicKeyHex)
		}

		publicKey := [fieldparams.BLSPubkeyLength]byte{}
		copy(publicKey[:], publicKeyBytes)

		publicKeys = append(publicKeys, publicKey)
	}

	return publicKeys, nil
}
