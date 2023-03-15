// Package rpc defines a gRPC server implementing the Ethereum consensus API as needed
// by validator clients and consumers of chain data.
package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/builder"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/cache"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/cache/depositcache"
	blockfeed "github.com/prysmaticlabs/prysm/v3/beacon-chain/core/feed/block"
	opfeed "github.com/prysmaticlabs/prysm/v3/beacon-chain/core/feed/operation"
	statefeed "github.com/prysmaticlabs/prysm/v3/beacon-chain/core/feed/state"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/db"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/attestations"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/blstoexec"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/slashings"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/synccommittee"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/voluntaryexits"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/p2p"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/eth/beacon"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/eth/debug"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/eth/events"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/eth/node"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/eth/validator"
	beaconv1alpha1 "github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/prysm/v1alpha1/beacon"
	debugv1alpha1 "github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/prysm/v1alpha1/debug"
	nodev1alpha1 "github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/prysm/v1alpha1/node"
	validatorv1alpha1 "github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/prysm/v1alpha1/validator"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/rpc/statefetcher"
	slasherservice "github.com/prysmaticlabs/prysm/v3/beacon-chain/slasher"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/state/stategen"
	chainSync "github.com/prysmaticlabs/prysm/v3/beacon-chain/sync"
	"github.com/prysmaticlabs/prysm/v3/config/features"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/prysmaticlabs/prysm/v3/io/logs"
	"github.com/prysmaticlabs/prysm/v3/monitoring/tracing"
	ethpbservice "github.com/prysmaticlabs/prysm/v3/proto/eth/service"
	ethpbv1alpha1 "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

const attestationBufferSize = 100

// Service defining an RPC server for a beacon node.
type Service struct {
	cfg                  *Config
	ctx                  context.Context
	cancel               context.CancelFunc
	listener             net.Listener
	grpcServer           *grpc.Server
	incomingAttestation  chan *ethpbv1alpha1.Attestation
	credentialError      error
	connectedRPCClients  map[net.Addr]bool
	clientConnectionLock sync.Mutex
}

// Config options for the beacon node RPC server.
type Config struct {
	ExecutionPayloadReconstructor execution.ExecutionPayloadReconstructor
	Host                          string
	Port                          string
	CertFlag                      string
	KeyFlag                       string
	BeaconMonitoringHost          string
	BeaconMonitoringPort          int
	BeaconDB                      db.HeadAccessDatabase
	ChainInfoFetcher              blockchain.ChainInfoFetcher
	HeadUpdater                   blockchain.HeadUpdater
	HeadFetcher                   blockchain.HeadFetcher
	CanonicalFetcher              blockchain.CanonicalFetcher
	ForkFetcher                   blockchain.ForkFetcher
	FinalizationFetcher           blockchain.FinalizationFetcher
	AttestationReceiver           blockchain.AttestationReceiver
	BlockReceiver                 blockchain.BlockReceiver
	ExecutionChainService         execution.Chain
	ChainStartFetcher             execution.ChainStartFetcher
	ExecutionChainInfoFetcher     execution.ChainInfoFetcher
	GenesisTimeFetcher            blockchain.TimeFetcher
	GenesisFetcher                blockchain.GenesisFetcher
	EnableDebugRPCEndpoints       bool
	MockEth1Votes                 bool
	AttestationsPool              attestations.Pool
	ExitPool                      voluntaryexits.PoolManager
	SlashingsPool                 slashings.PoolManager
	BLSToExecPool                 blstoexec.PoolManager
	SlashingChecker               slasherservice.SlashingChecker
	SyncCommitteeObjectPool       synccommittee.Pool
	BLSChangesPool                blstoexec.PoolManager
	SyncService                   chainSync.Checker
	Broadcaster                   p2p.Broadcaster
	PeersFetcher                  p2p.PeersProvider
	PeerManager                   p2p.PeerManager
	MetadataProvider              p2p.MetadataProvider
	DepositFetcher                depositcache.DepositFetcher
	PendingDepositFetcher         depositcache.PendingDepositsFetcher
	StateNotifier                 statefeed.Notifier
	BlockNotifier                 blockfeed.Notifier
	OperationNotifier             opfeed.Notifier
	StateGen                      *stategen.State
	MaxMsgSize                    int
	ExecutionEngineCaller         execution.EngineCaller
	ProposerIdsCache              *cache.ProposerPayloadIDsCache
	OptimisticModeFetcher         blockchain.OptimisticModeFetcher
	BlockBuilder                  builder.BlockBuilder
}

