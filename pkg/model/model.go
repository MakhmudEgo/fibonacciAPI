package model

type Request struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func NewRequest() *Request {
	return &Request{}
}
