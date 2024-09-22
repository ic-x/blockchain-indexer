# blockchain-indexer
A CLI tool written in Go to index Ethereum blockchain blocks. It fetches block data via RPC, starting from a specified block, and saves details (block number, hash, transaction count) to a file. Supports concurrent block processing with goroutines
