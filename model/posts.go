package model

import "time"

type ReportPostRequest struct {
	PostId  int64  `json:"post_id"`
	Reason  string `json:"reason" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type ReportPostResponse struct {
	PostId  int64  `json:"post_id"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type Post struct {
	ID           int64     `json:"id"`
	User         User      `json:"user"`
	Content      string    `json:"content"`
	ImageUrl     string    `json:"image_url"`
	LikeCount    int32     `json:"like_count"`
	CommentCount int32     `json:"comment_count"`
	RepostCount  int32     `json:"repost_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostComment struct {
	ID           int64     `json:"id"`
	PostId       int64     `json:"post_id"`
	User         User      `json:"user"`
	Content      string    `json:"content"`
	ImageUrl     string    `json:"image_url"`
	LikeCount    int32     `json:"like_count"`
	ReplyCount   int32     `json:"reply_count"`
	IsPostAuthor bool      `json:"is_post_author"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostCommentReply struct {
	ID            int64     `json:"id"`
	PostCommentId int64     `json:"post_comment_id"`
	User          User      `json:"user"`
	Content       string    `json:"content"`
	ImageUrl      string    `json:"image_url"`
	LikeCount     int32     `json:"like_count"`
	IsPostAuthor  bool      `json:"is_post_author"`
	UpdatedAt     time.Time `json:"updated_at"`
}
