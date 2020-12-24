package device

import (
	"context"
	"github.com/SwallowJ/loggo"
	"sync"
)

var logger = loggo.New("device")

//StartMonitor 开始监控
func StartMonitor(cx context.Context) {
	logger.Info("开始采集")
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(cx)
	defer cancel()

	wg.Add(1)
	go saveCPUInfo(ctx, wg)

	wg.Wait()

}
