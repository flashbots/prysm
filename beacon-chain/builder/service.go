package builder

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/prysmaticlabs/prysm/beacon-chain/powchain"
	"github.com/prysmaticlabs/prysm/network"
	enginev1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
)

var log = logrus.WithField("prefix", "builder")

type config struct {
	currHttpEndpoint network.Endpoint
	httpEndpoints    []network.Endpoint
}

type Caller interface {
	NewPayloadAttributes(ctx context.Context, attrs *enginev1.BuilderPayloadAttributes)
}

type Service struct {
	cfg       *config
	ctx       context.Context
	cancel    context.CancelFunc
	rpcClient powchain.RPCClient
}

func NewService(ctx context.Context, opts ...Option) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop()

	s := &Service{
		cfg:    &config{},
		ctx:    ctx,
		cancel: cancel,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Start the powchain service's main event loop.
func (s *Service) Start() {
	// TODO: open connections to all endpoints and periodically reopen closed ones
	client, err := powchain.NewRPCClientWithAuth(s.ctx, s.cfg.currHttpEndpoint)
	if err != nil {
		log.WithError(err).Error("Could not connect to builder endpoint")
	}
	s.rpcClient = client

	go s.run(s.ctx.Done())
}

// Stop the web3 service's main event loop and associated goroutines.
func (s *Service) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
	}
	if s.rpcClient != nil {
		s.rpcClient.Close()
	}
	return nil
}

func (s *Service) Status() error {
	// FIXME
	return nil
}

func (s *Service) run(done <-chan struct{}) {
	// TODO: open connections to all endpoints and periodically reopen closed ones
	for {
		select {
		case <-done:
			s.rpcClient.Close()
			log.Debug("Context closed, exiting goroutine")
			return
		}
	}
}

func (s *Service) NewPayloadAttributes(ctx context.Context, attrs *enginev1.BuilderPayloadAttributes) {
}
