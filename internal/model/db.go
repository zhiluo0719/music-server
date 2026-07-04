package model

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var dialector gorm.Dialector

	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") || os.Getenv("DATABASE_URL") != "" {
		pgDSN := dsn
		if pgDSN == "" || pgDSN == "./data/music.db" {
			pgDSN = os.Getenv("DATABASE_URL")
		}
		if pgDSN == "" {
			pgHost := getEnv("PGHOST", "localhost")
			pgPort := getEnv("PGPORT", "5432")
			pgUser := getEnv("PGUSER", "postgres")
			pgPass := getEnv("PGPASSWORD", "")
			pgDB := getEnv("PGDATABASE", "railway")
			pgDSN = "host=" + pgHost + " user=" + pgUser + " password=" + pgPass + " dbname=" + pgDB + " port=" + pgPort + " sslmode=require TimeZone=Asia/Shanghai"
		}
		dialector = postgres.Open(pgDSN)
		log.Println("使用 PostgreSQL 数据库")
	} else {
		dir := filepath.Dir(dsn)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("创建数据库目录失败: %v", err)
		}
		dialector = sqlite.Open(dsn)
		log.Printf("使用 SQLite 数据库: %s", dsn)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	if err := DB.AutoMigrate(&Song{}, &Playlist{}, &PlaylistSong{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库初始化完成")
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
