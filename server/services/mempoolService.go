package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type mempoolService struct {
	provider       NetworkProvider
	errFactory     *errFactory
	txsTransformer *transactionsTransformer
}

// NewMempoolService will create a new instance of mempoolAPIService
func NewMempoolService(provider NetworkProvider) server.MempoolAPIServicer {
	return &mempoolService{
		provider:       provider,
		errFactory:     newErrFactory(),
		txsTransformer: newTransactionsTransformer(provider),
	}
}

// Mempool is not implemented yet
func (service *mempoolService) Mempool(context.Context, *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	return nil, newErrFactory().newErr(ErrNotImplemented)
}

// MempoolTransaction will return operations for a transaction that is in pool
func (service *mempoolService) MempoolTransaction(
	_ context.Context,
	request *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	tx, err := service.provider.GetMempoolTransactionByHash(request.TransactionIdentifier.Hash)
	if err != nil {
		return nil, service.errFactory.newErrWithOriginal(ErrCannotParsePoolTransaction, err)
	}
	if tx == nil {
		return nil, service.errFactory.newErr(ErrTransactionIsNotInPool)
	}

	rosettaTx := service.txsTransformer.mempoolMoveBalanceTxToRosettaTx(tx)
	if err != nil {
		return nil, service.errFactory.newErr(ErrCannotParsePoolTransaction)
	}

	return &types.MempoolTransactionResponse{
		Transaction: rosettaTx,
	}, nil
}
