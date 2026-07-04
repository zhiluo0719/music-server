package model

import "time"

type Playlist struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	CoverPath   string    `json:"cover_path" gorm:"type:varchar(500);default:''"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Songs       []Song    `json:"songs,omitempty" gorm:"-"`
}

func (Playlist) TableName() string {
	return "playlists"
}

type PlaylistSong struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PlaylistID uint      `json:"playlist_id" gorm:"not null;index"`
	SongID     uint      `json:"song_id" gorm:"not null;index"`
	SortOrder  int       `json:"sort_order" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PlaylistSong) TableName() string {
	return "playlist_songs"
}

type AddSongRequest struct {
	SongID uint `json:"song_id" binding:"required"`
}
