package storage

import (
	"github.com/ic-x/blockchain-indexer/internal/blockchain"
)

type Storage interface {
	SaveBlock(data *blockchain.BlockData) error
}
