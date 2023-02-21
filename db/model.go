package db

import (
	"time"
)

type User struct {
	ID            int64  `json:"id"`
	Username      string `json:"name"`
	Password      string
	FollowCount   int64 `json:"follow_count"`
	FollowerCount int64 `json:"follower_count"`
	IsFollow      bool  `json:"is_follow" gorm:"-:all"`
}

type Video struct {
	ID            int64 `json:"id"`
	AuthorId      int64
	AuthorName    string
	CreateTime    time.Time
	Author        User   `json:"author" gorm:"foreignKey:AuthorId"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite" gorm:"-:all"`
	Title         string `json:"title"`
}

type Comment struct {
	ID         int64 `json:"id"`
	UserId     int64
	UserName   string
	VideoId    int64
	Content    string
	CreateTime string
}

type Like struct {
	ID      int64
	UserId  int64
	VideoId int64
}
