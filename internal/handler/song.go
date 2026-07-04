package handler

import (
	"fmt"
	"io"
	"music-backend/internal/model"
	"music-backend/internal/service"
	"music-backend/pkg/response"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	svc service.SongService
}

func NewSongHandler() *SongHandler {
	return &SongHandler{}
}

func (h *SongHandler) List(c *gin.Context) {
	var q model.SongListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	songs, total, err := h.svc.List(q)
	if err != nil {
		response.InternalError(c, "查询歌曲列表失败")
		return
	}
	response.SuccessPage(c, songs, total, q.Page, q.PageSize)
}

func (h *SongHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	song, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.InternalError(c, "查询歌曲失败")
		return
	}
	if song == nil {
		response.NotFound(c, "歌曲不存在")
		return
	}
	response.Success(c, song)
}

func (h *SongHandler) Create(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择音频文件")
		return
	}

	audioDir := "./uploads/audio"
	os.MkdirAll(audioDir, 0755)

	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(audioDir, fileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.InternalError(c, "文件保存失败")
		return
	}

	song := model.Song{
		Title:      c.PostForm("title"),
		Artist:     c.PostForm("artist"),
		Album:      c.PostForm("album"),
		Genre:      c.PostForm("genre"),
		FilePath:   savePath,
		FileSize:   file.Size,
		FileFormat: strings.TrimPrefix(ext, "."),
	}

	if song.Title == "" {
		song.Title = strings.TrimSuffix(file.Filename, ext)
	}

	coverFile, err := c.FormFile("cover")
	if err == nil {
		coverDir := "./uploads/covers"
		os.MkdirAll(coverDir, 0755)
		coverExt := strings.ToLower(filepath.Ext(coverFile.Filename))
		coverName := fmt.Sprintf("%d%s", time.Now().UnixNano(), coverExt)
		coverPath := filepath.Join(coverDir, coverName)
		if err := c.SaveUploadedFile(coverFile, coverPath); err == nil {
			song.CoverPath = coverPath
		}
	}

	if err := h.svc.Create(&song); err != nil {
		os.Remove(savePath)
		response.InternalError(c, "保存歌曲信息失败")
		return
	}

	response.Success(c, song)
}

func (h *SongHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	song, err := h.svc.GetByID(uint(id))
	if err != nil || song == nil {
		response.NotFound(c, "歌曲不存在")
		return
	}

	var input model.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	song.Title = input.Title
	song.Artist = input.Artist
	song.Album = input.Album
	song.Genre = input.Genre

	if err := h.svc.Update(song); err != nil {
		response.InternalError(c, "更新歌曲失败")
		return
	}
	response.Success(c, song)
}

func (h *SongHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	song, err := h.svc.GetByID(uint(id))
	if err != nil || song == nil {
		response.NotFound(c, "歌曲不存在")
		return
	}

	os.Remove(song.FilePath)
	if song.CoverPath != "" {
		os.Remove(song.CoverPath)
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除歌曲失败")
		return
	}
	response.Success(c, nil)
}

func (h *SongHandler) Stream(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	song, err := h.svc.GetByID(uint(id))
	if err != nil || song == nil {
		response.NotFound(c, "歌曲不存在")
		return
	}

	file, err := os.Open(song.FilePath)
	if err != nil {
		response.NotFound(c, "音频文件不存在")
		return
	}
	defer file.Close()

	stat, _ := file.Stat()
	c.Header("Content-Type", "audio/"+song.FileFormat)
	c.Header("Content-Length", strconv.FormatInt(stat.Size(), 10))
	c.Header("Accept-Ranges", "bytes")

	io.Copy(c.Writer, file)
}

func (h *SongHandler) GetArtists(c *gin.Context) {
	artists, err := h.svc.GetArtists()
	if err != nil {
		response.InternalError(c, "查询艺术家失败")
		return
	}
	response.Success(c, artists)
}

func (h *SongHandler) GetAlbums(c *gin.Context) {
	albums, err := h.svc.GetAlbums()
	if err != nil {
		response.InternalError(c, "查询专辑失败")
		return
	}
	response.Success(c, albums)
}

func (h *SongHandler) GetGenres(c *gin.Context) {
	genres, err := h.svc.GetGenres()
	if err != nil {
		response.InternalError(c, "查询流派失败")
		return
	}
	response.Success(c, genres)
}
