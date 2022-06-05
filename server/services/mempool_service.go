package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type mempoolAPIService struct {
	provider  NetworkProvider
	txsParser *transactionsParser
}

// NewMempoolApiService will create a new instance of mempoolAPIService
func NewMempoolApiService(provider NetworkProvider) server.MempoolAPIServicer {
	return &mempoolAPIService{
		provider:  provider,
		txsParser: newTransactionParser(provider),
	}
}

// Mempool is not implemented yet
func (mas *mempoolAPIService) Mempool(context.Context, *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	return nil, ErrNotImplemented
}

// MempoolTransaction will return operations for a transaction that is in pool
func (mas *mempoolAPIService) MempoolTransaction(
	_ context.Context,
	request *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	tx, ok := mas.provider.GetTransactionByHashFromPool(request.TransactionIdentifier.Hash)
	if !ok {
		return nil, ErrTransactionIsNotInPool
	}

	rosettaTx, err := mas.txsParser.parseTx(tx, true)
	if err != nil {
		return nil, ErrCannotParsePoolTransaction
	}

	return &types.MempoolTransactionResponse{
		Transaction: rosettaTx,
	}, nil

}