// NewService instantiates a new RPC service instance that will
// be registered into a running beacon node.
func NewService(ctx context.Context, cfg *Config) *Service {
	ctx, cancel := context.WithCancel(ctx)
	s := &Service{
		cfg:                 cfg,
		ctx:                 ctx,
		cancel:              cancel,
		incomingAttestation: make(chan *ethpbv1alpha1.Attestation, params.BeaconConfig().DefaultBufferSize),
		connectedRPCClients: make(map[net.Addr]bool),
	}

	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.WithError(err).Errorf("Could not listen to port in Start() %s", address)
	}
	s.listener = lis
	log.WithField("address", address).Info("gRPC server listening on port")

	opts := []grpc.ServerOption{
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.StreamInterceptor(middleware.ChainStreamServer(
			recovery.StreamServerInterceptor(
				recovery.WithRecoveryHandlerContext(tracing.RecoveryHandlerFunc),
			),
			grpcprometheus.StreamServerInterceptor,
			grpcopentracing.StreamServerInterceptor(),
			s.validatorStreamConnectionInterceptor,
		)),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandlerContext(tracing.RecoveryHandlerFunc),
			),
			grpcprometheus.UnaryServerInterceptor,
			grpcopentracing.UnaryServerInterceptor(),
			s.validatorUnaryConnectionInterceptor,
		)),
		grpc.MaxRecvMsgSize(s.cfg.MaxMsgSize),
	}
	if s.cfg.CertFlag != "" && s.cfg.KeyFlag != "" {
		creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFlag, s.cfg.KeyFlag)
		if err != nil {
			log.WithError(err).Fatal("Could not load TLS keys")
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		log.Warn("You are using an insecure gRPC server. If you are running your beacon node and " +
			"validator on the same machines, you can ignore this message. If you want to know " +
			"how to enable secure connections, see: https://docs.prylabs.network/docs/prysm-usage/secure-grpc")
	}
	s.grpcServer = grpc.NewServer(opts...)

	return s
}

// paranoid build time check to ensure ChainInfoFetcher implements required interfaces
var _ stategen.CanonicalChecker = blockchain.ChainInfoFetcher(nil)
var _ stategen.CurrentSlotter = blockchain.ChainInfoFetcher(nil)

