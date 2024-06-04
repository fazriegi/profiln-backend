package model

type Response struct {
	Status Status `json:"status"`
	Data   any    `json:"data"`
}

type Status struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	IsSuccess bool   `json:"is_success"`
}

type PaginationRequest struct {
	Page    int
	Limit   int
	OrderBy string
}

type PaginationResponse struct {
	Page             int   `json:"page"`
	TotalPages       int   `json:"total_pages"`
	TotalRows        int64 `json:"total_rows"`
	CurrentRowsCount int   `json:"current_rows_count"`
}

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	AvatarUrl  string `json:"avatar_url"`
	Fullname   string `json:"fullname"`
	Bio        string `json:"bio"`
	OpenToWork bool   `json:"open_to_work"`
}

type School struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type Skill struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type Company struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type IssuingOrganization struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type JobPosition struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
