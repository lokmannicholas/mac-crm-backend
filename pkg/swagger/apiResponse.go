package swagger

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type APIErrorResponse struct {
	Detail  string `json:"detail"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
