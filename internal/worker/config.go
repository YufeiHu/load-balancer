package worker

import "time"

const (
	Version            = "0.1-snapshot"
	basePath           = "/v1/worker"
	serverShutdownTime = 500 * time.Millisecond
)
