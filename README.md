# Blockchain Indexer

## Overview

**Blockchain Indexer** is a CLI tool written in Go for indexing Ethereum blockchain blocks. It fetches block data via RPC, starting from a specified block, and stores the details (block number, hash, transaction count) in a file. The tool supports concurrent block processing with goroutines for enhanced performance.

## Configuration

The configuration can be set via a `config.yaml` file or through command-line flags. Below is an example configuration file:

```yaml
rpc: "https://holesky.infura.io/v3/<YOUR_API_KEY>"
start: 1
out: "blocks.log"
```

### Configurable Parameters

- `rpc` : RPC URL to connect to an Ethereum node (required).
- `start` : Block number to start parsing from (default: 0).
- `end` : Block number to stop parsing (default: -1, meaning no end).
- `live`: If set to `true`, the indexer will continuously process newly mined blocks (default: `false`).
- `allow_future_start` : Allows the indexer to start from future blocks that haven't been mined yet (default: `false`).
- `retry_interval`: Interval (in seconds) to retry failed requests (default: 2 seconds).
- `block_buffer_size` : Buffer size for the block processing channel (default: 0, unbuffered).
- `headers_buffer_size` : Buffer size for the header processing channel (default: 0, unbuffered).

### Command-Line Flags

The following flags are available for configuring the tool directly from the command line:

- `--rpc` : RPC URL to connect to an Ethereum node (required if not set in `config.yaml`).
- `--start` : Block number to start parsing from.
- `--out` : Output file to store the block details (e.g., `blocks.log`).
- `--end` : Block number to stop parsing (default: -1, meaning no end).
- `--live` : If set, the indexer will continue to process blocks as they are mined.
- `--allow-future-start` : Allows starting from blocks that are not yet mined.

## How to Run

### Running Locally (Without Docker)

To run the application locally, you can use command-line flags or a configuration file:

```bash
go run main.go run --rpc=https://holesky.infura.io/v3/<YOUR_API_KEY> --start=1 --out=blocks.log
```

Alternatively, use WebSocket for RPC:

```bash
go run main.go run --rpc=wss://holesky.infura.io/ws/v3/<YOUR_API_KEY> --start=1 --out=blocks.log
```

### Running with Docker Compose

To run the application using Docker Compose, use the following commands:

```bash
docker-compose up --build
docker-compose up
```

## Technologies Used

- **Go-Ethereum (geth)**: Ethereum client libraries for interacting with Ethereum nodes.
- **Cobra**: Library for building CLI applications.
- **Viper**: Library for configuration management via `config.yaml`.

## Example RPC URLs

You can use the following RPC URLs to connect to the Holesky Ethereum testnet:

- **HTTP RPC**: `https://holesky.infura.io/v3/<YOUR_API_KEY>`
- **WebSocket RPC**: `wss://holesky.infura.io/ws/v3/<YOUR_API_KEY>`

These URLs are provided by Infura and can be used to connect to the Holesky test network.
