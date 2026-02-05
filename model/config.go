package model

import (
	"fmt"
	"os"
	"time"
	"zpic-client/helper"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 声明全局变量
// 全局数据库连接
var DB *gorm.DB

// 初始化数据库连接
func InitDB() {
	var err error
	//如果数据库文件不存在，会自动创建
	dsn := "data/db/zpic.db3"
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		SkipDefaultTransaction: true, // 减少写入开销（对 SQLite 通常更快）
		PrepareStmt:            true, // 重复 SQL 更快（会占用少量内存）
	})
	//如果出现错误，抛出错误并终止执行
	if err != nil {
		// 打印错误
		fmt.Println("Database connection failed!")
		// 写入日志
		helper.WriteLog(fmt.Sprintf("%s", err))
		// 终止执行
		os.Exit(1)
	} else {
		// 显式设置SQLite配置以确保WAL模式生效
		sqlDB, err := DB.DB()
		if err == nil {
			execPragma := func(q string) {
				if _, e := sqlDB.Exec(q); e != nil {
					helper.WriteLog(fmt.Sprintf("sqlite pragma failed: %s, err=%v", q, e))
				}
			}

			// 推荐：稳定性 + 性能（WAL 下常用组合）
			execPragma("PRAGMA journal_mode=WAL;")
			execPragma("PRAGMA synchronous=NORMAL;")
			execPragma("PRAGMA temp_store=memory;")
			execPragma("PRAGMA cache_size=-64000;")

			// 关键：并发/高频写入时避免频繁报 locked（单位 ms）
			execPragma("PRAGMA busy_timeout=5000;")

			// 关键：控制 WAL 文件增长（页数；可按写入量调整）
			execPragma("PRAGMA wal_autocheckpoint=1000;")

			// 可选：读多写少时常有收益（按需调整；0=关闭）
			execPragma("PRAGMA mmap_size=268435456;") // 256MB

			// 按需：开启外键约束（会有少量开销，但保证一致性）
			execPragma("PRAGMA foreign_keys=ON;")

			// 连接池：SQLite 写连接限制保持不变，但补齐生命周期更稳
			sqlDB.SetMaxOpenConns(1)
			sqlDB.SetMaxIdleConns(1)
			sqlDB.SetConnMaxLifetime(0)
			sqlDB.SetConnMaxIdleTime(5 * time.Minute)
		} else {
			helper.WriteLog(fmt.Sprintf("failed to get sql.DB: %v", err))
		}

		//自动迁移数据库表结构
		// 迁移多个表
		err = DB.AutoMigrate(&ZPurls{})
		if err != nil {
			fmt.Println("failed to migrate database!")
			// 写入错误日志
			helper.WriteLog(fmt.Sprintf("%s", err))
		}
		fmt.Print("Database connection succeeded!\n")
	}

}
