package monitor

import (
	"context"
	"monitor/src/conf"
	"monitor/src/database"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

func saveNetInfo(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	query := `INSERT INTO NetInfo 
		(name,bytesSent,bytesRecv,packetsSent,packetsRecv,errin,errout,dropin,
			dropout,fifoin,fifoout,sentSpeed,recvSpeed,createTime)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

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
			stat, _ := net.IOCountersWithContext(ctx, false)
			begin := stat[0]
			time.Sleep(time.Duration(conf.Config.Interval) * time.Second)
			stat, _ = net.IOCountersWithContext(ctx, false)
			end := stat[0]
			now := time.Now().Format("2006-01-02 15:04:05")

			var sentSpeed, recvSpeed float64
			if begin.BytesSent < end.BytesSent {
				sentSpeed = float64(end.BytesSent-begin.BytesSent) / float64(conf.Config.Interval)
			}

			if begin.BytesRecv < end.BytesRecv {
				recvSpeed = float64(end.BytesRecv-begin.BytesRecv) / float64(conf.Config.Interval)
			}

			if _, err := stmt.Exec(
				end.Name,
				end.BytesSent,
				end.BytesRecv,
				end.PacketsSent,
				end.PacketsRecv,
				end.Errin,
				end.Errout,
				end.Dropin,
				end.Dropout,
				end.Fifoin,
				end.Fifoout,
				sentSpeed,
				recvSpeed,
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

				logger.Printf("保存Net[%d]条数据\n", count)

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
