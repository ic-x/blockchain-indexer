package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ic-x/blockchain-indexer/internal/blockchain"
	"github.com/ic-x/blockchain-indexer/internal/storage"
)

type BlockWorker struct {
	blockchain        blockchain.Blockchain
	storage           storage.Storage
	blockBufferSize   int
	headersBufferSize int
}

func NewBlockWorker(bc blockchain.Blockchain, st storage.Storage, blockBufferSize, headersBufferSize int) *BlockWorker {
	return &BlockWorker{
		blockchain:        bc,
		storage:           st,
		blockBufferSize:   blockBufferSize,
		headersBufferSize: headersBufferSize,
	}
}

func (w *BlockWorker) saveBlocks(blockCh chan *blockchain.BlockData) {
	for block := range blockCh {
		err := w.storage.SaveBlock(block)
		if err != nil {
			log.Printf("Failed to save block %d: %v", block.Number, err)
		}
	}
}

// checkStartBlock - checks and waits until startBlock <= latestBlock or returns an error
func (w *BlockWorker) checkStartBlock(ctx context.Context, startBlock int64, allowFutureStart bool) error {
	for {
		latestBlockNumber, err := w.blockchain.GetLatestBlockNumber(ctx)
		if err != nil {
			return err
		}

		if startBlock <= latestBlockNumber {
			return nil // Condition met, ready to start
		}

		if !allowFutureStart {
			return fmt.Errorf("start block %d is in the future (latest block: %d) and allowFutureStart is false", startBlock, latestBlockNumber)
		}

		log.Printf("Start block %d is in the future (latest block: %d), waiting...", startBlock, latestBlockNumber)
		time.Sleep(5 * time.Second) // TODO: Refactor me
	}
}

func (w *BlockWorker) Start(ctx context.Context, startBlock int64, allowFutureStart bool, endBlock int64, retryInterval int) {
	// Check the starting block
	if err := w.checkStartBlock(ctx, startBlock, allowFutureStart); err != nil {
		log.Fatalf("Failed to start: %v", err)
		return
	}

	var wg sync.WaitGroup
	blockCh := make(chan *blockchain.BlockData, w.blockBufferSize)

	// Start a goroutine for saving blocks
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.saveBlocks(blockCh)
	}()

	// Parse blocks in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for blockNumber := startBlock; ; blockNumber++ {
			if endBlock > 0 && blockNumber > endBlock {
				log.Printf("Reached end block %d, stopping...", endBlock)
				break
			}

			for {
				block, err := w.blockchain.GetBlockByNumber(ctx, blockNumber)
				if err != nil {
					log.Printf("Failed to fetch block %d: %v. Retrying in %d seconds...", blockNumber, err, retryInterval)
					time.Sleep(time.Duration(retryInterval) * time.Second)
					continue
				}

				blockCh <- block
				log.Printf("Processed block %d", blockNumber)
				break
			}
		}
	}()

	wg.Wait()
	close(blockCh)
}

// StartWithSubscription - processes blocks starting from startBlock, then subscribes to new ones
func (w *BlockWorker) StartWithSubscription(ctx context.Context, startBlock int64, allowFutureStart bool, endBlock int64) {
	// Check the starting block
	if err := w.checkStartBlock(ctx, startBlock, allowFutureStart); err != nil {
		log.Fatalf("Failed to start with subscription: %v", err)
		return
	}

	var wg sync.WaitGroup
	blockCh := make(chan *blockchain.BlockData, w.blockBufferSize)
	headers := make(chan *blockchain.BlockHeader, w.headersBufferSize)
	lastProcessedBlock := startBlock - 1

	// Goroutine for saving blocks
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.saveBlocks(blockCh)
	}()

	// Goroutine for sequentially processing blocks
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Process blocks starting from startBlock
		for blockNumber := startBlock; ; blockNumber++ {
			if endBlock > 0 && blockNumber > endBlock {
				log.Printf("Reached end block: %d", endBlock)
				return
			}

			block, err := w.blockchain.GetBlockByNumber(ctx, blockNumber)
			if err != nil {
				log.Printf("Failed to fetch block %d: %v", blockNumber, err)
				break
			}

			blockCh <- block
			lastProcessedBlock = blockNumber
			log.Printf("Processed block %d", blockNumber)
		}

		// Subscribe to new blocks
		sub, err := w.blockchain.SubscribeNewBlocks(ctx, headers)
		if err != nil {
			log.Fatalf("Failed to subscribe to new blocks: %v", err)
		}
		defer sub.Unsubscribe()

		for {
			select {
			case err := <-sub.Err():
				log.Fatalf("Subscription error: %v", err)
			case header := <-headers:
				newBlockNumber := header.Number

				// Check if we missed any blocks
				if newBlockNumber > lastProcessedBlock+1 {
					for blockNumber := lastProcessedBlock + 1; blockNumber < newBlockNumber; blockNumber++ {
						block, err := w.blockchain.GetBlockByNumber(ctx, blockNumber)
						if err != nil {
							log.Fatalf("Failed to fetch missed block %d: %v", blockNumber, err)
							continue
						}

						blockCh <- block
						lastProcessedBlock = blockNumber
						log.Printf("Processed missed block %d", blockNumber)
					}
				}

				// Process the new block
				block, err := w.blockchain.GetBlockByHash(ctx, header.Hash)
				if err != nil {
					log.Fatalf("Failed to fetch new block by hash: %v", err)
					continue
				}

				blockCh <- block
				lastProcessedBlock = newBlockNumber
				log.Printf("Processed new block %d", newBlockNumber)
			}
		}
	}()

	wg.Wait()
	close(blockCh)
}

func (w *BlockWorker) StartLive(ctx context.Context, endBlock int64, retryInterval int) {
	latestBlockNumber, err := w.blockchain.GetLatestBlockNumber(ctx)
	if err != nil {
		log.Printf("Failed to get the latest block number: %v", err)
		return
	}

	w.Start(ctx, latestBlockNumber, false, endBlock, retryInterval)
}

func (w *BlockWorker) StartLiveWithSubscription(ctx context.Context, endBlock int64) {
	latestBlockNumber, err := w.blockchain.GetLatestBlockNumber(ctx)
	if err != nil {
		log.Printf("Failed to get the latest block number: %v", err)
		return
	}
	w.StartWithSubscription(ctx, latestBlockNumber, false, endBlock)
}
