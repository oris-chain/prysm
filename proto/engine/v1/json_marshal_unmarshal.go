package enginev1

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	fieldparams "github.com/prysmaticlabs/prysm/config/fieldparams"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
)

// PayloadIDBytes defines a custom type for Payload IDs used by the engine API
// client with proper JSON Marshal and Unmarshal methods to hex.
type PayloadIDBytes [8]byte

// MarshalJSON --
func (b PayloadIDBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hexutil.Bytes(b[:]))
}

// ExecutionBlock is the response kind received by the eth_getBlockByHash and
// eth_getBlockByNumber endpoints via JSON-RPC.
type ExecutionBlock struct {
	gethtypes.Header
	Hash            common.Hash              `json:"hash"`
	Transactions    []*gethtypes.Transaction `json:"transactions"`
	TotalDifficulty string                   `json:"totalDifficulty"`
}

func (e *ExecutionBlock) MarshalJSON() ([]byte, error) {
	decoded := make(map[string]interface{})
	encodedHeader, err := e.Header.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(encodedHeader, &decoded); err != nil {
		return nil, err
	}
	encodedTxs, err := json.Marshal(e.Transactions)
	if err != nil {
		return nil, err
	}
	decoded["hash"] = e.Hash.String()
	decoded["transactions"] = string(encodedTxs)
	decoded["totalDifficulty"] = e.TotalDifficulty
	return json.Marshal(decoded)
}

func (e *ExecutionBlock) UnmarshalJSON(enc []byte) error {
	if err := e.Header.UnmarshalJSON(enc); err != nil {
		return err
	}
	decoded := make(map[string]interface{})
	if err := json.Unmarshal(enc, &decoded); err != nil {
		return err
	}
	blockHashStr, ok := decoded["hash"].(string)
	if !ok {
		return errors.New("expected `hash` field in JSON response")
	}
	e.Hash = common.HexToHash(blockHashStr)
	e.TotalDifficulty, ok = decoded["totalDifficulty"].(string)
	if !ok {
		return errors.New("expected `totalDifficulty` field in JSON response")
	}
	txsList, ok := decoded["transactions"].([]interface{})
	if !ok {
		return nil
	}
	// If the block contains a list of transactions, we JSON unmarshal
	// them into a list of geth transaction objects.
	txs := make([]*gethtypes.Transaction, len(txsList))
	for i, tx := range txsList {
		t := &gethtypes.Transaction{}
		encodedTx, err := json.Marshal(tx)
		if err != nil {
			return errors.Wrapf(err, "could not marshal tx %v", tx)
		}
		if err := json.Unmarshal(encodedTx, &t); err != nil {
			return errors.Wrapf(err, "could not marshal tx %s", string(encodedTx))
		}
		txs[i] = t
	}
	e.Transactions = txs
	return nil
}

// UnmarshalJSON --
func (b *PayloadIDBytes) UnmarshalJSON(enc []byte) error {
	hexBytes := hexutil.Bytes(make([]byte, 0))
	if err := json.Unmarshal(enc, &hexBytes); err != nil {
		return err
	}
	res := [8]byte{}
	copy(res[:], hexBytes)
	*b = res
	return nil
}

type executionPayloadJSON struct {
	ParentHash    *common.Hash    `json:"parentHash"`
	FeeRecipient  *common.Address `json:"feeRecipient"`
	StateRoot     *common.Hash    `json:"stateRoot"`
	ReceiptsRoot  *common.Hash    `json:"receiptsRoot"`
	LogsBloom     *hexutil.Bytes  `json:"logsBloom"`
	PrevRandao    *common.Hash    `json:"prevRandao"`
	BlockNumber   *hexutil.Uint64 `json:"blockNumber"`
	GasLimit      *hexutil.Uint64 `json:"gasLimit"`
	GasUsed       *hexutil.Uint64 `json:"gasUsed"`
	Timestamp     *hexutil.Uint64 `json:"timestamp"`
	ExtraData     hexutil.Bytes   `json:"extraData"`
	BaseFeePerGas string          `json:"baseFeePerGas"`
	BlockHash     *common.Hash    `json:"blockHash"`
	Transactions  []hexutil.Bytes `json:"transactions"`
}

