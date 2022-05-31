package builder

import (
	"github.com/prysmaticlabs/prysm/beacon-chain/powchain"
	"github.com/prysmaticlabs/prysm/network"
)

type Option func(s *Service) error

func WithHttpEndpoints(endpointStrings []string) Option {
	return func(s *Service) error {
		stringEndpoints := powchain.DedupEndpoints(endpointStrings)
		endpoints := make([]network.Endpoint, len(stringEndpoints))
		for i, e := range stringEndpoints {
			endpoints[i] = powchain.HttpEndpoint(e)
		}
		// Select first http endpoint in the provided list.
		var currEndpoint network.Endpoint
		if len(endpointStrings) > 0 {
			currEndpoint = endpoints[0]
		}
		s.cfg.httpEndpoints = endpoints
		s.cfg.currHttpEndpoint = currEndpoint
		return nil
	}
}
