package monitor

import (
	"context"
	"monitor/src/conf"
	"monitor/src/database"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

func saveMemoryInfo(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	query := `INSERT INTO MemoryInfo 
		(total,available,used,usedPercent,free,active,inactive,wired,buffers,cached,
			activefile,inactivefile,activeanon,inactiveanon,unevictable,createTime)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

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
			memory, _ := mem.VirtualMemoryWithContext(ctx)
			stat, _ := mem.VirtualMemoryExWithContext(ctx)

			now := time.Now().Format("2006-01-02 15:04:05")
			if _, err := stmt.Exec(
				memory.Total,
				memory.Available,
				memory.Used,
				memory.UsedPercent,
				memory.Free,
				memory.Active,
				memory.Inactive,
				memory.Wired,
				memory.Buffers,
				memory.Cached,
				stat.ActiveFile,
				stat.InactiveFile,
				stat.ActiveAnon,
				stat.InactiveAnon,
				stat.Unevictable,
				now,
			); err != nil {
				logger.Error(err)
				return
			}

			time.Sleep(time.Duration(conf.Config.Interval) * time.Second)

			count++
			if count == conf.Config.Clickhouse.Nums {
				if err := tx.Commit(); err != nil {
					logger.Error(err)
				}
				stmt.Close()
				logger.Printf("保存Memory[%d]条数据\n", count)

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
