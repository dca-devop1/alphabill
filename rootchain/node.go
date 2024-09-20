package rootchain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alphabill-org/alphabill-go-base/types"
	"github.com/alphabill-org/alphabill/internal/debug"
	"github.com/alphabill-org/alphabill/logger"
	"github.com/alphabill-org/alphabill/network"
	"github.com/alphabill-org/alphabill/network/protocol/certification"
	"github.com/alphabill-org/alphabill/network/protocol/handshake"
	"github.com/alphabill-org/alphabill/observability"
	"github.com/alphabill-org/alphabill/rootchain/consensus"
	"github.com/alphabill-org/alphabill/rootchain/partitions"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

type (
	PartitionNet interface {
		Send(ctx context.Context, msg any, receivers ...peer.ID) error
		ReceivedChannel() <-chan any
	}

	Observability interface {
		Meter(name string, opts ...metric.MeterOption) metric.Meter
		Tracer(name string, options ...trace.TracerOption) trace.Tracer
		Logger() *slog.Logger
	}

	ConsensusManager interface {
		// RequestCertification accepts certification requests with proof of quorum or no-quorum.
		RequestCertification(ctx context.Context, cr consensus.IRChangeRequest) error
		// CertificationResult read the channel to receive certification results
		CertificationResult() <-chan *types.UnicityCertificate
		// GetLatestUnicityCertificate get the latest certification for partition (maybe should/can be removed)
		GetLatestUnicityCertificate(id types.SystemID) (*types.UnicityCertificate, error)
		// Run consensus algorithm
		Run(ctx context.Context) error
	}

	Node struct {
		peer             *network.Peer // p2p network host for partition
		partitions       partitions.PartitionConfiguration
		incomingRequests *CertRequestBuffer
		subscription     *Subscriptions
		net              PartitionNet
		consensusManager ConsensusManager

		log    *slog.Logger
		tracer trace.Tracer

		bcrCount metric.Int64Counter // Block Certification Request count
	}
)

// New creates a new instance of the root chain node
func New(
	p *network.Peer,
	pNet PartitionNet,
	ps partitions.PartitionConfiguration,
	cm ConsensusManager,
	observe Observability,
) (*Node, error) {
	if p == nil {
		return nil, fmt.Errorf("partition listener is nil")
	}
	if pNet == nil {
		return nil, fmt.Errorf("network is nil")
	}

	meter := observe.Meter("rootchain.node", metric.WithInstrumentationAttributes(observability.PeerID("node.id", p.ID())))
	node := &Node{
		peer:             p,
		partitions:       ps,
		incomingRequests: NewCertificationRequestBuffer(),
		subscription:     NewSubscriptions(meter),
		net:              pNet,
		consensusManager: cm,
		log:              observe.Logger(),
		tracer:           observe.Tracer("rootchain.node"),
	}
	if err := node.initMetrics(meter); err != nil {
		return nil, fmt.Errorf("initializing metrics: %w", err)
	}
	return node, nil
}

func (v *Node) initMetrics(m metric.Meter) (err error) {
	v.bcrCount, err = m.Int64Counter("block.cert.req", metric.WithDescription("Number of Block Certification Requests processed"))
	if err != nil {
		return fmt.Errorf("creating Block Certification Requests counter: %w", err)
	}

	return nil
}

func (v *Node) Run(ctx context.Context) error {
	v.log.InfoContext(ctx, fmt.Sprintf("Starting root node. Addresses=%v; BuildInfo=%s", v.peer.MultiAddresses(), debug.ReadBuildInfo()))
	g, gctx := errgroup.WithContext(ctx)
	// Run root consensus algorithm
	g.Go(func() error { return v.consensusManager.Run(gctx) })
	// Start receiving messages from partition nodes
	g.Go(func() error { return v.loop(gctx) })
	// Start handling certification responses
	g.Go(func() error { return v.handleConsensus(gctx) })
	return g.Wait()
}

func (v *Node) GetPeer() *network.Peer {
	return v.peer
}

// loop handles messages from different goroutines.
func (v *Node) loop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-v.net.ReceivedChannel():
			if !ok {
				return fmt.Errorf("partition channel closed")
			}
			v.log.LogAttrs(ctx, logger.LevelTrace, fmt.Sprintf("received %T", msg), logger.Data(msg))
			switch mt := msg.(type) {
			case *certification.BlockCertificationRequest:
				if err := v.onBlockCertificationRequest(ctx, mt); err != nil {
					v.log.LogAttrs(ctx, slog.LevelWarn, fmt.Sprintf("handling block certification request from %s", mt.NodeIdentifier), logger.Error(err))
				}
			case *handshake.Handshake:
				if err := v.onHandshake(ctx, mt); err != nil {
					v.log.LogAttrs(ctx, slog.LevelWarn, fmt.Sprintf("handling handshake from %s", mt.NodeIdentifier), logger.Error(err))
				}
			default:
				v.log.LogAttrs(ctx, slog.LevelWarn, fmt.Sprintf("message %T not supported.", msg))
			}
		}
	}
}

func (v *Node) sendResponse(ctx context.Context, nodeID string, uc *types.UnicityCertificate) error {
	ctx, span := v.tracer.Start(ctx, "node.sendResponse")
	defer span.End()

	peerID, err := peer.Decode(nodeID)
	if err != nil {
		return fmt.Errorf("invalid receiver id: %w", err)
	}
	return v.net.Send(ctx, uc, peerID)
}

