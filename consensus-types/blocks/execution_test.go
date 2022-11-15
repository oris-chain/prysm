package blocks_test

import (
	"testing"

	fieldparams "github.com/prysmaticlabs/prysm/v3/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v3/crypto/hash"
	"github.com/prysmaticlabs/prysm/v3/encoding/ssz"
	enginev1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v3/testing/assert"
	"github.com/prysmaticlabs/prysm/v3/testing/require"
)

func TestWrapExecutionPayload(t *testing.T) {
	data := &enginev1.ExecutionPayload{GasUsed: 54}
	wsb, err := blocks.WrappedExecutionPayload(data)
	require.NoError(t, err)

	assert.DeepEqual(t, data, wsb.Proto())
}

func TestWrapExecutionPayloadHeader(t *testing.T) {
	data := &enginev1.ExecutionPayloadHeader{GasUsed: 54}
	wsb, err := blocks.WrappedExecutionPayloadHeader(data)
	require.NoError(t, err)

	assert.DeepEqual(t, data, wsb.Proto())
}

func TestWrapExecutionPayload_IsNil(t *testing.T) {
	_, err := blocks.WrappedExecutionPayload(nil)
	require.Equal(t, blocks.ErrNilObjectWrapped, err)

	data := &enginev1.ExecutionPayload{GasUsed: 54}
	wsb, err := blocks.WrappedExecutionPayload(data)
	require.NoError(t, err)

	assert.Equal(t, false, wsb.IsNil())
}

func TestWrapExecutionPayloadHeader_IsNil(t *testing.T) {
	_, err := blocks.WrappedExecutionPayloadHeader(nil)
	require.Equal(t, blocks.ErrNilObjectWrapped, err)

	data := &enginev1.ExecutionPayloadHeader{GasUsed: 54}
	wsb, err := blocks.WrappedExecutionPayloadHeader(data)
	require.NoError(t, err)

	assert.Equal(t, false, wsb.IsNil())
}

func TestWrapExecutionPayload_SSZ(t *testing.T) {
	wsb := createWrappedPayload(t)
	rt, err := wsb.HashTreeRoot()
	assert.NoError(t, err)
	assert.NotEmpty(t, rt)

	var b []byte
	b, err = wsb.MarshalSSZTo(b)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(b))
	encoded, err := wsb.MarshalSSZ()
	require.NoError(t, err)
	assert.NotEqual(t, 0, wsb.SizeSSZ())
	assert.NoError(t, wsb.UnmarshalSSZ(encoded))
}

func TestWrapExecutionPayloadHeader_SSZ(t *testing.T) {
	wsb := createWrappedPayloadHeader(t)
	rt, err := wsb.HashTreeRoot()
	assert.NoError(t, err)
	assert.NotEmpty(t, rt)

	var b []byte
	b, err = wsb.MarshalSSZTo(b)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(b))
	encoded, err := wsb.MarshalSSZ()
	require.NoError(t, err)
	assert.NotEqual(t, 0, wsb.SizeSSZ())
	assert.NoError(t, wsb.UnmarshalSSZ(encoded))
}

func TestWrapExecutionPayloadCapella(t *testing.T) {
	data := &enginev1.ExecutionPayloadCapella{
		ParentHash:    []byte("parenthash"),
		FeeRecipient:  []byte("feerecipient"),
		StateRoot:     []byte("stateroot"),
		ReceiptsRoot:  []byte("receiptsroot"),
		LogsBloom:     []byte("logsbloom"),
		PrevRandao:    []byte("prevrandao"),
		BlockNumber:   11,
		GasLimit:      22,
		GasUsed:       33,
		Timestamp:     44,
		ExtraData:     []byte("extradata"),
		BaseFeePerGas: []byte("basefeepergas"),
		BlockHash:     []byte("blockhash"),
		Transactions:  [][]byte{[]byte("transaction")},
		Withdrawals: []*enginev1.Withdrawal{{
			WithdrawalIndex:  55,
			ValidatorIndex:   66,
			ExecutionAddress: []byte("executionaddress"),
			Amount:           77,
		}},
	}
	payload, err := blocks.WrappedExecutionPayloadCapella(data)
	require.NoError(t, err)

	assert.DeepEqual(t, data, payload.Proto())
}

