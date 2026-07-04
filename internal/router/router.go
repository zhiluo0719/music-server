package router

import (
	"music-backend/internal/handler"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Range")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Accept-Ranges, Content-Range")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	r.MaxMultipartMemory = 100 << 20

	r.Static("/uploads", "./uploads")

	songH := handler.NewSongHandler()
	playlistH := handler.NewPlaylistHandler()

	api := r.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.GET("", songH.List)
			songs.POST("", songH.Create)
			songs.GET("/:id", songH.GetByID)
			songs.PUT("/:id", songH.Update)
			songs.DELETE("/:id", songH.Delete)
			songs.GET("/:id/stream", songH.Stream)
		}

		api.GET("/artists", songH.GetArtists)
		api.GET("/albums", songH.GetAlbums)
		api.GET("/genres", songH.GetGenres)

		playlists := api.Group("/playlists")
		{
			playlists.GET("", playlistH.List)
			playlists.POST("", playlistH.Create)
			playlists.GET("/:id", playlistH.GetByID)
			playlists.PUT("/:id", playlistH.Update)
			playlists.DELETE("/:id", playlistH.Delete)
			playlists.POST("/:id/songs", playlistH.AddSong)
			playlists.DELETE("/:id/songs/:songId", playlistH.RemoveSong)
		}
	}

	publicDir := "./public"
	if _, err := os.Stat(publicDir); err == nil {
		assetsDir := filepath.Join(publicDir, "assets")
		if _, err := os.Stat(assetsDir); err == nil {
			r.Static("/assets", assetsDir)
		}
		r.StaticFile("/", filepath.Join(publicDir, "index.html"))
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(publicDir, "index.html"))
		})
	}

	return r
}
