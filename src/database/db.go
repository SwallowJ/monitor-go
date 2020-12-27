package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	// "fmt"
	"monitor/src/conf"
	"monitor/src/dialect"

	_ "github.com/ClickHouse/clickhouse-go" //clickhouse 连接驱动
	"github.com/SwallowJ/loggo"
)

var (
	logger = loggo.New("database")
	//ClickHouseEngine clickhouse连接
	ClickHouseEngine *sql.DB
	isFaild          bool
	ck               = &dialect.ClickhouseType{}
)

//InitClickhouse 初始化clickhouse
//创建相关数据库, 表
func InitClickhouse(ctx context.Context) error {
	config := conf.Config.Clickhouse

	dbConnStr := fmt.Sprintf(`tcp://%s:%d?user=%s&password=%s`,
		config.Host, config.Port, config.User, config.Password,
	)

	connect, err := sql.Open("clickhouse", dbConnStr)
	if err != nil {
		return err
	}
	defer connect.Close()

	if err := connect.Ping(); err != nil {
		return err
	}

	dbSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.Database)
	logger.Debug(dbSQL)
	if _, err := connect.ExecContext(ctx, dbSQL); err != nil {
		return err
	}

	connStr := fmt.Sprintf("%s&database=%s", dbConnStr, config.Database)
	ClickHouseEngine, err = sql.Open("clickhouse", connStr)
	if err != nil {
		return err
	}

	if err := connect.Ping(); err != nil {
		return err
	}

	tx, _ := ClickHouseEngine.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	createTable(ctx, tx, new(CPUInfo))
	createTable(ctx, tx, new(DiskInfo))
	createTable(ctx, tx, new(NetInfo))
	createTable(ctx, tx, new(MemoryInfo))

	if isFaild {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return errors.New("创建表失败")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func createTable(ctx context.Context, tx *sql.Tx, dest interface{}) {
	defer func() {
		if err := recover(); err != nil {
			isFaild = true
			logger.Error(err)
		}
	}()

	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	dbname := modelType.Name()

	var data uint8

	existSQL := fmt.Sprintf("EXISTS  %s", dbname)
	logger.Debug(existSQL)
	if err := tx.QueryRowContext(ctx, existSQL).Scan(&data); err != nil {
		logger.Fatal(err)
	}

	if data == 0 {
		var sql, order strings.Builder
		sql.WriteString(fmt.Sprintf("\nCREATE TABLE %s (\n", dbname))

		for i := 0; i < modelType.NumField(); i++ {
			p := modelType.Field(i)

			if !p.Anonymous {
				name := p.Tag.Get("json")
				sql.WriteByte('\t')
				sql.WriteString(name)
				order.WriteString(fmt.Sprintf("`%s`", name))
				sql.WriteString("\t\t")
				sql.WriteString(ck.DataTypeOf(p.Type))
				if i != modelType.NumField()-1 {
					order.WriteByte(',')
					sql.WriteString(",\n")
				}
			}
		}

		sql.WriteString(fmt.Sprintf("\n)engine=MergeTree()\nORDER BY (%s)", order.String()))
		logger.Debug(sql.String())

		stmt, _ := tx.Prepare(sql.String())
		defer stmt.Close()

		if _, err := stmt.Exec(); err != nil {
			logger.Error(err)
			isFaild = true
		}
	}
}

//CreateTx 创建事务
func CreateTx(ctx context.Context, query string) (*sql.Tx, *sql.Stmt, error) {
	logger.Debug(query)
	tx, err := ClickHouseEngine.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, nil, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, nil, err
	}

	return tx, stmt, nil
}
