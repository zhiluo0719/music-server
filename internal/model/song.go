package model

import "time"

type Song struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"type:varchar(255);not null"`
	Artist     string    `json:"artist" gorm:"type:varchar(255);default:''"`
	Album      string    `json:"album" gorm:"type:varchar(255);default:''"`
	Genre      string    `json:"genre" gorm:"type:varchar(100);default:''"`
	Duration   int       `json:"duration" gorm:"default:0"`
	FilePath   string    `json:"file_path" gorm:"type:varchar(500);not null"`
	FileSize   int64     `json:"file_size" gorm:"default:0"`
	FileFormat string    `json:"file_format" gorm:"type:varchar(20);default:''"`
	CoverPath  string    `json:"cover_path" gorm:"type:varchar(500);default:''"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Song) TableName() string {
	return "songs"
}

type SongListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Artist   string `form:"artist"`
	Album    string `form:"album"`
	Genre    string `form:"genre"`
}
