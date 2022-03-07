package main

import (
    "context"
    "fmt"
    "load-balancer/internal/config"
    "load-balancer/internal/helper/context_helper"
    "load-balancer/internal/helper/log_helper"
    "load-balancer/internal/load_balancer"
    "net"
    "os"

    "github.com/spf13/cobra"
)

func main() {
    rootCmd := cobra.Command{
        Use:          "load-balancer",
        SilenceUsage: true,
    }

    runCmd := &cobra.Command{
        Use:     "run",
        Short:   "Load balancer",
        Long:    "Load balancer",
        Version: config.Version,
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, cancel := context_helper.GetInterruptableCtx(context.Background())
            defer cancel()

            listener, err := net.Listen("tcp", config.LbListenAddr)
            if err != nil {
                log_helper.LogError(ctx, "Failed to listen to address: %s. Error: %v", config.LbListenAddr, err)
                return err
            }
            log_helper.LogInfo(ctx, "Listening to address: %s", config.LbListenAddr)

            err = load_balancer.Run(ctx, listener)
            if err != nil {
                log_helper.LogError(ctx, "Failed to run load balancer. Error: %v", err)
                return err
            }

            return nil
        },
    }

    rootCmd.AddCommand(runCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Printf("Failed to run the load balancer. Error: %v\n", err)
        os.Exit(-1)
    }
}
