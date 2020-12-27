package monitor

import (
	"context"
	"monitor/src/conf"
	"monitor/src/database"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

func saveDiskInfo(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	query := `INSERT INTO DiskInfo 
			(path,fstype,total,free,used,speed,usedPercent,createTime)
			VALUES (?,?,?,?,?,?,?,?)`
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
			begin, err := disk.UsageWithContext(ctx, "/")
			if err != nil {
				logger.Error(err)
				return
			}
			time.Sleep(time.Duration(conf.Config.Interval) * time.Second)
			end, err := disk.UsageWithContext(ctx, "/")
			if err != nil {
				logger.Error(err)
				return
			}

			var speed float64
			if begin.Used < end.Used {
				speed = float64(end.Used-begin.Used) / float64(conf.Config.Interval)
			}

			now := time.Now().Format("2006-01-02 15:04:05")
			if _, err := stmt.Exec(
				end.Path,
				end.Fstype,
				end.Total,
				end.Free,
				end.Used,
				speed,
				end.UsedPercent,
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
				logger.Printf("保存DISK[%d]条数据\n", count)

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
