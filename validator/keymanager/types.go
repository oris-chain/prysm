package keymanager

import (
	"context"
	"fmt"

	"github.com/prysmaticlabs/prysm/async/event"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"
)

// IKeymanager defines a struct which can be used to manage keys with important
// actions such as listing public keys, signing data, and subscribing to key changes.
type IKeymanager interface {
	KeysFetcher
	Signer
	Importer
	KeyChangeSubscriber
}

// KeysFetcher for validating private and public keys.
type KeysFetcher interface {
	FetchValidatingPrivateKeys(ctx context.Context) ([][32]byte, error)
	PublicKeysFetcher
}

// PublicKeysFetcher for validating public keys.
type PublicKeysFetcher interface {
	FetchValidatingPublicKeys(ctx context.Context) ([][48]byte, error)
}

// Signer allows signing messages using a validator private key.
type Signer interface {
	Sign(context.Context, *validatorpb.SignRequest) (bls.Signature, error)
}

// Importer can import new keystores into the keymanager.
type Importer interface {
	ImportKeystores(ctx context.Context, keystores []*Keystore, importsPassword string) error
}

// KeyChangeSubscriber allows subscribing to changes made to the underlying keys.
type KeyChangeSubscriber interface {
	SubscribeAccountChanges(pubKeysChan chan [][48]byte) event.Subscription
}

// Keystore json file representation as a Go struct.
type Keystore struct {
	Crypto  map[string]interface{} `json:"crypto"`
	ID      string                 `json:"uuid"`
	Pubkey  string                 `json:"pubkey"`
	Version uint                   `json:"version"`
	Name    string                 `json:"name"`
}

// Kind defines an enum for either imported, derived, or remote-signing
// keystores for Prysm wallets.
type Kind int

const (
	// Imported keymanager defines an on-disk, encrypted keystore-capable store.
	Imported Kind = iota
	// Derived keymanager using a hierarchical-deterministic algorithm.
	Derived
	// Remote keymanager capable of remote-signing data.
	Remote
)

// IncorrectPasswordErrMsg defines a common error string representing an EIP-2335
// keystore password was incorrect.
const IncorrectPasswordErrMsg = "invalid checksum"

// String marshals a keymanager kind to a string value.
func (k Kind) String() string {
	switch k {
	case Derived:
		return "derived"
	case Imported:
		return "direct"
	case Remote:
		return "remote"
	default:
		return fmt.Sprintf("%d", int(k))
	}
}

// ParseKind from a raw string, returning a keymanager kind.
func ParseKind(k string) (Kind, error) {
	switch k {
	case "derived":
		return Derived, nil
	case "direct":
		return Imported, nil
	case "remote":
		return Remote, nil
	default:
		return 0, fmt.Errorf("%s is not an allowed keymanager", k)
	}
}
