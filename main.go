package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/SwallowJ/loggo"

	// "monitor/device"
	"monitor/src/database"
	"monitor/src/monitor"
	"sync"
)

var (
	logger *loggo.Logger
	wg     = &sync.WaitGroup{}
)

func init() {
	loggo.SetServiceName("monitor")
	logger = loggo.New("main")
}

func main() {
	logger.Info("Start")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-signalCh
		logger.Info("接收到退出信号")
		cancel()
	}()

	if err := database.InitClickhouse(ctx); err != nil {
		logger.Fatal(err)
	}

	monitor.StartMonitor(ctx, wg)

	wg.Wait()
}
