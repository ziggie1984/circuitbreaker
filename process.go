package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lightningnetwork/lnd/routing/route"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	rpcTimeout = 10 * time.Second
	ctxb       = context.Background()
)

const burstSize = 10

type lndclient interface {
	getIdentity() (route.Vertex, error)

	listChannels() (map[uint64]*channel, error)

	getNodeAlias(key route.Vertex) (string, error)

	subscribeHtlcEvents(ctx context.Context) (htlcEventsClient, error)

	htlcInterceptor(ctx context.Context) (htlcInterceptorClient, error)

	getPendingIncomingHtlcs(ctx context.Context, peer *route.Vertex) (
		map[route.Vertex]map[circuitKey]struct{}, error)
}

type circuitKey struct {
	channel uint64
	htlc    uint64
}

type interceptEvent struct {
	circuitKey
	resume func(bool) error
}

type resolvedEvent struct {
	circuitKey
	settled bool
}

type rateCounters struct {
	counters map[route.Vertex]*peerState
}

type rateCountersRequest struct {
	counters chan *rateCounters
}

type process struct {
	client lndclient
	limits *Limits
	log    *zap.SugaredLogger

	interceptChan           chan interceptEvent
	resolveChan             chan resolvedEvent
	updateLimitChan         chan updateLimitEvent
	rateCountersRequestChan chan rateCountersRequest

	identity route.Vertex
	chanMap  map[uint64]*channel
	aliasMap map[route.Vertex]string

	peerCtrls map[route.Vertex]*peerController

	burstSize int

	// Testing hook
	resolvedCallback func()
}

func NewProcess(client lndclient, log *zap.SugaredLogger, limits *Limits) *process {
	return &process{
		log:                     log,
		client:                  client,
		interceptChan:           make(chan interceptEvent),
		resolveChan:             make(chan resolvedEvent),
		updateLimitChan:         make(chan updateLimitEvent),
		rateCountersRequestChan: make(chan rateCountersRequest),
		chanMap:                 make(map[uint64]*channel),
		aliasMap:                make(map[route.Vertex]string),
		peerCtrls:               make(map[route.Vertex]*peerController),
		limits:                  limits,
		burstSize:               burstSize,
	}
}

type updateLimitEvent struct {
	limit *Limit
	peer  *route.Vertex
}

