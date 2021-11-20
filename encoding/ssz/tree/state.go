package tree

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
)

var (
	BLSPubkeyType = view.BasicVectorType(view.ByteType, 48)
	ValidatorType = view.ContainerType("Validator", []view.FieldDef{
		{"pubkey", BLSPubkeyType},
		{"withdrawal_credentials", view.RootType},
		{"effective_balance", view.Uint64Type}, // Balance at stake
		{"slashed", view.BoolType},
		{"activation_eligibility_epoch", view.Uint64Type}, // When criteria for activation were met
		{"activation_epoch", view.Uint64Type},
		{"exit_epoch", view.Uint64Type},
		{"withdrawable_epoch", view.Uint64Type}, // When validator can withdraw funds
	})
	ForkType = view.ContainerType("Fork", []view.FieldDef{
		{"previous_version", view.Bytes4Type},
		{"current_version", view.Bytes4Type},
		{"epoch", view.Uint64Type},
	})
	BeaconBlockHeaderType = view.ContainerType("BeaconBlockHeader", []view.FieldDef{
		{"slot", view.Uint64Type},
		{"proposer_index", view.Uint64Type},
		{"parent_root", view.RootType},
		{"state_root", view.RootType},
		{"body_root", view.RootType},
	})
	BlockRootsType      = view.VectorType(view.RootType, 8192)
	StateRootsType      = view.VectorType(view.RootType, 8192)
	HistoricalRootsType = view.ListType(view.RootType, 16777216)
	Eth1DataType        = view.ContainerType("Eth1Data", []view.FieldDef{
		{"deposit_root", view.RootType},
		{"deposit_count", view.Uint64Type},
		{"block_hash", view.RootType},
	})
	Eth1DataVotesType     = view.ComplexListType(Eth1DataType, 2048)
	ValidatorsType        = view.ComplexListType(ValidatorType, 1099511627776)
	BalancesType          = view.BasicListType(view.Uint64Type, 1099511627776)
	RandaoMixesType       = view.VectorType(view.RootType, 65536)
	SlashingsType         = view.BasicVectorType(view.Uint64Type, 8192)
	ParticipationType     = view.BasicListType(view.ByteType, 1099511627776)
	JustificationBitsType = view.BitVectorType(4)
	CheckpointType        = view.ContainerType("Checkpoint", []view.FieldDef{
		{"epoch", view.Uint64Type},
		{"root", view.RootType},
	})
	InactivityScoresType  = view.BasicListType(view.Uint64Type, 1099511627776)
	SyncCommitteeKeysType = view.VectorType(BLSPubkeyType, 512)
	SyncCommitteeType     = view.ContainerType("SyncCommittee", []view.FieldDef{
		{"pubkeys", SyncCommitteeKeysType},
		{"aggregate_pubkey", BLSPubkeyType},
	})
	BeaconStateAltairType = view.ContainerType("BeaconStateAltair", []view.FieldDef{
		{"genesis_time", view.Uint64Type},
		{"genesis_validators_root", view.RootType},
		{"slot", view.Uint64Type},
		{"fork", ForkType},
		{"latest_block_header", BeaconBlockHeaderType},
		{"block_roots", BlockRootsType},
		{"state_roots", StateRootsType},
		{"historical_roots", HistoricalRootsType},
		{"eth1_data", Eth1DataType},
		{"eth1_data_votes", Eth1DataVotesType},
		{"eth1_deposit_index", view.Uint64Type},
		{"validators", ValidatorsType},
		{"balances", BalancesType},
		{"randao_mixes", RandaoMixesType},
		{"slashings", SlashingsType},
		{"previous_epoch_participation", ParticipationType},
		{"current_epoch_participation", ParticipationType},
		{"justification_bits", JustificationBitsType},
		{"previous_justified_checkpoint", CheckpointType},
		{"current_justified_checkpoint", CheckpointType},
		{"finalized_checkpoint", CheckpointType},
		{"inactivity_scores", InactivityScoresType},
		{"current_sync_committee", SyncCommitteeType},
		{"next_sync_committee", SyncCommitteeType},
	})
)

