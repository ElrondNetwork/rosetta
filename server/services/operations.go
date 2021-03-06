package services

import "github.com/coinbase/rosetta-sdk-go/types"

const (
	opGenesisBalanceMovement = "GenesisBalanceMovement"
	opTransfer               = "Transfer"
	opFee                    = "Fee"
	opReward                 = "Reward"
	opScResult               = "SmartContractResult"
	opFeeOfInvalidTx         = "FeeOfInvalidTransaction"
	opFeeRefund              = "FeeRefund"
)

var (
	// SupportedOperationTypes is a list of the supported operations.
	SupportedOperationTypes = []string{
		opTransfer,
		opFee,
		opReward,
		opScResult,
		opFeeOfInvalidTx,
		opGenesisBalanceMovement,
		opFeeRefund,
	}

	opStatusSuccess = "Success"

	supportedOperationStatuses = []*types.OperationStatus{
		{
			Status:     opStatusSuccess,
			Successful: true,
		},
	}
)

func indexOperations(operations []*types.Operation) {
	for index, operation := range operations {
		operation.OperationIdentifier = indexToOperationIdentifier(index)
	}
}

func populateStatusOfOperations(operations []*types.Operation) {
	for _, operation := range operations {
		// TODO: Improve this, perhaps use a clone?
		operation.Status = &opStatusSuccess
	}
}
