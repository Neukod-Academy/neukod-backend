package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//HTTP

type HttpResponseBody struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func (res HttpResponseBody) UpdateHttpResponse(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	body := HttpResponseBody{
		Status:  res.Status,
		Message: res.Message,
		Data:    res.Data,
	}
	byteBody, err := json.Marshal(body)
	if err != nil {
		http.Error(writer, "Failed to encode to json file", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(body.Status)
	writer.Write(byteBody)
}

func HttpReqReader[T interface{}](r *http.Request) (T, error) {
	var data T
	if ByteBody, err := io.ReadAll(r.Body); err != nil {
		return data, fmt.Errorf("an error ocured: %s", err.Error())
	} else {
		if err = json.Unmarshal(ByteBody, &data); err != nil {
			return data, fmt.Errorf("an error ocured: %s", err.Error())
		}
		return data, nil
	}
}
