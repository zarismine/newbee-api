package web

type JsonResult struct {
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Success   bool        `json:"success"`
}

func JsonData(data interface{}) *JsonResult {
	return &JsonResult{
		ErrorCode: 0,
		Data:      data,
		Success:   true,
	}
}

func JsonSuccess() *JsonResult {
	return &JsonResult{
		ErrorCode: 0,
		Data:      nil,
		Success:   true,
	}
}

func JsonError(err error) *JsonResult {
	if err == nil {
		return JsonSuccess()
	}
	if e, ok := err.(*CodeError); ok {
		return &JsonResult{
			ErrorCode: e.Code,
			Message:   e.Message,
			Data:      e.Data,
			Success:   false,
		}
	}
	return &JsonResult{
		ErrorCode: 0,
		Message:   err.Error(),
		Data:      nil,
		Success:   false,
	}
}