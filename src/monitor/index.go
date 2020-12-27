package monitor

import (
	"context"
	"monitor/src/conf"
	"sync"
	"time"

	"github.com/SwallowJ/loggo"
	"github.com/shirou/gopsutil/v3/net"
)

var logger = loggo.New("monitor")

//StartMonitor 开始监控
func StartMonitor(ctx context.Context, wg *sync.WaitGroup) {
	logger.Info("开始采集")

	go saveCPUInfo(ctx, wg)
	go saveDiskInfo(ctx, wg)
	go saveMemoryInfo(ctx, wg)
	go saveNetInfo(ctx, wg)

}

func test(ctx context.Context) {
	stat, _ := net.IOCountersWithContext(ctx, false)
	begin := stat[0]
	time.Sleep(time.Duration(conf.Config.Interval) * time.Second)
	stat, _ = net.IOCountersWithContext(ctx, false)
	end := stat[0]

	var sentSpeed float64

	if begin.BytesSent < end.BytesSent {
		sentSpeed = float64(end.BytesSent-begin.BytesSent) / float64(conf.Config.Interval)
	}

	logger.Info(sentSpeed)
}
