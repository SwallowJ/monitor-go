package device

import (
	"context"
	"github.com/SwallowJ/loggo"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

var logger = loggo.New("main")

//CPUInfo 信息
type CPUInfo struct {
	cpu.TimesStat
	cores   int
	logical int
	percent float64
}

//GetCputInfo 获取cpu信息
func GetCputInfo(ctx context.Context) {
	cors, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		logger.Error(err)
	}

	logical, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		logger.Error(err)
	}

	percent, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		logger.Error(err)
	}

	times, err := cpu.TimesWithContext(ctx, false)
	if err != nil {
		logger.Error(err)
	}

	logger.Info(times[0].String(), times[0].Total())

	logger.Info("核心数：", cors)
	logger.Info("线程数：", logical)
	logger.Info("使用率：", percent[0])
}
