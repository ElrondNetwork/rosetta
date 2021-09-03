package provider

import (
	"github.com/ElrondNetwork/elrond-proxy-go/api"
)

func NewOfflineElrondProvider(elrondFacade api.ElrondProxyHandler, networkConfig *NetworkConfig) (*ElrondProvider, error) {
	elrondProxy, ok := elrondFacade.(ElrondProxyClient)
	if !ok {
		return nil, ErrInvalidElrondProxyHandler
	}

	return &ElrondProvider{
		client:                    elrondProxy,
		genesisTime:               networkConfig.StartTime,
		roundDurationMilliseconds: networkConfig.RoundDuration,
	}, nil
}