func (p *process) UpdateLimit(ctx context.Context, peer *route.Vertex,
	limit *Limit) error {

	if peer == nil && limit == nil {
		return errors.New("cannot clear default limit")
	}

	update := updateLimitEvent{
		limit: limit,
		peer:  peer,
	}

	select {
	case p.updateLimitChan <- update:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *process) Run(ctx context.Context) error {
	p.log.Info("CircuitBreaker started")

	var err error

	p.identity, err = p.client.getIdentity()
	if err != nil {
		return err
	}

	p.log.Infow("Connected to lnd node",
		"pubkey", p.identity.String())

	group, ctx := errgroup.WithContext(ctx)

	stream, err := p.client.subscribeHtlcEvents(ctx)
	if err != nil {
		return err
	}

	interceptor, err := p.client.htlcInterceptor(ctx)
	if err != nil {
		return err
	}

	p.log.Info("Interceptor/notification handlers registered")

	group.Go(func() error {
		err := p.processHtlcEvents(ctx, stream)
		if err != nil {
			return fmt.Errorf("htlc events error: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		err := p.processInterceptor(ctx, interceptor)
		if err != nil {
			return fmt.Errorf("interceptor error: %w", err)
		}

		return err
	})

	group.Go(func() error {
		return p.eventLoop(ctx)
	})

	return group.Wait()
}

func (p *process) getPeerController(ctx context.Context, peer route.Vertex,
	startGo func(func() error)) *peerController {

	ctrl, ok := p.peerCtrls[peer]
	if ok {
		return ctrl
	}

	// If the peer does not yet exist, initialize it with no pending htlcs.
	htlcs := make(map[circuitKey]struct{})

	return p.createPeerController(ctx, peer, startGo, htlcs)
}

func (p *process) createPeerController(ctx context.Context, peer route.Vertex,
	startGo func(func() error), htlcs map[circuitKey]struct{}) *peerController {

	peerCfg, ok := p.limits.PerPeer[peer]
	if !ok {
		peerCfg = p.limits.Default
	}

	cfg := &peerControllerCfg{
		logger:    p.log,
		limit:     peerCfg,
		burstSize: p.burstSize,
		htlcs:     htlcs,
		lnd:       p.client,
		pubKey:    peer,
	}
	ctrl := newPeerController(cfg)

	startGo(func() error {
		return ctrl.run(ctx)
	})

	p.peerCtrls[peer] = ctrl

	return ctrl
}

func (p *process) eventLoop(ctx context.Context) error {
	// Create a group to attach peer goroutines to.
	group, ctx := errgroup.WithContext(ctx)
	defer func() {
		_ = group.Wait()
	}()

	// Retrieve all pending htlcs from lnd.
	htlcsPerPeer, err := p.client.getPendingIncomingHtlcs(ctx, nil)
	if err != nil {
		return err
	}

	// Initialize peer controllers with currently pending htlcs.
	for peer, htlcs := range htlcsPerPeer {
		p.createPeerController(ctx, peer, group.Go, htlcs)
	}

	for {
		select {
		case interceptEvent := <-p.interceptChan:
			chanInfo, err := p.getChanInfo(interceptEvent.channel)
			if err != nil {
				return err
			}

			ctrl := p.getPeerController(ctx, chanInfo.peer, group.Go)

			peerEvent := peerInterceptEvent{
				interceptEvent: interceptEvent,
				peerInitiated:  !chanInfo.initiator,
			}
			if err := ctrl.process(ctx, peerEvent); err != nil {
				return err
			}

		case resolvedEvent := <-p.resolveChan:
			chanInfo, err := p.getChanInfo(resolvedEvent.channel)
			if err != nil {
				return err
			}

			ctrl := p.getPeerController(ctx, chanInfo.peer, group.Go)

			if err := ctrl.resolved(ctx, resolvedEvent); err != nil {
				return err
			}

			if p.resolvedCallback != nil {
				p.resolvedCallback()
			}

		case update := <-p.updateLimitChan:
			switch {
			// Update sets default limit.
			case update.peer == nil:
				p.limits.Default = *update.limit

				// Update all controllers that have no specific limit.
				for node, ctrl := range p.peerCtrls {
					_, ok := p.limits.PerPeer[node]
					if ok {
						continue
					}

					err := ctrl.updateLimit(ctx, *update.limit)
					if err != nil {
						return err
					}
				}

			// Update sets specific limit.
			case update.limit != nil:
				p.limits.PerPeer[*update.peer] = *update.limit

				// Update specific controller if it exists.
				ctrl, ok := p.peerCtrls[*update.peer]
				if ok {
					err := ctrl.updateLimit(ctx, *update.limit)
					if err != nil {
						return err
					}
				}

			// Update clears limit.
			case update.limit == nil:
				delete(p.limits.PerPeer, *update.peer)

				// Apply default limit to peer controller.
				ctrl, ok := p.peerCtrls[*update.peer]
				if ok {
					err := ctrl.updateLimit(ctx, p.limits.Default)
					if err != nil {
						return err
					}
				}
			}

		case req := <-p.rateCountersRequestChan:
			allCounts := make(map[route.Vertex]*peerState)
			for node, ctrl := range p.peerCtrls {
				state, err := ctrl.state(ctx)
				if err != nil {
					return err
				}

				allCounts[node] = state
			}

			req.counters <- &rateCounters{
				counters: allCounts,
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *process) getRateCounters(ctx context.Context) (
	map[route.Vertex]*peerState, error) {

	replyChan := make(chan *rateCounters)

	select {
	case p.rateCountersRequestChan <- rateCountersRequest{
		counters: replyChan,
	}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case reply := <-replyChan:
		return reply.counters, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *process) processHtlcEvents(ctx context.Context,
	stream htlcEventsClient) error {

	for {
		event, err := stream.recv()
		if err != nil {
			return err
		}

		select {
		case p.resolveChan <- *event:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *process) processInterceptor(ctx context.Context,
	interceptor htlcInterceptorClient) error {

	for {
		event, err := interceptor.recv()
		if err != nil {
			return err
		}

		key := event.circuitKey

		resume := func(resume bool) error {
			return interceptor.send(&interceptResponse{
				key:    key,
				resume: resume,
			})
		}

		select {
		case p.interceptChan <- interceptEvent{
			circuitKey: key,
			resume:     resume,
		}:

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *process) getChanInfo(channel uint64) (*channel, error) {
	// Try to look up from the cache.
	ch, ok := p.chanMap[channel]
	if ok {
		return ch, nil
	}

	// Cache miss. Retrieve all channels and update the cache.
	channels, err := p.client.listChannels()
	if err != nil {
		return nil, err
	}

	for chanId, ch := range channels {
		p.chanMap[chanId] = ch
	}

	// Try looking up the channel again.
	ch, ok = p.chanMap[channel]
	if ok {
		return ch, nil
	}

	// Channel not found.
	return nil, fmt.Errorf("incoming channel %v not found", channel)
}