func (v *Node) onHandshake(ctx context.Context, req *handshake.Handshake) error {
	if err := req.IsValid(); err != nil {
		return fmt.Errorf("invalid handshake request: %w", err)
	}
	latestUnicityCertificate, err := v.consensusManager.GetLatestUnicityCertificate(req.SystemIdentifier)
	if err != nil {
		return fmt.Errorf("reading partition %s certificate: %w", req.SystemIdentifier, err)
	}
	if err = v.sendResponse(ctx, req.NodeIdentifier, latestUnicityCertificate); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return nil
}

/*
onBlockCertificationRequest handles Certification Request from partition nodes.
Partition nodes can only extend the stored/certified state.
*/
func (v *Node) onBlockCertificationRequest(ctx context.Context, req *certification.BlockCertificationRequest) (rErr error) {
	ctx, span := v.tracer.Start(ctx, "node.onBlockCertificationRequest")
	defer span.End()

	sysID := req.SystemIdentifier
	if sysID == 0 {
		v.bcrCount.Add(ctx, 1, metric.WithAttributeSet(attribute.NewSet(attribute.String("status", "err.sysid"))))
		return fmt.Errorf("request contains invalid partition identifier %s", sysID)
	}
	defer func() {
		if rErr != nil {
			span.RecordError(rErr)
			span.SetStatus(codes.Error, rErr.Error())
		}
		partition := observability.Partition(sysID)
		span.SetAttributes(partition)
		v.bcrCount.Add(ctx, 1, metric.WithAttributeSet(attribute.NewSet(observability.ErrStatus(rErr), partition)))
	}()

	pdr, pTrustBase, err := v.partitions.GetInfo(sysID, req.RootRound())
	if err != nil {
		return fmt.Errorf("reading partition info: %w", err)
	}
	if err := pdr.IsValidShard(req.Shard); err != nil {
		return fmt.Errorf("invalid shard: %w", err)
	}
	if err = pTrustBase.Verify(req.NodeIdentifier, req); err != nil {
		return fmt.Errorf("partition %s node %v rejected: %w", sysID, req.NodeIdentifier, err)
	}
	latestUnicityCertificate, err := v.consensusManager.GetLatestUnicityCertificate(sysID)
	if err != nil {
		return fmt.Errorf("reading last certified state: %w", err)
	}
	v.subscription.Subscribe(sysID, req.NodeIdentifier)
	if err = consensus.CheckBlockCertificationRequest(req, latestUnicityCertificate); err != nil {
		err = fmt.Errorf("invalid block certification request: %w", err)
		if se := v.sendResponse(ctx, req.NodeIdentifier, latestUnicityCertificate); se != nil {
			err = errors.Join(err, fmt.Errorf("sending latest cert: %w", se))
		}
		return err
	}
	// check if consensus is already achieved, then store, but it will not be used as proof
	if res := v.incomingRequests.IsConsensusReceived(sysID, pTrustBase); res != QuorumInProgress {
		// stale request buffer, but no need to add extra proof
		if _, _, err = v.incomingRequests.Add(sysID, req, pTrustBase); err != nil {
			return fmt.Errorf("stale block certification request, could not be stored: %w", err)
		}
		return nil
	}
	// store new request and see if quorum is achieved
	res, proof, err := v.incomingRequests.Add(sysID, req, pTrustBase)
	if err != nil {
		return fmt.Errorf("storing request: %w", err)
	}
	var reason consensus.CertReqReason
	switch res {
	case QuorumAchieved:
		v.log.DebugContext(ctx, fmt.Sprintf("partition %s reached consensus, new InputHash: %X", sysID, proof[0].InputRecord.Hash))
		reason = consensus.Quorum
	case QuorumNotPossible:
		v.log.DebugContext(ctx, fmt.Sprintf("partition %s consensus not possible, repeat UC", sysID))
		reason = consensus.QuorumNotPossible
	case QuorumInProgress:
		v.log.DebugContext(ctx, fmt.Sprintf("partition %s quorum not yet reached, but possible in the future", sysID))
		return nil
	}
	if err = v.consensusManager.RequestCertification(ctx,
		consensus.IRChangeRequest{
			SystemIdentifier: sysID,
			Reason:           reason,
			Requests:         proof,
		}); err != nil {
		return fmt.Errorf("requesting certification: %w", err)
	}
	return nil
}

// handleConsensus - receives consensus results and delivers certificates to subscribers
func (v *Node) handleConsensus(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case uc, ok := <-v.consensusManager.CertificationResult():
			if !ok {
				return fmt.Errorf("consensus channel closed")
			}
			v.onCertificationResult(ctx, uc)
		}
	}
}

func (v *Node) onCertificationResult(ctx context.Context, certificate *types.UnicityCertificate) {
	sysID := certificate.UnicityTreeCertificate.SystemIdentifier
	if sysID == 0 {
		v.log.WarnContext(ctx, "certificate has invalid partition id")
		return
	}
	// remember to clear the incoming buffer to accept new nodeRequest
	// NB! this will try and reset the store also in the case when system id is unknown, but this is fine
	defer func() {
		v.incomingRequests.Clear(sysID)
		v.log.LogAttrs(ctx, logger.LevelTrace, fmt.Sprintf("Resetting request store for partition '%s'", sysID))
	}()

	subscribed := v.subscription.Get(sysID)
	v.log.DebugContext(ctx, fmt.Sprintf("sending unicity certificate to partition %s, IR Hash: %X, Block Hash: %X",
		certificate.UnicityTreeCertificate.SystemIdentifier, certificate.InputRecord.Hash, certificate.InputRecord.BlockHash))
	// send response to all registered nodes
	for _, node := range subscribed {
		if err := v.sendResponse(ctx, node, certificate); err != nil {
			v.log.WarnContext(ctx, "sending certification result", logger.Error(err))
		}
		v.subscription.ResponseSent(sysID, node)
	}
}
