package provider

import (
	"testing"

	dataBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestGatherInvalidTransactions(t *testing.T) {
	// Block N-1
	previousBlock := &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				Type: dataBlock.InvalidBlock.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "aaaa"},
				},
			},
		},
	}

	// Block N
	block := &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				ProcessingType: dataBlock.Scheduled.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "bbbb"},
					{Hash: "cccc"},
				},
			},
			{
				Type: dataBlock.InvalidBlock.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "bbbb"},
				},
			},
		},
	}

	// Block N+1
	nextBlock := &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				ProcessingType: dataBlock.Processed.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "cccc"},
				},
			},
			{
				Type: dataBlock.InvalidBlock.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "bbbb"},
					{Hash: "eeee"},
				},
			},
		},
	}

	// "aaaa" is ignored for N, since it produces effects in N-1
	// "eeee" is ignored for N, since it produces effects in N+1
	invalidTxs := gatherInvalidTransactions(previousBlock, block, nextBlock)
	require.Len(t, invalidTxs, 1)
	require.Equal(t, "bbbb", invalidTxs[0].Hash)
}

func TestGatherInvalidTransactions_WhenIntraShardIsMissingInPreviousBlock(t *testing.T) {
	// Edge case on start of epoch.

	// Block N-1
	previousBlock := &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				ProcessingType: dataBlock.Scheduled.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "aaaa"},
				},
			},
			// "aaaa" is an invalid transaction that produces effects in N-1,
			// but the intra-shard, "invalid" miniblock is missing (not provided by the API).
		},
	}

	// Block N
	block := &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				ProcessingType: dataBlock.Scheduled.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "abab"},
					{Hash: "cccc"},
				},
			},
			{
				Type: dataBlock.InvalidBlock.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "aaaa"},
				},
			},
			{
				// Intra-shard miniblock, holds both "aaaa" (scheduled in N-1, with effects in N)
				// and "abab" (scheduled in N, with effects in N)
				Type: dataBlock.InvalidBlock.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "aaaa"},
					{Hash: "abab"},
				},
			},
		},
	}

	// Block N+1
	nextBlock := &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				ProcessingType: dataBlock.Processed.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "cccc"},
				},
			},
			{
				Type: dataBlock.InvalidBlock.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "abab"},
				},
			},
		},
	}

	// "aaaa" is ignored for N, since it produces effects in N-1
	invalidTxs := gatherInvalidTransactions(previousBlock, block, nextBlock)
	require.Len(t, invalidTxs, 1)
	require.Equal(t, "abab", invalidTxs[0].Hash)
}

func TestDoSimplifyBlockWithScheduledTransactions_WithRespectToConstructionState(t *testing.T) {
	// Edge case on cross-shard miniblocks, both scheduled and final.

	// Empty, trivial blocks at N-1 and N+1
	previousBlock := &data.Block{MiniBlocks: []*data.MiniBlock{}}
	nextBlock := &data.Block{MiniBlocks: []*data.MiniBlock{}}

	// Scheduled & Final (won't be removed)
	block := &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				ProcessingType:    dataBlock.Scheduled.String(),
				ConstructionState: dataBlock.Final.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "aaaa"},
					{Hash: "bbbb"},
				},
			},
		},
	}

	doSimplifyBlockWithScheduledTransactions(previousBlock, block, nextBlock)
	require.Len(t, block.MiniBlocks, 1)
	require.Len(t, block.MiniBlocks[0].Transactions, 2)
	require.Equal(t, "aaaa", block.MiniBlocks[0].Transactions[0].Hash)
	require.Equal(t, "bbbb", block.MiniBlocks[0].Transactions[1].Hash)

	// Scheduled & !Final (will be removed)
	block = &data.Block{
		MiniBlocks: []*data.MiniBlock{
			{
				ProcessingType: dataBlock.Scheduled.String(),
				Transactions: []*data.FullTransaction{
					{Hash: "aaaa"},
					{Hash: "bbbb"},
				},
			},
		},
	}

	doSimplifyBlockWithScheduledTransactions(previousBlock, block, nextBlock)
	require.Len(t, block.MiniBlocks, 0)
}
