package blockchain

import (
    "context"
    "math/big"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
    client *ethclient.Client
}

func NewEthereum(client *ethclient.Client) *Ethereum {
    return &Ethereum{client: client}
}

func (e *Ethereum) GetBlockByNumber(ctx context.Context, number int64) (*BlockData, error) {
    block, err := e.client.BlockByNumber(ctx, big.NewInt(number))
    if err != nil {
        return nil, err
    }

    return &BlockData{
        Number:       block.Number().Int64(),
        Hash:         block.Hash().Hex(),
        TxCount:      len(block.Transactions()),
        Timestamp:    int64(block.Time()),
        ParentHash:   block.ParentHash().Hex(),
        Nonce:        block.Nonce(),
        Miner:        block.Coinbase().Hex(),
        GasUsed:      block.GasUsed(),
        GasLimit:     block.GasLimit(),
        Size:         block.Size(),
        ExtraData:    block.Extra(),
        Difficulty:   block.Difficulty().Uint64(),
        ReceiptsRoot: block.ReceiptHash().Hex(),
    }, nil
}

func (e *Ethereum) GetBlockByHash(ctx context.Context, hash string) (*BlockData, error) {
    blockHash := common.HexToHash(hash)
    block, err := e.client.BlockByHash(ctx, blockHash)
    if err != nil {
        return nil, err
    }

    return &BlockData{
        Number:       block.Number().Int64(),
        Hash:         block.Hash().Hex(),
        TxCount:      len(block.Transactions()),
        Timestamp:    int64(block.Time()),
        ParentHash:   block.ParentHash().Hex(),
        Nonce:        block.Nonce(),
        Miner:        block.Coinbase().Hex(),
        GasUsed:      block.GasUsed(),
        GasLimit:     block.GasLimit(),
        Size:         block.Size(),
        ExtraData:    block.Extra(),
        Difficulty:   block.Difficulty().Uint64(),
        ReceiptsRoot: block.ReceiptHash().Hex(),
    }, nil
}

func (e *Ethereum) GetLatestBlockNumber(ctx context.Context) (int64, error) {
    header, err := e.client.HeaderByNumber(ctx, nil)
    if err != nil {
        return 0, err
    }
    return header.Number.Int64(), nil
}

func (e *Ethereum) SubscribeNewBlocks(ctx context.Context, headers chan<- *BlockHeader) (Subscription, error) {
    ethHeaders := make(chan *types.Header)
    sub, err := e.client.SubscribeNewHead(ctx, ethHeaders)
    if err != nil {
        return nil, err
    }

    go func() {
        for ethHeader := range ethHeaders {
            headers <- &BlockHeader{
                Hash:   ethHeader.Hash().Hex(),
                Number: ethHeader.Number.Int64(),
            }
        }
    }()

    return &EthereumSubscription{sub: sub}, nil
}

// EthereumSubscription - structure for subscription management.
type EthereumSubscription struct {
    sub ethereum.Subscription
}

// Unsubscribe - unsubscribe from receiving new blocks.
func (s *EthereumSubscription) Unsubscribe() {
    s.sub.Unsubscribe()
}

// Err - receive error channel.
func (s *EthereumSubscription) Err() <-chan error {
    return s.sub.Err()
}