// MarshalJSON --
func (e *ExecutionPayload) MarshalJSON() ([]byte, error) {
	transactions := make([]hexutil.Bytes, len(e.Transactions))
	for i, tx := range e.Transactions {
		transactions[i] = tx
	}
	baseFee := new(big.Int).SetBytes(bytesutil.ReverseByteOrder(e.BaseFeePerGas))
	baseFeeHex := hexutil.EncodeBig(baseFee)
	pHash := common.BytesToHash(e.ParentHash)
	sRoot := common.BytesToHash(e.StateRoot)
	recRoot := common.BytesToHash(e.ReceiptsRoot)
	prevRan := common.BytesToHash(e.PrevRandao)
	bHash := common.BytesToHash(e.BlockHash)
	blockNum := hexutil.Uint64(e.BlockNumber)
	gasLimit := hexutil.Uint64(e.GasLimit)
	gasUsed := hexutil.Uint64(e.GasUsed)
	timeStamp := hexutil.Uint64(e.Timestamp)
	recipient := common.BytesToAddress(e.FeeRecipient)
	logsBloom := hexutil.Bytes(e.LogsBloom)
	return json.Marshal(executionPayloadJSON{
		ParentHash:    &pHash,
		FeeRecipient:  &recipient,
		StateRoot:     &sRoot,
		ReceiptsRoot:  &recRoot,
		LogsBloom:     &logsBloom,
		PrevRandao:    &prevRan,
		BlockNumber:   &blockNum,
		GasLimit:      &gasLimit,
		GasUsed:       &gasUsed,
		Timestamp:     &timeStamp,
		ExtraData:     e.ExtraData,
		BaseFeePerGas: baseFeeHex,
		BlockHash:     &bHash,
		Transactions:  transactions,
	})
}

// UnmarshalJSON --
func (e *ExecutionPayload) UnmarshalJSON(enc []byte) error {
	dec := executionPayloadJSON{}
	if err := json.Unmarshal(enc, &dec); err != nil {
		return err
	}

	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for ExecutionPayload")
	}
	if dec.FeeRecipient == nil {
		return errors.New("missing required field 'feeRecipient' for ExecutionPayload")
	}
	if dec.StateRoot == nil {
		return errors.New("missing required field 'stateRoot' for ExecutionPayload")
	}
	if dec.ReceiptsRoot == nil {
		return errors.New("missing required field 'receiptsRoot' for ExecutableDataV1")
	}

	if dec.LogsBloom == nil {
		return errors.New("missing required field 'logsBloom' for ExecutionPayload")
	}
	if dec.PrevRandao == nil {
		return errors.New("missing required field 'prevRandao' for ExecutionPayload")
	}
	if dec.ExtraData == nil {
		return errors.New("missing required field 'extraData' for ExecutionPayload")
	}
	if dec.BlockHash == nil {
		return errors.New("missing required field 'blockHash' for ExecutionPayload")
	}
	if dec.Transactions == nil {
		return errors.New("missing required field 'transactions' for ExecutionPayload")
	}
	if dec.BlockNumber == nil {
		return errors.New("missing required field 'blockNumber' for ExecutionPayload")
	}
	if dec.Timestamp == nil {
		return errors.New("missing required field 'timestamp' for ExecutionPayload")
	}
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for ExecutionPayload")
	}
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for ExecutionPayload")
	}
	*e = ExecutionPayload{}
	e.ParentHash = dec.ParentHash.Bytes()
	e.FeeRecipient = dec.FeeRecipient.Bytes()
	e.StateRoot = dec.StateRoot.Bytes()
	e.ReceiptsRoot = dec.ReceiptsRoot.Bytes()
	e.LogsBloom = *dec.LogsBloom
	e.PrevRandao = dec.PrevRandao.Bytes()
	e.BlockNumber = uint64(*dec.BlockNumber)
	e.GasLimit = uint64(*dec.GasLimit)
	e.GasUsed = uint64(*dec.GasUsed)
	e.Timestamp = uint64(*dec.Timestamp)
	e.ExtraData = dec.ExtraData
	baseFee, err := hexutil.DecodeBig(dec.BaseFeePerGas)
	if err != nil {
		return err
	}
	e.BaseFeePerGas = bytesutil.PadTo(bytesutil.ReverseByteOrder(baseFee.Bytes()), fieldparams.RootLength)
	e.BlockHash = dec.BlockHash.Bytes()
	transactions := make([][]byte, len(dec.Transactions))
	for i, tx := range dec.Transactions {
		transactions[i] = tx
	}
	e.Transactions = transactions
	return nil
}

type payloadAttributesJSON struct {
	Timestamp             hexutil.Uint64 `json:"timestamp"`
	PrevRandao            hexutil.Bytes  `json:"prevRandao"`
	SuggestedFeeRecipient hexutil.Bytes  `json:"suggestedFeeRecipient"`
}

