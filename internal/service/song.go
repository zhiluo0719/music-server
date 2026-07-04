package service

import (
	"errors"
	"music-backend/internal/model"
	"strings"

	"gorm.io/gorm"
)

type SongService struct{}

func (s *SongService) List(q model.SongListQuery) ([]model.Song, int64, error) {
	db := model.DB.Model(&model.Song{})
	if q.Keyword != "" {
		kw := "%" + q.Keyword + "%"
		db = db.Where("title LIKE ? OR artist LIKE ? OR album LIKE ?", kw, kw, kw)
	}
	if q.Artist != "" {
		db = db.Where("artist = ?", q.Artist)
	}
	if q.Album != "" {
		db = db.Where("album = ?", q.Album)
	}
	if q.Genre != "" {
		db = db.Where("genre = ?", q.Genre)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}

	var songs []model.Song
	if err := db.Order("created_at DESC").Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&songs).Error; err != nil {
		return nil, 0, err
	}
	return songs, total, nil
}

func (s *SongService) GetByID(id uint) (*model.Song, error) {
	var song model.Song
	if err := model.DB.First(&song, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &song, nil
}

func (s *SongService) Create(song *model.Song) error {
	return model.DB.Create(song).Error
}

func (s *SongService) Update(song *model.Song) error {
	return model.DB.Save(song).Error
}

func (s *SongService) Delete(id uint) error {
	return model.DB.Delete(&model.Song{}, id).Error
}

func (s *SongService) GetArtists() ([]string, error) {
	var artists []string
	err := model.DB.Model(&model.Song{}).Distinct("artist").Where("artist != ''").Pluck("artist", &artists).Error
	return artists, err
}

func (s *SongService) GetAlbums() ([]string, error) {
	var albums []string
	err := model.DB.Model(&model.Song{}).Distinct("album").Where("album != ''").Pluck("album", &albums).Error
	return albums, err
}

func (s *SongService) GetGenres() ([]string, error) {
	var genres []string
	err := model.DB.Model(&model.Song{}).Distinct("genre").Where("genre != ''").Pluck("genre", &genres).Error
	return genres, err
}

func DetectFileFormat(filename string) string {
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	switch ext {
	case "mp3", "wav", "flac", "ogg", "aac", "wma", "m4a":
		return ext
	default:
		return ext
	}
}
