package response

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

type MetaData struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func Success(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error(err error) Response {
	return Response{
		Success: false,
		Error:   err.Error(),
	}
}

func SuccessWithData(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}
