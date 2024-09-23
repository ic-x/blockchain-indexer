package blockchain

import (
	"context"
)

type Blockchain interface {
	GetBlockByNumber(ctx context.Context, number int64) (*BlockData, error)
	GetBlockByHash(ctx context.Context, hash string) (*BlockData, error)
	GetLatestBlockNumber(ctx context.Context) (int64, error)
	SubscribeNewBlocks(ctx context.Context, headers chan<- *BlockHeader) (Subscription, error)
}

type BlockData struct {
	Number       int64  // Block number
	Hash         string // Block hash
	TxCount      int    // Number of transactions in the block
	Timestamp    int64  // Block timestamp
	ParentHash   string // Parent block hash
	Nonce        uint64 // Nonce value for block mining
	Miner        string // Address of the miner who mined the block
	GasUsed      uint64 // Gas used in the block
	GasLimit     uint64 // Gas limit for the block
	Size         uint64 // Block size in bytes
	ExtraData    []byte // Additional data in the block
	Difficulty   uint64 // Mining difficulty at the time the block was created
	ReceiptsRoot string // Root of the transaction receipts tree
}

type BlockHeader struct {
	Hash   string
	Number int64
}

// Subscription - interface for subscribing to new blocks.
type Subscription interface {
	Unsubscribe()
	Err() <-chan error
}