func TestWrapExecutionPayloadHeaderCapella(t *testing.T) {
	data := &enginev1.ExecutionPayloadHeaderCapella{
		ParentHash:       []byte("parenthash"),
		FeeRecipient:     []byte("feerecipient"),
		StateRoot:        []byte("stateroot"),
		ReceiptsRoot:     []byte("receiptsroot"),
		LogsBloom:        []byte("logsbloom"),
		PrevRandao:       []byte("prevrandao"),
		BlockNumber:      11,
		GasLimit:         22,
		GasUsed:          33,
		Timestamp:        44,
		ExtraData:        []byte("extradata"),
		BaseFeePerGas:    []byte("basefeepergas"),
		BlockHash:        []byte("blockhash"),
		TransactionsRoot: []byte("transactionsroot"),
		WithdrawalsRoot:  []byte("withdrawalsroot"),
	}
	payload, err := blocks.WrappedExecutionPayloadHeaderCapella(data)
	require.NoError(t, err)

	assert.DeepEqual(t, data, payload.Proto())

	txRoot, err := payload.TransactionsRoot()
	require.NoError(t, err)
	require.DeepEqual(t, txRoot, data.TransactionsRoot)

	wrRoot, err := payload.WithdrawalsRoot()
	require.NoError(t, err)
	require.DeepEqual(t, wrRoot, data.WithdrawalsRoot)
}

func TestWrapExecutionPayloadCapella_IsNil(t *testing.T) {
	_, err := blocks.WrappedExecutionPayloadCapella(nil)
	require.Equal(t, blocks.ErrNilObjectWrapped, err)

	data := &enginev1.ExecutionPayloadCapella{GasUsed: 54}
	payload, err := blocks.WrappedExecutionPayloadCapella(data)
	require.NoError(t, err)

	assert.Equal(t, false, payload.IsNil())
}

func TestWrapExecutionPayloadHeaderCapella_IsNil(t *testing.T) {
	_, err := blocks.WrappedExecutionPayloadHeaderCapella(nil)
	require.Equal(t, blocks.ErrNilObjectWrapped, err)

	data := &enginev1.ExecutionPayloadHeaderCapella{GasUsed: 54}
	payload, err := blocks.WrappedExecutionPayloadHeaderCapella(data)
	require.NoError(t, err)

	assert.Equal(t, false, payload.IsNil())
}

func TestWrapExecutionPayloadCapella_SSZ(t *testing.T) {
	payload := createWrappedPayloadCapella(t)
	rt, err := payload.HashTreeRoot()
	assert.NoError(t, err)
	assert.NotEmpty(t, rt)

	var b []byte
	b, err = payload.MarshalSSZTo(b)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(b))
	encoded, err := payload.MarshalSSZ()
	require.NoError(t, err)
	assert.NotEqual(t, 0, payload.SizeSSZ())
	assert.NoError(t, payload.UnmarshalSSZ(encoded))
}

func TestWrapExecutionPayloadHeaderCapella_SSZ(t *testing.T) {
	payload := createWrappedPayloadHeaderCapella(t)
	rt, err := payload.HashTreeRoot()
	assert.NoError(t, err)
	assert.NotEmpty(t, rt)

	var b []byte
	b, err = payload.MarshalSSZTo(b)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(b))
	encoded, err := payload.MarshalSSZ()
	require.NoError(t, err)
	assert.NotEqual(t, 0, payload.SizeSSZ())
	assert.NoError(t, payload.UnmarshalSSZ(encoded))
}

func createWrappedPayload(t testing.TB) interfaces.ExecutionData {
	wsb, err := blocks.WrappedExecutionPayload(&enginev1.ExecutionPayload{
		ParentHash:    make([]byte, fieldparams.RootLength),
		FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:     make([]byte, fieldparams.RootLength),
		ReceiptsRoot:  make([]byte, fieldparams.RootLength),
		LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:    make([]byte, fieldparams.RootLength),
		BlockNumber:   0,
		GasLimit:      0,
		GasUsed:       0,
		Timestamp:     0,
		ExtraData:     make([]byte, 0),
		BaseFeePerGas: make([]byte, fieldparams.RootLength),
		BlockHash:     make([]byte, fieldparams.RootLength),
		Transactions:  make([][]byte, 0),
	})
	require.NoError(t, err)
	return wsb
}