func VerifyProof(root [32]byte, proof [][]byte, leaf tree.Root, generalizedIndex tree.Gindex64) bool {
	h := leaf
	hFn := tree.GetHashFn()
	idx := generalizedIndex
	fmt.Println("")
	for _, elem := range proof {
		if idx%2 == 0 {
			fmt.Printf("Combi (H, elem) %s and %#x\n", h.String(), bytesutil.ToBytes32(elem))
			h = hFn(h, bytesutil.ToBytes32(elem))
		} else {
			fmt.Printf("Combi (elem, H) %#x and %s\n", bytesutil.ToBytes32(elem), h.String())
			h = hFn(bytesutil.ToBytes32(elem), h)
		}
		fmt.Printf("Got %s\n", h.String())
		idx = idx / 2
	}
	return h == root
}

func Proof(
	state view.View,
	index uint64,
) (proof [][]byte, generalizedIdx tree.Gindex64, err error) {
	cont, ok := state.(*view.ContainerView)
	if !ok {
		err = errors.New("not a container")
		return
	}
	depth := tree.CoverDepth(cont.FieldCount())
	generalizedIdx, err = tree.ToGindex64(index, depth)
	if err != nil {
		return
	}
	fmt.Println(depth, generalizedIdx)
	leaves := make(map[tree.Gindex64]struct{})
	leaves[generalizedIdx] = struct{}{}
	leavesSorted := make([]tree.Gindex64, 0, len(leaves))
	for g := range leaves {
		leavesSorted = append(leavesSorted, g)
	}
	sort.Slice(leavesSorted, func(i, j int) bool {
		return leavesSorted[i] < leavesSorted[j]
	})

	// Mark every gindex that is between the root and the leaves.
	interest := make(map[tree.Gindex64]struct{})
	for _, g := range leavesSorted {
		iter, _ := g.BitIter()
		n := tree.Gindex64(1)
		for {
			right, ok := iter.Next()
			if !ok {
				break
			}
			n *= 2
			if right {
				n += 1
			}
			interest[n] = struct{}{}
		}
	}
	witness := make(map[tree.Gindex64]struct{})
	// For every gindex that is covered, check if the sibling is covered, and if not, it's a witness
	for g := range interest {
		if _, ok := interest[g^1]; !ok {
			witness[g^1] = struct{}{}
		}
	}
	witnessSorted := make([]tree.Gindex64, 0, len(witness))
	for g := range witness {
		witnessSorted = append(witnessSorted, g)
	}
	sort.Slice(witnessSorted, func(i, j int) bool {
		return witnessSorted[i] < witnessSorted[j]
	})

	node := state.Backing()
	hFn := tree.GetHashFn()
	root := state.HashTreeRoot(hFn)
	fmt.Printf("Root %s\n", root.String())
	proof = make([][]byte, 0, len(witnessSorted))
	leaf, err2 := node.Getter(generalizedIdx)
	if err2 != nil {
		err = err2
		return
	}
	topRoot := node.MerkleRoot(hFn)
	fmt.Printf("Top root %s\n", topRoot.String())
	curr := leaf.MerkleRoot(hFn)
	fmt.Printf("%s\n", curr.String())
	for i := len(witnessSorted) - 1; i > 0; i-- {
		g := witnessSorted[i]
		n, err2 := node.Getter(g)
		if err2 != nil {
			err = err2
			return
		}
		root := n.MerkleRoot(hFn)
		fmt.Printf("sibling idx %d\n", g^1)
		sibling, err3 := node.Getter(g ^ 1)
		if err3 != nil {
			err = err3
			return
		}
		siblingRoot := sibling.MerkleRoot(hFn)
		fmt.Printf("%d witness index with proof (left %s, right %s)\n", g, siblingRoot.String(), root.String())
		fmt.Printf("Got %s\n", hFn(siblingRoot, root).String())
		proof = append(proof, root[:])
	}
	return
}
