package models

type Response struct {
	Code    int         `json:"code,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Title   string      `json:"title,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Response1 struct {
	Code    int         `json:"code,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Title   string      `json:"title,omitempty"`
	Message string      `json:"message,omitempty"`
	Hours   float64     `json:"hours,omitempty"`
	Sum     uint        `json:"sum,omitempty"`
}
