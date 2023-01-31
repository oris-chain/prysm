package depositsnapshot

import (
	"context"
	"math/big"
	"sort"
	"sync"

	fieldparams "github.com/prysmaticlabs/prysm/v3/config/fieldparams"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/wealdtech/go-bytesutil"
	"go.opencensus.io/trace"
)

type Cache struct {
	pendingDeposits   []*ethpb.DepositContainer
	deposits          []*ethpb.DepositContainer
	finalizedDeposits *FinalizedDeposits
	depositsByKey     map[[fieldparams.BLSPubkeyLength]byte][]*ethpb.DepositContainer
	depositsLock      sync.RWMutex
}

// FinalizedDeposits stores the trie of deposits that have been included
// in the beacon state up to the latest finalized checkpoint.
type FinalizedDeposits struct {
	Deposits        *depositTree
	MerkleTrieIndex int64
}

func (c *Cache) AllDeposits(ctx context.Context, untilBlk *big.Int) []*ethpb.Deposit {
	ctx, span := trace.StartSpan(ctx, "Cache.AllDeposits")
	defer span.End()
	c.depositsLock.RLock()
	defer c.depositsLock.RUnlock()

	return c.allDeposits(untilBlk)
}

func (c *Cache) allDeposits(untilBlk *big.Int) []*ethpb.Deposit {
	var deposits []*ethpb.Deposit
	for _, ctnr := range c.deposits {
		if untilBlk == nil || untilBlk.Uint64() >= ctnr.Eth1BlockHeight {
			deposits = append(deposits, ctnr.Deposit)
		}
	}
	return deposits
}

// AllDepositContainers returns all historical deposit containers.
func (c *Cache) AllDepositContainers(ctx context.Context) []*ethpb.DepositContainer {
	ctx, span := trace.StartSpan(ctx, "DepositsCache.AllDepositContainers")
	defer span.End()
	c.depositsLock.RLock()
	defer c.depositsLock.RUnlock()

	// Make a shallow copy of the deposits and return that. This way, the
	// caller can safely iterate over the returned list of deposits without
	// the possibility of new deposits showing up. If we were to return the
	// list without a copy, when a new deposit is added to the cache, it
	// would also be present in the returned value. This could result in a
	// race condition if the list is being iterated over.
	//
	// It's not necessary to make a deep copy of this list because the
	// deposits in the cache should never be modified. It is still possible
	// for the caller to modify one of the underlying deposits and modify
	// the cache, but that's not a race condition. Also, a deep copy would
	// take too long and use too much memory.
	deposits := make([]*ethpb.DepositContainer, len(c.deposits))
	copy(deposits, c.deposits)
	return deposits
}

func (c *Cache) DepositByPubkey(ctx context.Context, pubKey []byte) (*ethpb.Deposit, *big.Int) {
	ctx, span := trace.StartSpan(ctx, "DepositsCache.DepositByPubkey")
	defer span.End()
	c.depositsLock.RLock()
	defer c.depositsLock.RUnlock()

	var deposit *ethpb.Deposit
	var blockNum *big.Int
	deps, ok := c.depositsByKey[bytesutil.ToBytes48(pubKey)]
	if !ok || len(deps) == 0 {
		return deposit, blockNum
	}
	// We always return the first deposit if a particular
	// validator key has multiple deposits assigned to
	// it.
	deposit = deps[0].Deposit
	blockNum = big.NewInt(int64(deps[0].Eth1BlockHeight))
	return deposit, blockNum
}

func (c *Cache) DepositsNumberAndRootAtHeight(ctx context.Context, blockHeight *big.Int) (uint64, [32]byte) {
	ctx, span := trace.StartSpan(ctx, "DepositsCache.DepositsNumberAndRootAtHeight")
	defer span.End()
	c.depositsLock.RLock()
	defer c.depositsLock.RUnlock()
	heightIdx := sort.Search(len(c.deposits), func(i int) bool { return c.deposits[i].Eth1BlockHeight > blockHeight.Uint64() })
	// send the deposit root of the empty trie, if eth1follow distance is greater than the time of the earliest
	// deposit.
	if heightIdx == 0 {
		return 0, [32]byte{}
	}
	return uint64(heightIdx), bytesutil.ToBytes32(c.deposits[heightIdx-1].DepositRoot)
}

func (c *Cache) FinalizedDeposits(ctx context.Context) *FinalizedDeposits {
	ctx, span := trace.StartSpan(ctx, "DepositsCache.FinalizedDeposits")
	defer span.End()
	c.depositsLock.RLock()
	defer c.depositsLock.RUnlock()

	return &FinalizedDeposits{
		Deposits:        c.finalizedDeposits.Deposits,
		MerkleTrieIndex: c.finalizedDeposits.MerkleTrieIndex,
	}
}

func (c *Cache) NonFinalizedDeposits(ctx context.Context, lastFinalizedIndex int64, untilBlk *big.Int) []*ethpb.Deposit {
	ctx, span := trace.StartSpan(ctx, "DepositsCache.NonFinalizedDeposits")
	defer span.End()
	c.depositsLock.RLock()
	defer c.depositsLock.RUnlock()

	if c.finalizedDeposits == nil {
		return c.allDeposits(untilBlk)
	}

	var deposits []*ethpb.Deposit
	for _, d := range c.deposits {
		if (d.Index > lastFinalizedIndex) && (untilBlk == nil || untilBlk.Uint64() >= d.Eth1BlockHeight) {
			deposits = append(deposits, d.Deposit)
		}
	}

	return deposits
}
