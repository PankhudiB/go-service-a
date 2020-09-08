package request

type HelloBRequest struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
