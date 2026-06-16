package db

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"openerp/packages/pkg-logger"
)

var DB *gorm.DB

func InitDB() {
	log := logger.Ctx(context.Background())
	log.Info("正在连接全局 PostgreSQL 18 数据库引擎...")

	var err error
	dsn := "host=localhost user=postgres password=123456@a dbname=openerp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("PostgreSQL 数据库连接失败: " + err.Error())
		panic("Failed to connect to PostgreSQL database")
	}
	
	log.Info("🚀 PostgreSQL 18 (JSONB Ready) 连接成功!")
}
