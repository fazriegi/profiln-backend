package model

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
