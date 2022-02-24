package config

import "time"

const (
    Version = "0.1-snapshot"

    ListenAddr = "localhost:7777"
    HealthAddr = "localhost:7778"

    RequestTimeout = 10 * time.Second
)
