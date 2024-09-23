package storage

import (
    "fmt"
    "os"

    "github.com/ic-x/blockchain-indexer/internal/blockchain"
)

type FileStorage struct {
    file *os.File
}

func NewFileStorage(file *os.File) *FileStorage {
    return &FileStorage{file: file}
}

func (s *FileStorage) SaveBlock(data *blockchain.BlockData) error {
    _, err := s.file.WriteString(fmt.Sprintf(
        "Number: %d\nHash: %s\nTxCount: %d\nTimestamp: %d\nParentHash: %s\nNonce: %d\nMiner: %s\nGasUsed: %d\nGasLimit: %d\nSize: %d\nExtraData: %x\nDifficulty: %d\nReceiptsRoot: %s\n\n",
        data.Number, data.Hash, data.TxCount, data.Timestamp, data.ParentHash, data.Nonce, data.Miner, data.GasUsed, data.GasLimit, data.Size, data.ExtraData, data.Difficulty, data.ReceiptsRoot))
    return err
}
