package model

import "time"

type ListPostsResponse struct {
	ID           int64     `json:"id"`
	User         User      `json:"user"`
	Content      string    `json:"content"`
	ImageUrls    []string  `json:"image_urls"`
	LikeCount    int32     `json:"like_count"`
	CommentCount int32     `json:"comment_count"`
	RepostCount  int32     `json:"repost_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}
