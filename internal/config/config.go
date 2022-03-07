package config

import "time"

const (
    Version = "0.1-snapshot"

    LbListenAddr = "localhost:7777"
    LbHealthAddr = "localhost:7778"

    LbRequestTimeout = 10 * time.Second

    LbBasePath = "/v1/demo"
)
