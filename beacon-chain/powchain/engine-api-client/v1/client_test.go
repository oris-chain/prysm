package v1

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	pb "github.com/prysmaticlabs/prysm/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func TestClient_IPC(t *testing.T) {
	server := newTestIPCServer(t)
	defer server.Stop()
	rpcClient := rpc.DialInProc(server)
	defer rpcClient.Close()
	client := &Client{}
	client.rpc = rpcClient
	ctx := context.Background()
	fix := fixtures()

	t.Run("engine_getPayloadV1", func(t *testing.T) {
		want := fix["ExecutionPayload"].(*pb.ExecutionPayload)
		payloadId := [8]byte{1}
		resp, err := client.GetPayload(ctx, payloadId)
		require.NoError(t, err)
		require.DeepEqual(t, want, resp)
	})
	t.Run("engine_forkchoiceUpdatedV1", func(t *testing.T) {
		want := fix["ForkchoiceUpdatedResponse"].(*ForkchoiceUpdatedResponse)
		resp, err := client.ForkchoiceUpdated(ctx, &pb.ForkchoiceState{}, &pb.PayloadAttributes{})
		require.NoError(t, err)
		require.DeepEqual(t, want.Status, resp.Status)
		require.DeepEqual(t, want.PayloadId, resp.PayloadId)
	})
	t.Run("engine_newPayloadV1", func(t *testing.T) {
		want := fix["PayloadStatus"].(*pb.PayloadStatus)
		resp, err := client.NewPayload(ctx, &pb.ExecutionPayload{})
		require.NoError(t, err)
		require.DeepEqual(t, want, resp)
	})
}

func TestClient_HTTP(t *testing.T) {

}

type customError struct {
	code int
}

func (c *customError) ErrorCode() int {
	return c.code
}

func (c *customError) Error() string {
	return "something went wrong"
}

type dataError struct {
	code int
	data interface{}
}

func (c *dataError) ErrorCode() int {
	return c.code
}

func (c *dataError) Error() string {
	return "something went wrong"
}

func (c *dataError) ErrorData() interface{} {
	return c.data
}

func Test_handleRPCError(t *testing.T) {
	got := handleRPCError(nil)
	require.Equal(t, true, got == nil)

	var tests = []struct {
		name             string
		expected         error
		expectedContains string
		given            error
	}{
		{
			name:             "not an rpc error",
			expectedContains: "got an unexpected error",
			given:            errors.New("foo"),
		},
		{
			name:             "ErrParse",
			expectedContains: ErrParse.Error(),
			given:            &customError{code: -32700},
		},
		{
			name:             "ErrInvalidRequest",
			expectedContains: ErrInvalidRequest.Error(),
			given:            &customError{code: -32600},
		},
		{
			name:             "ErrMethodNotFound",
			expectedContains: ErrMethodNotFound.Error(),
			given:            &customError{code: -32601},
		},
		{
			name:             "ErrInvalidParams",
			expectedContains: ErrInvalidParams.Error(),
			given:            &customError{code: -32602},
		},
		{
			name:             "ErrInternal",
			expectedContains: ErrInternal.Error(),
			given:            &customError{code: -32603},
		},
		{
			name:             "ErrUnknownPayload",
			expectedContains: ErrUnknownPayload.Error(),
			given:            &customError{code: -32001},
		},
		{
			name:             "ErrServer unexpected no data",
			expectedContains: "got an unexpected error",
			given:            &customError{code: -32000},
		},
		{
			name:             "ErrServer with data",
			expectedContains: ErrServer.Error(),
			given:            &dataError{code: -32000, data: 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleRPCError(tt.given)
			require.ErrorContains(t, tt.expectedContains, got)
		})
	}
}

func newTestIPCServer(t *testing.T) *rpc.Server {
	server := rpc.NewServer()
	err := server.RegisterName("engine", new(testEngineService))
	require.NoError(t, err)
	return server
}

func fixtures() map[string]interface{} {
	foo := bytesutil.ToBytes32([]byte("foo"))
	bar := bytesutil.PadTo([]byte("bar"), 20)
	baz := bytesutil.PadTo([]byte("baz"), 256)
	executionPayloadFixture := &pb.ExecutionPayload{
		ParentHash:    foo[:],
		FeeRecipient:  bar,
		StateRoot:     foo[:],
		ReceiptsRoot:  foo[:],
		LogsBloom:     baz,
		Random:        foo[:],
		BlockNumber:   1,
		GasLimit:      1,
		GasUsed:       1,
		Timestamp:     1,
		ExtraData:     foo[:],
		BaseFeePerGas: foo[:],
		BlockHash:     foo[:],
		Transactions:  [][]byte{foo[:]},
	}
	status := &pb.PayloadStatus{
		Status:          pb.PayloadStatus_ACCEPTED,
		LatestValidHash: foo[:],
		ValidationError: "",
	}
	forkChoiceResp := &ForkchoiceUpdatedResponse{
		Status:    status,
		PayloadId: [8]byte{1},
	}
	return map[string]interface{}{
		"ExecutionPayload":          executionPayloadFixture,
		"PayloadStatus":             status,
		"ForkchoiceUpdatedResponse": forkChoiceResp,
	}
}

type testEngineService struct{}

type echoArgs struct {
	S string
}

type echoResult struct {
	String string
	Int    int
	Args   *echoArgs
}

type testError struct{}

func (testError) Error() string          { return "testError" }
func (testError) ErrorCode() int         { return 444 }
func (testError) ErrorData() interface{} { return "testError data" }

func (s *testEngineService) NoArgsRets() {}

func (s *testEngineService) GetPayloadV1(
	ctx context.Context, payloadId [8]byte,
) *pb.ExecutionPayload {
	fix := fixtures()
	return fix["ExecutionPayload"].(*pb.ExecutionPayload)
}

func (s *testEngineService) ForkchoiceUpdatedV1(
	ctx context.Context, state *pb.ForkchoiceState, attrs *pb.PayloadAttributes,
) *ForkchoiceUpdatedResponse {
	fix := fixtures()
	return fix["ForkchoiceUpdatedResponse"].(*ForkchoiceUpdatedResponse)
}

func (s *testEngineService) NewPayloadV1(
	ctx context.Context, payload *pb.ExecutionPayload,
) *pb.PayloadStatus {
	fix := fixtures()
	return fix["PayloadStatus"].(*pb.PayloadStatus)
}
