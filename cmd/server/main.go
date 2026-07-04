package main

import (
	"fmt"
	"io"
	"log"
	"music-backend/internal/model"
	"music-backend/internal/router"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	os.MkdirAll("./uploads/audio", 0755)
	os.MkdirAll("./uploads/covers", 0755)
	os.MkdirAll("./logs", 0755)

	logFile := filepath.Join("logs", fmt.Sprintf("server_%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		multi := io.MultiWriter(os.Stdout, f)
		log.SetOutput(multi)
		gin.DefaultWriter = multi
		gin.DefaultErrorWriter = multi
	} else {
		log.Printf("无法创建日志文件: %v, 使用标准输出", err)
	}

	dsn := "./data/music.db"
	if envDSN := os.Getenv("DATABASE_URL"); envDSN != "" {
		dsn = envDSN
		log.Printf("检测到 Railway PostgreSQL: %s...", dsn[:min(len(dsn), 50)])
	}

	absDSN, _ := filepath.Abs(dsn)
	log.Printf("数据库路径: %s", absDSN)
	log.Printf("日志文件: %s", logFile)

	model.InitDB(dsn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	r := router.SetupRouter()
	log.Printf("服务启动在 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
