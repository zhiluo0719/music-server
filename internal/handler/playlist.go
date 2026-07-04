package handler

import (
	"music-backend/internal/model"
	"music-backend/internal/service"
	"music-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	svc service.PlaylistService
}

func NewPlaylistHandler() *PlaylistHandler {
	return &PlaylistHandler{}
}

func (h *PlaylistHandler) List(c *gin.Context) {
	playlists, err := h.svc.List()
	if err != nil {
		response.InternalError(c, "查询播放列表失败")
		return
	}
	response.Success(c, playlists)
}

func (h *PlaylistHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	playlist, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.InternalError(c, "查询播放列表失败")
		return
	}
	if playlist == nil {
		response.NotFound(c, "播放列表不存在")
		return
	}
	response.Success(c, playlist)
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	var input model.Playlist
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.Create(&input); err != nil {
		response.InternalError(c, "创建播放列表失败")
		return
	}
	response.Success(c, input)
}

func (h *PlaylistHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	playlist, err := h.svc.GetByID(uint(id))
	if err != nil || playlist == nil {
		response.NotFound(c, "播放列表不存在")
		return
	}

	var input model.Playlist
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	playlist.Name = input.Name
	playlist.Description = input.Description

	if err := h.svc.Update(playlist); err != nil {
		response.InternalError(c, "更新播放列表失败")
		return
	}
	response.Success(c, playlist)
}

func (h *PlaylistHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	playlist, err := h.svc.GetByID(uint(id))
	if err != nil || playlist == nil {
		response.NotFound(c, "播放列表不存在")
		return
	}
	if err := h.svc.Delete(playlist.ID); err != nil {
		response.InternalError(c, "删除播放列表失败")
		return
	}
	response.Success(c, nil)
}

func (h *PlaylistHandler) AddSong(c *gin.Context) {
	playlistID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的播放列表ID")
		return
	}

	var req model.AddSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.AddSong(uint(playlistID), req.SongID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *PlaylistHandler) RemoveSong(c *gin.Context) {
	playlistID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的播放列表ID")
		return
	}
	songID, err := strconv.ParseUint(c.Param("songId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的歌曲ID")
		return
	}

	if err := h.svc.RemoveSong(uint(playlistID), uint(songID)); err != nil {
		response.InternalError(c, "移除歌曲失败")
		return
	}
	response.Success(c, nil)
}
