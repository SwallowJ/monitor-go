package device

import (
	"context"
	"github.com/SwallowJ/loggo"
	"github.com/shirou/gopsutil/v3/cpu"
	// "github.com/shirou/gopsutil/process"
	"time"
)

var logger = loggo.New("main")

//CPUInfo 信息
type CPUInfo struct {
	Cors      int     `json:"cors"`    //内核
	Logical   int     `json:"logical"` //逻辑处理器
	Percent   float64 `json:"percent"` //利用率
	CPU       string  `json:"cpu"`     //cpu
	User      float64 `json:"user"`    //用户态空间运行时间
	System    float64 `json:"system"`  //内核空间运行时间
	Idle      float64 `json:"idle"`    //空闲时间
	Nice      float64 `json:"nice"`    //用户空间进程的CPU的调度优先级
	Iowait    float64 `json:"iowait"`  //读写等待状态时间
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}

//GetCputInfo 获取cpu信息
func GetCputInfo(ctx context.Context) {

	Info := &CPUInfo{}

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

	logger.Info(times)

	Info.Cors = cors
	Info.Logical = logical
	logger.Info("核心数：", cors)
	logger.Info("线程数：", logical)
	logger.Info("使用率：", percent[0])

	// proce, _ := process.NewProcess()
}