// MarshalJSON --
func (p *PayloadAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(payloadAttributesJSON{
		Timestamp:             hexutil.Uint64(p.Timestamp),
		PrevRandao:            p.PrevRandao,
		SuggestedFeeRecipient: p.SuggestedFeeRecipient,
	})
}

// UnmarshalJSON --
func (p *PayloadAttributes) UnmarshalJSON(enc []byte) error {
	dec := payloadAttributesJSON{}
	if err := json.Unmarshal(enc, &dec); err != nil {
		return err
	}
	*p = PayloadAttributes{}
	p.Timestamp = uint64(dec.Timestamp)
	p.PrevRandao = dec.PrevRandao
	p.SuggestedFeeRecipient = dec.SuggestedFeeRecipient
	return nil
}

type payloadStatusJSON struct {
	LatestValidHash *common.Hash `json:"latestValidHash"`
	Status          string       `json:"status"`
	ValidationError *string      `json:"validationError"`
}

// MarshalJSON --
func (p *PayloadStatus) MarshalJSON() ([]byte, error) {
	var latestHash *common.Hash
	if p.LatestValidHash != nil {
		hash := common.Hash(bytesutil.ToBytes32(p.LatestValidHash))
		latestHash = (*common.Hash)(&hash)
	}
	return json.Marshal(payloadStatusJSON{
		LatestValidHash: latestHash,
		Status:          p.Status.String(),
		ValidationError: &p.ValidationError,
	})
}

// UnmarshalJSON --
func (p *PayloadStatus) UnmarshalJSON(enc []byte) error {
	dec := payloadStatusJSON{}
	if err := json.Unmarshal(enc, &dec); err != nil {
		return err
	}
	*p = PayloadStatus{}
	if dec.LatestValidHash != nil {
		p.LatestValidHash = dec.LatestValidHash[:]
	}
	p.Status = PayloadStatus_Status(PayloadStatus_Status_value[dec.Status])
	if dec.ValidationError != nil {
		p.ValidationError = *dec.ValidationError
	}
	return nil
}

type transitionConfigurationJSON struct {
	TerminalTotalDifficulty string        `json:"terminalTotalDifficulty"`
	TerminalBlockHash       hexutil.Bytes `json:"terminalBlockHash"`
	TerminalBlockNumber     string        `json:"terminalBlockNumber"`
}

// MarshalJSON --
func (t *TransitionConfiguration) MarshalJSON() ([]byte, error) {
	num := new(big.Int).SetBytes(t.TerminalBlockNumber)
	numHex := hexutil.EncodeBig(num)
	return json.Marshal(transitionConfigurationJSON{
		TerminalTotalDifficulty: t.TerminalTotalDifficulty,
		TerminalBlockHash:       t.TerminalBlockHash,
		TerminalBlockNumber:     numHex,
	})
}

// UnmarshalJSON --
func (t *TransitionConfiguration) UnmarshalJSON(enc []byte) error {
	dec := transitionConfigurationJSON{}
	if err := json.Unmarshal(enc, &dec); err != nil {
		return err
	}
	*t = TransitionConfiguration{}
	num, err := hexutil.DecodeBig(dec.TerminalBlockNumber)
	if err != nil {
		return err
	}
	t.TerminalTotalDifficulty = dec.TerminalTotalDifficulty
	t.TerminalBlockHash = dec.TerminalBlockHash
	t.TerminalBlockNumber = num.Bytes()
	return nil
}

type forkchoiceStateJSON struct {
	HeadBlockHash      hexutil.Bytes `json:"headBlockHash"`
	SafeBlockHash      hexutil.Bytes `json:"safeBlockHash"`
	FinalizedBlockHash hexutil.Bytes `json:"finalizedBlockHash"`
}

// MarshalJSON --
func (f *ForkchoiceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(forkchoiceStateJSON{
		HeadBlockHash:      f.HeadBlockHash,
		SafeBlockHash:      f.SafeBlockHash,
		FinalizedBlockHash: f.FinalizedBlockHash,
	})
}

// UnmarshalJSON --
func (f *ForkchoiceState) UnmarshalJSON(enc []byte) error {
	dec := forkchoiceStateJSON{}
	if err := json.Unmarshal(enc, &dec); err != nil {
		return err
	}
	*f = ForkchoiceState{}
	f.HeadBlockHash = dec.HeadBlockHash
	f.SafeBlockHash = dec.SafeBlockHash
	f.FinalizedBlockHash = dec.FinalizedBlockHash
	return nil
}