// Start the gRPC server.
func (s *Service) Start() {
	grpcprometheus.EnableHandlingTimeHistogram()

	var stateCache stategen.CachedGetter
	if s.cfg.StateGen != nil {
		stateCache = s.cfg.StateGen.CombinedCache()
	}
	withCache := stategen.WithCache(stateCache)
	ch := stategen.NewCanonicalHistory(s.cfg.BeaconDB, s.cfg.ChainInfoFetcher, s.cfg.ChainInfoFetcher, withCache)

	validatorServer := &validatorv1alpha1.Server{
		Ctx:                    s.ctx,
		AttestationCache:       cache.NewAttestationCache(),
		AttPool:                s.cfg.AttestationsPool,
		ExitPool:               s.cfg.ExitPool,
		HeadFetcher:            s.cfg.HeadFetcher,
		HeadUpdater:            s.cfg.HeadUpdater,
		ForkFetcher:            s.cfg.ForkFetcher,
		GenesisFetcher:         s.cfg.GenesisFetcher,
		FinalizationFetcher:    s.cfg.FinalizationFetcher,
		TimeFetcher:            s.cfg.GenesisTimeFetcher,
		BlockFetcher:           s.cfg.ExecutionChainService,
		DepositFetcher:         s.cfg.DepositFetcher,
		ChainStartFetcher:      s.cfg.ChainStartFetcher,
		Eth1InfoFetcher:        s.cfg.ExecutionChainService,
		OptimisticModeFetcher:  s.cfg.OptimisticModeFetcher,
		SyncChecker:            s.cfg.SyncService,
		StateNotifier:          s.cfg.StateNotifier,
		BlockNotifier:          s.cfg.BlockNotifier,
		OperationNotifier:      s.cfg.OperationNotifier,
		P2P:                    s.cfg.Broadcaster,
		BlockReceiver:          s.cfg.BlockReceiver,
		MockEth1Votes:          s.cfg.MockEth1Votes,
		Eth1BlockFetcher:       s.cfg.ExecutionChainService,
		PendingDepositsFetcher: s.cfg.PendingDepositFetcher,
		SlashingsPool:          s.cfg.SlashingsPool,
		StateGen:               s.cfg.StateGen,
		SyncCommitteePool:      s.cfg.SyncCommitteeObjectPool,
		BLSChangesPool:         s.cfg.BLSChangesPool,
		ReplayerBuilder:        ch,
		ExecutionEngineCaller:  s.cfg.ExecutionEngineCaller,
		BeaconDB:               s.cfg.BeaconDB,
		ProposerSlotIndexCache: s.cfg.ProposerIdsCache,
		BlockBuilder:           s.cfg.BlockBuilder,
	}
	validatorServerV1 := &validator.Server{
		HeadFetcher:           s.cfg.HeadFetcher,
		HeadUpdater:           s.cfg.HeadUpdater,
		TimeFetcher:           s.cfg.GenesisTimeFetcher,
		SyncChecker:           s.cfg.SyncService,
		OptimisticModeFetcher: s.cfg.OptimisticModeFetcher,
		AttestationsPool:      s.cfg.AttestationsPool,
		PeerManager:           s.cfg.PeerManager,
		Broadcaster:           s.cfg.Broadcaster,
		V1Alpha1Server:        validatorServer,
		StateFetcher: &statefetcher.StateProvider{
			BeaconDB:           s.cfg.BeaconDB,
			ChainInfoFetcher:   s.cfg.ChainInfoFetcher,
			GenesisTimeFetcher: s.cfg.GenesisTimeFetcher,
			StateGenService:    s.cfg.StateGen,
			ReplayerBuilder:    ch,
		},
		SyncCommitteePool:      s.cfg.SyncCommitteeObjectPool,
		BLSChangesPool:         s.cfg.BLSChangesPool,
		ProposerSlotIndexCache: s.cfg.ProposerIdsCache,
	}

	nodeServer := &nodev1alpha1.Server{
		LogsStreamer:         logs.NewStreamServer(),
		StreamLogsBufferSize: 1000, // Enough to handle bursts of beacon node logs for gRPC streaming.
		BeaconDB:             s.cfg.BeaconDB,
		Server:               s.grpcServer,
		SyncChecker:          s.cfg.SyncService,
		GenesisTimeFetcher:   s.cfg.GenesisTimeFetcher,
		PeersFetcher:         s.cfg.PeersFetcher,
		PeerManager:          s.cfg.PeerManager,
		GenesisFetcher:       s.cfg.GenesisFetcher,
		POWChainInfoFetcher:  s.cfg.ExecutionChainInfoFetcher,
		BeaconMonitoringHost: s.cfg.BeaconMonitoringHost,
		BeaconMonitoringPort: s.cfg.BeaconMonitoringPort,
	}
	nodeServerV1 := &node.Server{
		BeaconDB:                  s.cfg.BeaconDB,
		Server:                    s.grpcServer,
		SyncChecker:               s.cfg.SyncService,
		OptimisticModeFetcher:     s.cfg.OptimisticModeFetcher,
		GenesisTimeFetcher:        s.cfg.GenesisTimeFetcher,
		PeersFetcher:              s.cfg.PeersFetcher,
		PeerManager:               s.cfg.PeerManager,
		MetadataProvider:          s.cfg.MetadataProvider,
		HeadFetcher:               s.cfg.HeadFetcher,
		ExecutionChainInfoFetcher: s.cfg.ExecutionChainInfoFetcher,
	}

	beaconChainServer := &beaconv1alpha1.Server{
		Ctx:                         s.ctx,
		BeaconDB:                    s.cfg.BeaconDB,
		AttestationsPool:            s.cfg.AttestationsPool,
		SlashingsPool:               s.cfg.SlashingsPool,
		HeadUpdater:                 s.cfg.HeadUpdater,
		OptimisticModeFetcher:       s.cfg.OptimisticModeFetcher,
		HeadFetcher:                 s.cfg.HeadFetcher,
		FinalizationFetcher:         s.cfg.FinalizationFetcher,
		CanonicalFetcher:            s.cfg.CanonicalFetcher,
		ChainStartFetcher:           s.cfg.ChainStartFetcher,
		DepositFetcher:              s.cfg.DepositFetcher,
		BlockFetcher:                s.cfg.ExecutionChainService,
		GenesisTimeFetcher:          s.cfg.GenesisTimeFetcher,
		StateNotifier:               s.cfg.StateNotifier,
		BlockNotifier:               s.cfg.BlockNotifier,
		AttestationNotifier:         s.cfg.OperationNotifier,
		Broadcaster:                 s.cfg.Broadcaster,
		StateGen:                    s.cfg.StateGen,
		SyncChecker:                 s.cfg.SyncService,
		ReceivedAttestationsBuffer:  make(chan *ethpbv1alpha1.Attestation, attestationBufferSize),
		CollectedAttestationsBuffer: make(chan []*ethpbv1alpha1.Attestation, attestationBufferSize),
		ReplayerBuilder:             ch,
		BLSChangesPool:              s.cfg.BLSChangesPool,
	}
	beaconChainServerV1 := &beacon.Server{
		CanonicalHistory:   ch,
		BeaconDB:           s.cfg.BeaconDB,
		AttestationsPool:   s.cfg.AttestationsPool,
		SlashingsPool:      s.cfg.SlashingsPool,
		ChainInfoFetcher:   s.cfg.ChainInfoFetcher,
		GenesisTimeFetcher: s.cfg.GenesisTimeFetcher,
		BlockNotifier:      s.cfg.BlockNotifier,
		OperationNotifier:  s.cfg.OperationNotifier,
		Broadcaster:        s.cfg.Broadcaster,
		BlockReceiver:      s.cfg.BlockReceiver,
		StateGenService:    s.cfg.StateGen,
		StateFetcher: &statefetcher.StateProvider{
			BeaconDB:           s.cfg.BeaconDB,
			ChainInfoFetcher:   s.cfg.ChainInfoFetcher,
			GenesisTimeFetcher: s.cfg.GenesisTimeFetcher,
			StateGenService:    s.cfg.StateGen,
			ReplayerBuilder:    ch,
		},
		OptimisticModeFetcher:         s.cfg.OptimisticModeFetcher,
		HeadFetcher:                   s.cfg.HeadFetcher,
		VoluntaryExitsPool:            s.cfg.ExitPool,
		V1Alpha1ValidatorServer:       validatorServer,
		SyncChecker:                   s.cfg.SyncService,
		ExecutionPayloadReconstructor: s.cfg.ExecutionPayloadReconstructor,
		BLSChangesPool:                s.cfg.BLSChangesPool,
		FinalizationFetcher:           s.cfg.FinalizationFetcher,
	}
	ethpbv1alpha1.RegisterNodeServer(s.grpcServer, nodeServer)
	ethpbservice.RegisterBeaconNodeServer(s.grpcServer, nodeServerV1)
	ethpbv1alpha1.RegisterHealthServer(s.grpcServer, nodeServer)
	ethpbv1alpha1.RegisterBeaconChainServer(s.grpcServer, beaconChainServer)
	ethpbservice.RegisterBeaconChainServer(s.grpcServer, beaconChainServerV1)
	ethpbservice.RegisterEventsServer(s.grpcServer, &events.Server{
		Ctx:               s.ctx,
		StateNotifier:     s.cfg.StateNotifier,
		BlockNotifier:     s.cfg.BlockNotifier,
		OperationNotifier: s.cfg.OperationNotifier,
	})
	if s.cfg.EnableDebugRPCEndpoints {
		log.Info("Enabled debug gRPC endpoints")
		debugServer := &debugv1alpha1.Server{
			GenesisTimeFetcher: s.cfg.GenesisTimeFetcher,
			BeaconDB:           s.cfg.BeaconDB,
			StateGen:           s.cfg.StateGen,
			HeadFetcher:        s.cfg.HeadFetcher,
			PeerManager:        s.cfg.PeerManager,
			PeersFetcher:       s.cfg.PeersFetcher,
			ReplayerBuilder:    ch,
		}
		debugServerV1 := &debug.Server{
			BeaconDB:    s.cfg.BeaconDB,
			HeadFetcher: s.cfg.HeadFetcher,
			StateFetcher: &statefetcher.StateProvider{
				BeaconDB:           s.cfg.BeaconDB,
				ChainInfoFetcher:   s.cfg.ChainInfoFetcher,
				GenesisTimeFetcher: s.cfg.GenesisTimeFetcher,
				StateGenService:    s.cfg.StateGen,
				ReplayerBuilder:    ch,
			},
			OptimisticModeFetcher: s.cfg.OptimisticModeFetcher,
			ForkFetcher:           s.cfg.ForkFetcher,
			FinalizationFetcher:   s.cfg.FinalizationFetcher,
		}
		ethpbv1alpha1.RegisterDebugServer(s.grpcServer, debugServer)
		ethpbservice.RegisterBeaconDebugServer(s.grpcServer, debugServerV1)
	}
	ethpbv1alpha1.RegisterBeaconNodeValidatorServer(s.grpcServer, validatorServer)
	ethpbservice.RegisterBeaconValidatorServer(s.grpcServer, validatorServerV1)
	// Register reflection service on gRPC server.
	reflection.Register(s.grpcServer)

	go func() {
		if s.listener != nil {
			if err := s.grpcServer.Serve(s.listener); err != nil {
				log.WithError(err).Errorf("Could not serve gRPC")
			}
		}
	}()
}

