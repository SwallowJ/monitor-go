package monitor

import (
	"context"
	"monitor/src/conf"
	"monitor/src/database"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func saveCPUInfo(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	query := `INSERT INTO CPUInfo 
		(cors,logical,percent,cpu,user,system,idle,nice,iowait,irq,softirq,
			steal,guest,guestNice,createTime)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	var count int
	tx, stmt, err := database.CreateTx(ctx, query)
	if err != nil {
		logger.Error(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			if err := tx.Commit(); err != nil {
				logger.Error(err)
			}
			stmt.Close()
			return
		default:
			cors, _ := cpu.CountsWithContext(ctx, false)
			logical, _ := cpu.CountsWithContext(ctx, true)
			percent, _ := cpu.PercentWithContext(ctx, time.Second, false)

			timeStats, _ := cpu.TimesWithContext(ctx, false)
			timeStat := timeStats[0]

			now := time.Now().Format("2006-01-02 15:04:05")
			if _, err := stmt.Exec(
				cors,
				logical,
				percent[0],
				timeStat.CPU,
				timeStat.User,
				timeStat.System,
				timeStat.Idle,
				timeStat.Nice,
				timeStat.Iowait,
				timeStat.Irq,
				timeStat.Softirq,
				timeStat.Steal,
				timeStat.Guest,
				timeStat.GuestNice,
				now,
			); err != nil {
				logger.Error(err)
				return
			}

			count++
			if count == conf.Config.Clickhouse.Nums {
				if err := tx.Commit(); err != nil {
					logger.Error(err)
				}
				stmt.Close()
				logger.Printf("保存CPU[%d]条数据\n", count)

				count = 0
				tx, stmt, err = database.CreateTx(ctx, query)
				if err != nil {
					logger.Error(err)
					return
				}
			}
		}
	}

}
