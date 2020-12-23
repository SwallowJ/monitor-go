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

	// getCputInfo(ctx)
	// device.GetCputInfo(ctx)
	wg.Add(1)
	go test(ctx)

	wg.Wait()
}

func test(ctx context.Context) {
	defer wg.Done()

	// for i := 0; i < 100; i++ {
	// }
	device.GetCputInfo(ctx)
}
