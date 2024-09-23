package run

import (
    "context"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ic-x/blockchain-indexer/internal/blockchain"
    "github.com/ic-x/blockchain-indexer/internal/storage"
    "github.com/ic-x/blockchain-indexer/internal/worker"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

func determineLinkType(link string) string {
    switch {
    case strings.HasPrefix(link, "https://"):
        return "HTTPS"
    case strings.HasPrefix(link, "http://"):
        return "HTTP"
    case strings.HasPrefix(link, "wss://"):
        return "WebSocket Secure"
    case strings.HasPrefix(link, "ws://"):
        return "WebSocket"
    default:
        return "Unknown"
    }
}

func validateParams(retryInterval, blockBufferSize, headersBufferSize, startBlock, endBlock int64) error {
    if retryInterval < 0 {
        return fmt.Errorf("retryInterval must be >= 0, got %d", retryInterval)
    }
    if blockBufferSize < 0 {
        return fmt.Errorf("blockBufferSize must be >= 0, got %d", blockBufferSize)
    }
    if headersBufferSize < 0 {
        return fmt.Errorf("headersBufferSize must be >= 0, got %d", headersBufferSize)
    }
    if startBlock < 0 {
        return fmt.Errorf("startBlock must be >= 0, got %d", startBlock)
    }
    if endBlock > 0 && startBlock > endBlock {
        return fmt.Errorf("startBlock must be <= endBlock, got startBlock = %d and endBlock = %d", startBlock, endBlock)
    }
    return nil
}

var RunCmd = &cobra.Command{
    Use:   "run",
    Short: "Start the blockchain indexer",
    Run: func(cmd *cobra.Command, args []string) {
        // Extracting command line flags
        rpcURL, err := cmd.Flags().GetString("rpc")
        if err != nil {
            log.Fatalf("Error reading --rpc flag: %v", err)
        }
        startBlock, err := cmd.Flags().GetInt64("start")
        if err != nil {
            log.Fatalf("Error reading --start flag: %v", err)
        }
        endBlock, err := cmd.Flags().GetInt64("end")
        if err != nil {
            log.Fatalf("Error reading --end flag: %v", err)
        }
        live, err := cmd.Flags().GetBool("live")
        if err != nil {
            log.Fatalf("Error reading --live flag: %v", err)
        }
        allowFutureStart, err := cmd.Flags().GetBool("allow-future-start")
        if err != nil {
            log.Fatalf("Error reading --allow-future-start flag: %v", err)
        }
    
        // Loading parameters from config file (if flags are not provided)
        if rpcURL == "" {
            rpcURL = viper.GetString("rpc")
            if rpcURL == "" {
                log.Fatal("RPC URL is required (set in config.yaml or use --rpc flag)")
            }
        }
    
        retryInterval := int64(viper.GetInt("retry_interval"))
        blockBufferSize := int64(viper.GetInt("block_buffer_size"))
        headersBufferSize := int64(viper.GetInt("headers_buffer_size"))
        outputFile := viper.GetString("out")
    
        // Checking for conflicting parameters
        if live && (startBlock != 0 || allowFutureStart) {
            log.Fatal("When 'live' mode is enabled, 'start' and 'allow-future-start' cannot be set")
        }

        // Validate parameters using the new function
        if err := validateParams(retryInterval, blockBufferSize, headersBufferSize, startBlock, endBlock); err != nil {
            log.Fatalf("Parameter validation failed: %v", err)
        }
    
        // Connecting to RPC
        client, err := ethclient.Dial(rpcURL)
        if err != nil {
            log.Fatalf("Failed to connect to the Ethereum client: %v", err)
        }
    
        // Opening the file for logging (using the value from the 'out' flag)
        file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatalf("Failed to open output file: %v", err)
        }
        defer file.Close()
    
        // Creating storage and blockchain instances
        storage := storage.NewFileStorage(file)
        eth := blockchain.NewEthereum(client)
    
        // Starting the worker to process blocks
        ctx := context.Background()
        worker := worker.NewBlockWorker(eth, storage, int(blockBufferSize), int(headersBufferSize))
    
        linkType := determineLinkType(rpcURL)
        switch linkType {
        case "HTTPS", "HTTP":
            if live {
                worker.StartLive(ctx, endBlock, int(retryInterval))
            } else {
                worker.Start(ctx, startBlock, allowFutureStart, endBlock, int(retryInterval))
            }
        case "WebSocket", "WebSocket Secure":
            if live {
                worker.StartLiveWithSubscription(ctx, endBlock)
            } else {
                worker.StartWithSubscription(ctx, startBlock, allowFutureStart, endBlock)
            }
        default:
            log.Fatalf("Unknown link type: %s. Supported types are HTTP, HTTPS, ws://, and wss://", linkType)
        }
    },
}

func init() {
    // Defining command line flags
    RunCmd.Flags().String("rpc", "", "RPC URL (required)")
    RunCmd.Flags().Int64("start", 0, "Start block")
    RunCmd.Flags().Int64("end", -1, "End block")
    RunCmd.Flags().Bool("live", false, "Live mode")
    RunCmd.Flags().Bool("allow-future-start", false, "Allow future start block")
    RunCmd.Flags().String("out", "blocks.log", "Output file for storing the logs")

    // Binding flags to Viper for config.yaml
    viper.BindPFlag("rpc", RunCmd.Flags().Lookup("rpc"))
    viper.BindPFlag("start", RunCmd.Flags().Lookup("start"))
    viper.BindPFlag("end", RunCmd.Flags().Lookup("end"))
    viper.BindPFlag("live", RunCmd.Flags().Lookup("live"))
    viper.BindPFlag("allow_future_start", RunCmd.Flags().Lookup("allow-future-start"))
    viper.BindPFlag("out", RunCmd.Flags().Lookup("out"))
}
