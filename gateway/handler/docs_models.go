package handler

type _BaseResponse struct {
	Data   any    `json:"data"`                     // 数据
	Status string `json:"status" example:"success"` // 状态
}

type _ResponseError struct {
	Error  string `json:"error" example:"err msg"` // 错误信息
	Status string `json:"status" example:"failed"` // 状态
}
