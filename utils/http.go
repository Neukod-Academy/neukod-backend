package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpResponseBody struct {
	Status  int
	Message string
	Data    interface{}
}

func (res HttpResponseBody) UpdateHttpResponse(writer http.ResponseWriter, newStatus int, newData interface{}) {
	body := HttpResponseBody{
		Status: newStatus,
		Data:   newData,
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
