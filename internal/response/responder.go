package response

import "net/http"

type ApiResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}, message string) (resp ApiResponse, status int) {
	resp = ApiResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	return resp, http.StatusOK
}

func Error(message string, statusCode int, data interface{}) (resp ApiResponse, status int) {
	if data == nil {
		data = make(map[string]string, 0)
	}
	resp = ApiResponse{
		Status:  false,
		Message: message,
		Data:    data,
	}
	return resp, statusCode
}