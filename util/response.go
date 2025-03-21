package util

import (
	"encoding/json"
)

type Response struct {
	Status int
	Error  error
	Data   interface{}
}

func GetBytesResponse(status int, data interface{}) ([]byte, error) {
	resp := Response{
		Status: status,
		Data:   data,
		Error:  nil,
	}

	return json.Marshal(&resp)
}
