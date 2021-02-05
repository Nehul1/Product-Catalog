package model

type Brand struct {
	Id    int    `json:"-"` //`json:"id omit empty"`
	Brand string `json:"name"`
}
