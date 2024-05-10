package model

import "time"

type ListPostsResponse struct {
	ID           int       `json:"id"`
	User         User      `json:"user"`
	Content      string    `json:"content"`
	ImageUrl     string    `json:"image_url"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	RepostCount  int       `json:"repost_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}
