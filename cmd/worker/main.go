package main

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "load-balancer/internal/helper/context_helper"
    "load-balancer/internal/helper/log_helper"
    "load-balancer/internal/worker"
    "os"
    "strconv"
    "time"
)

func parseInputs(listenAddr, workerIdInput, timeoutSecondInput string) (workerId int, timeout time.Duration, err error) {
    if listenAddr == "" {
        err = fmt.Errorf("invalid listenAddr: %s", listenAddr)
        return
    }

    workerId, err = strconv.Atoi(workerIdInput)
    if err != nil {
        err = fmt.Errorf("invalid workerId: %s", workerIdInput)
        return
    }

    timeoutSecond, err := strconv.Atoi(timeoutSecondInput)
    if err != nil {
        err = fmt.Errorf("invalid timeout: %s", timeoutSecondInput)
        return
    }

    timeout = time.Duration(timeoutSecond) * time.Second
    return
}

func main() {
    rootCmd := cobra.Command{
        Use:           "worker",
        SilenceUsage:  true,
        SilenceErrors: true,
    }

    var (
        listenAddr       string
        workerIdStr      string
        timeoutSecondStr string
    )

    runCmd := &cobra.Command{
        Use:     "run",
        Short:   "Worker",
        Long:    "Worker",
        Version: worker.Version,
        RunE: func(cmd *cobra.Command, args []string) (err error) {
            ctx, cancel := context_helper.GetInterruptableCtx(context.Background())
            ctx = context_helper.SetWorkRequestId(ctx, fmt.Sprintf("WorkerID=%s", workerIdStr))
            defer cancel()

            workerId, timeout, err := parseInputs(listenAddr, workerIdStr, timeoutSecondStr)
            if err != nil {
                log_helper.LogError(ctx, "Failed to parse input parameters. Error: %v", err)
                return
            }

            log_helper.LogInfo(ctx, "Start running worker on %s with timeout=%v", listenAddr, timeout)
            err = worker.Run(ctx, workerId, listenAddr, timeout)
            if err != nil {
                log_helper.LogError(ctx, "Failed to run worker. Error: %v", err)
                return
            }

            return
        },
    }

    rootCmd.PersistentFlags().StringVar(&listenAddr, "listenAddr", "", "worker listening address")
    rootCmd.PersistentFlags().StringVar(&workerIdStr, "workerId", "", "worker ID")
    rootCmd.PersistentFlags().StringVar(&timeoutSecondStr, "timeoutSecond", "", "worker timeout")

    rootCmd.AddCommand(runCmd)

    if err := rootCmd.Execute(); err != nil {
        os.Exit(-1)
    }
}
