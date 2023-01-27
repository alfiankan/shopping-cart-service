package transport

type BaseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
