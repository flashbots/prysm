package builder

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/prysmaticlabs/prysm/consensus-types/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

var log = logrus.WithField("prefix", "builder")

type config struct {
}

type Caller interface {
	OnNewPayload(ctx context.Context, postStateHeader *ethpb.ExecutionPayloadHeader, blk interfaces.SignedBeaconBlock)
}

type Service struct {
	cfg    *config
	ctx    context.Context
	cancel context.CancelFunc
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
	go s.run(s.ctx.Done())
}

// Stop the web3 service's main event loop and associated goroutines.
func (s *Service) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
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
			log.Debug("Context closed, exiting goroutine")
			return
		}
	}
}

func (s *Service) OnNewPayload(ctx context.Context, postStateHeader *ethpb.ExecutionPayloadHeader, blk interfaces.SignedBeaconBlock) {
	// TODO: adjust the attributes to what's actually required
	_ = &enginev1.BuilderPayloadAttributes{
		Timestamp:  postStateHeader.Timestamp,
		Slot:       uint64(blk.Block().Slot()),
		PrevRandao: blk.Block().Body().RandaoReveal(),
	}
}
