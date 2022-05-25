package powchain

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	newPayloadLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "new_payload_v1_latency_milliseconds",
			Help:    "Captures RPC latency for newPayloadV1 in milliseconds",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000},
		},
	)
	getPayloadLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "get_payload_v1_latency_milliseconds",
			Help:    "Captures RPC latency for getPayloadV1 in milliseconds",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000},
		},
	)
	forkchoiceUpdatedLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "forkchoice_updated_v1_latency_milliseconds",
			Help:    "Captures RPC latency for forkchoiceUpdatedV1 in milliseconds",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000},
		},
	)
	executionBlockByHashWithTxsLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "execution_block_by_hash_with_txs_latency_milliseconds",
			Help:    "Captures RPC latency for retrieving in milliseconds",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000},
		},
	)
	executionPayloadReconstructionLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "execution_payload_reconstruction_latency_milliseconds",
			Help:    "Captures RPC latency for retrieving in milliseconds",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000},
		},
	)
	reconstructedExecutionPayloadCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "reconstructed_execution_payload_count",
		Help: "Count the number of execution payloads that are reconstructed using JSON-RPC from payload headers",
	})
)
