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
