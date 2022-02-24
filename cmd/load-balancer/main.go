package main

import (
    "context"
    "fmt"
    "load-balancer/internal/config"
    "load-balancer/internal/context_helper"
    "load-balancer/internal/log_helper"
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

            _, err := net.Listen("tcp", config.ListenAddr)
            if err != nil {
                log_helper.LogError(ctx, "Failed to listen to address: %s. Error is: %v", config.ListenAddr, err)
                return err
            }
            log_helper.LogInfo(ctx, "Listening to address: %s", config.ListenAddr)

            //// now, we can bind to hc port. we prefer this before connecting and
            //// advertising ourselves to switchboard.
            //hcLn, err := zsys.TryBind(ctx, hcAddr, zsys.DefaultBindBackoff())
            //if err != nil {
            //    svcLn.Close()
            //    logger.Errorf("failed to bind to hc addr %s: %v", hcAddr, err)
            //    return err
            //}
            //
            //cfg := csa.DefaultCSAConfig(swUdsPath)
            //if err := csa.CreateCSA(ctx, cfg, svcLn, hcLn); err != nil {
            //    return err
            //}

            return nil
        },
    }

    rootCmd.AddCommand(runCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Printf("Failed to start the load balancer. Error is: %v\n", err)
        os.Exit(-1)
    }
}
