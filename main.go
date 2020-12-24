package main

import (
	"context"
	"github.com/SwallowJ/loggo"
	"sync"

	"monitor/device"
)

var logger *loggo.Logger
var wg sync.WaitGroup

func init() {
	loggo.SetServiceName("monitor")
	logger = loggo.New("main")
}

func main() {
	logger.Info("Start")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	device.StartMonitor(ctx)
}
