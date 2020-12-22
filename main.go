package main

import (
	"context"
	"github.com/SwallowJ/loggo"

	"monitor/device"
)

var logger = loggo.New("main")

func main() {
	logger.Info("Start")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// getCputInfo(ctx)
	device.GetCputInfo(ctx)

}