func createWrappedPayloadHeader(t testing.TB) interfaces.ExecutionData {
	wsb, err := blocks.WrappedExecutionPayloadHeader(&enginev1.ExecutionPayloadHeader{
		ParentHash:       make([]byte, fieldparams.RootLength),
		FeeRecipient:     make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:        make([]byte, fieldparams.RootLength),
		ReceiptsRoot:     make([]byte, fieldparams.RootLength),
		LogsBloom:        make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:       make([]byte, fieldparams.RootLength),
		BlockNumber:      0,
		GasLimit:         0,
		GasUsed:          0,
		Timestamp:        0,
		ExtraData:        make([]byte, 0),
		BaseFeePerGas:    make([]byte, fieldparams.RootLength),
		BlockHash:        make([]byte, fieldparams.RootLength),
		TransactionsRoot: make([]byte, fieldparams.RootLength),
	})
	require.NoError(t, err)
	return wsb
}

func createWrappedPayloadCapella(t testing.TB) interfaces.ExecutionData {
	payload, err := blocks.WrappedExecutionPayloadCapella(&enginev1.ExecutionPayloadCapella{
		ParentHash:    make([]byte, fieldparams.RootLength),
		FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:     make([]byte, fieldparams.RootLength),
		ReceiptsRoot:  make([]byte, fieldparams.RootLength),
		LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:    make([]byte, fieldparams.RootLength),
		BlockNumber:   0,
		GasLimit:      0,
		GasUsed:       0,
		Timestamp:     0,
		ExtraData:     make([]byte, 0),
		BaseFeePerGas: make([]byte, fieldparams.RootLength),
		BlockHash:     make([]byte, fieldparams.RootLength),
		Transactions:  make([][]byte, 0),
		Withdrawals:   make([]*enginev1.Withdrawal, 0),
	})
	require.NoError(t, err)
	return payload
}

func createWrappedPayloadHeaderCapella(t testing.TB) interfaces.ExecutionData {
	payload, err := blocks.WrappedExecutionPayloadHeaderCapella(&enginev1.ExecutionPayloadHeaderCapella{
		ParentHash:       make([]byte, fieldparams.RootLength),
		FeeRecipient:     make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:        make([]byte, fieldparams.RootLength),
		ReceiptsRoot:     make([]byte, fieldparams.RootLength),
		LogsBloom:        make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:       make([]byte, fieldparams.RootLength),
		BlockNumber:      0,
		GasLimit:         0,
		GasUsed:          0,
		Timestamp:        0,
		ExtraData:        make([]byte, 0),
		BaseFeePerGas:    make([]byte, fieldparams.RootLength),
		BlockHash:        make([]byte, fieldparams.RootLength),
		TransactionsRoot: make([]byte, fieldparams.RootLength),
		WithdrawalsRoot:  make([]byte, fieldparams.RootLength),
	})
	require.NoError(t, err)
	return payload
}

func TestCopyExecutionDataHeader(t *testing.T) {
	data := &enginev1.ExecutionPayloadHeader{GasUsed: 54}
	wsb, err := blocks.WrappedExecutionPayloadHeader(data)
	require.NoError(t, err)

	headerCopy, err := blocks.CopyExecutionDataHeader(wsb)
	require.NoError(t, err)

	assert.DeepEqual(t, data, headerCopy.Proto())
}

func TestCopyExecutionDataHeaderFromFullPayload(t *testing.T) {
	data := &enginev1.ExecutionPayload{GasUsed: 54}
	wsb, err := blocks.WrappedExecutionPayload(data)
	require.NoError(t, err)

	headerCopy, err := blocks.CopyExecutionDataHeader(wsb)
	require.NoError(t, err)

	txRoot, err := ssz.TransactionsRoot(nil)
	require.NoError(t, err)

	assert.DeepSSZEqual(t, &enginev1.ExecutionPayloadHeader{GasUsed: 54, TransactionsRoot: txRoot[:]}, headerCopy.Proto().(*enginev1.ExecutionPayloadHeader))
}