// Stop the service.
func (s *Service) Stop() error {
	s.cancel()
	if s.listener != nil {
		s.grpcServer.GracefulStop()
		log.Debug("Initiated graceful stop of gRPC server")
	}
	return nil
}

// Status returns nil or credentialError
func (s *Service) Status() error {
	optimistic, err := s.cfg.OptimisticModeFetcher.IsOptimistic(s.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to check if service is optimistic")
	}
	if optimistic {
		return errors.New("service is optimistic, validators can't perform duties " +
			"please check if execution layer is fully synced")
	}
	if s.cfg.SyncService.Syncing() {
		return errors.New("syncing")
	}
	if s.credentialError != nil {
		return s.credentialError
	}
	return nil
}

// Stream interceptor for new validator client connections to the beacon node.
func (s *Service) validatorStreamConnectionInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	s.logNewClientConnection(ss.Context())
	return handler(srv, ss)
}

// Unary interceptor for new validator client connections to the beacon node.
func (s *Service) validatorUnaryConnectionInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	s.logNewClientConnection(ctx)
	return handler(ctx, req)
}

func (s *Service) logNewClientConnection(ctx context.Context) {
	if features.Get().DisableGRPCConnectionLogs {
		return
	}
	if clientInfo, ok := peer.FromContext(ctx); ok {
		// Check if we have not yet observed this grpc client connection
		// in the running beacon node.
		s.clientConnectionLock.Lock()
		defer s.clientConnectionLock.Unlock()
		if !s.connectedRPCClients[clientInfo.Addr] {
			log.WithFields(logrus.Fields{
				"addr": clientInfo.Addr.String(),
			}).Infof("gRPC client connected to beacon node")
			s.connectedRPCClients[clientInfo.Addr] = true
		}
	}
}
