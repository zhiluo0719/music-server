package service

import (
	"errors"
	"music-backend/internal/model"

	"gorm.io/gorm"
)

type PlaylistService struct{}

func (s *PlaylistService) List() ([]model.Playlist, error) {
	var playlists []model.Playlist
	if err := model.DB.Order("created_at DESC").Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}

func (s *PlaylistService) GetByID(id uint) (*model.Playlist, error) {
	var playlist model.Playlist
	if err := model.DB.First(&playlist, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var psList []model.PlaylistSong
	if err := model.DB.Where("playlist_id = ?", id).Order("sort_order ASC").Find(&psList).Error; err != nil {
		return nil, err
	}

	if len(psList) > 0 {
		songIDs := make([]uint, len(psList))
		for i, ps := range psList {
			songIDs[i] = ps.SongID
		}

		var songs []model.Song
		if err := model.DB.Where("id IN ?", songIDs).Find(&songs).Error; err != nil {
			return nil, err
		}

		songMap := make(map[uint]model.Song, len(songs))
		for _, s := range songs {
			songMap[s.ID] = s
		}

		playlist.Songs = make([]model.Song, 0, len(psList))
		for _, ps := range psList {
			if s, ok := songMap[ps.SongID]; ok {
				playlist.Songs = append(playlist.Songs, s)
			}
		}
	} else {
		playlist.Songs = []model.Song{}
	}

	return &playlist, nil
}

func (s *PlaylistService) Create(playlist *model.Playlist) error {
	return model.DB.Create(playlist).Error
}

func (s *PlaylistService) Update(playlist *model.Playlist) error {
	return model.DB.Save(playlist).Error
}

func (s *PlaylistService) Delete(id uint) error {
	tx := model.DB.Begin()
	if err := tx.Where("playlist_id = ?", id).Delete(&model.PlaylistSong{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&model.Playlist{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *PlaylistService) AddSong(playlistID, songID uint) error {
	var count int64
	if err := model.DB.Model(&model.PlaylistSong{}).
		Where("playlist_id = ? AND song_id = ?", playlistID, songID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("歌曲已在播放列表中")
	}

	var maxOrder int
	model.DB.Model(&model.PlaylistSong{}).
		Where("playlist_id = ?", playlistID).
		Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	ps := model.PlaylistSong{
		PlaylistID: playlistID,
		SongID:     songID,
		SortOrder:  maxOrder + 1,
	}
	return model.DB.Create(&ps).Error
}

func (s *PlaylistService) RemoveSong(playlistID, songID uint) error {
	return model.DB.Where("playlist_id = ? AND song_id = ?", playlistID, songID).
		Delete(&model.PlaylistSong{}).Error
}