func TestCopyExecutionDataHeaderCapella(t *testing.T) {
	data := &enginev1.ExecutionPayloadHeaderCapella{
		ParentHash:       []byte("parenthash"),
		FeeRecipient:     []byte("feerecipient"),
		StateRoot:        []byte("stateroot"),
		ReceiptsRoot:     []byte("receiptsroot"),
		LogsBloom:        []byte("logsbloom"),
		PrevRandao:       []byte("prevrandao"),
		BlockNumber:      11,
		GasLimit:         22,
		GasUsed:          33,
		Timestamp:        44,
		ExtraData:        []byte("extradata"),
		BaseFeePerGas:    []byte("basefeepergas"),
		BlockHash:        []byte("blockhash"),
		TransactionsRoot: []byte("transactionsroot"),
		WithdrawalsRoot:  []byte("withdrawalsroot"),
	}
	payload, err := blocks.WrappedExecutionPayloadHeaderCapella(data)
	require.NoError(t, err)

	headerCopy, err := blocks.CopyExecutionDataHeader(payload)
	require.NoError(t, err)

	assert.DeepEqual(t, data, headerCopy.Proto())

	txRoot, err := headerCopy.TransactionsRoot()
	require.NoError(t, err)
	require.DeepEqual(t, txRoot, data.TransactionsRoot)

	wrRoot, err := headerCopy.WithdrawalsRoot()
	require.NoError(t, err)
	require.DeepEqual(t, wrRoot, data.WithdrawalsRoot)
}

func TestCopyExecutionDataHeaderCapellaFromFullPayload(t *testing.T) {
	data := &enginev1.ExecutionPayloadCapella{
		ParentHash:    []byte("parenthash"),
		FeeRecipient:  []byte("feerecipient"),
		StateRoot:     []byte("stateroot"),
		ReceiptsRoot:  []byte("receiptsroot"),
		LogsBloom:     []byte("logsbloom"),
		PrevRandao:    []byte("prevrandao"),
		BlockNumber:   11,
		GasLimit:      22,
		GasUsed:       33,
		Timestamp:     44,
		ExtraData:     []byte("extradata"),
		BaseFeePerGas: []byte("basefeepergas"),
		BlockHash:     []byte("blockhash"),
		Transactions:  [][]byte{[]byte("transaction")},
		Withdrawals: []*enginev1.Withdrawal{{
			WithdrawalIndex:  55,
			ValidatorIndex:   66,
			ExecutionAddress: []byte("executionaddress"),
			Amount:           77,
		}},
	}
	payload, err := blocks.WrappedExecutionPayloadCapella(data)
	require.NoError(t, err)

	headerCopy, err := blocks.CopyExecutionDataHeader(payload)
	require.NoError(t, err)

	txRoot, err := ssz.TransactionsRoot(data.Transactions)
	require.NoError(t, err)
	wRoot, err := ssz.WithdrawalSliceRoot(hash.CustomSHA256Hasher(), data.Withdrawals, fieldparams.MaxWithdrawalsPerPayload)
	require.NoError(t, err)
	dataHeader := &enginev1.ExecutionPayloadHeaderCapella{
		ParentHash:       data.ParentHash,
		FeeRecipient:     data.FeeRecipient,
		StateRoot:        data.StateRoot,
		ReceiptsRoot:     data.ReceiptsRoot,
		LogsBloom:        data.LogsBloom,
		PrevRandao:       data.PrevRandao,
		BlockNumber:      data.BlockNumber,
		GasLimit:         data.GasLimit,
		GasUsed:          data.GasUsed,
		Timestamp:        data.Timestamp,
		ExtraData:        data.ExtraData,
		BaseFeePerGas:    data.BaseFeePerGas,
		BlockHash:        data.BlockHash,
		TransactionsRoot: txRoot[:],
		WithdrawalsRoot:  wRoot[:],
	}
	assert.DeepSSZEqual(t, dataHeader, headerCopy.Proto())
}
