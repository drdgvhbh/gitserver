package internal

import (
	"encoding/json"
	"net/http"
)

type ResponsePayload struct {
	Data map[string]interface{} `json:"data"`
	Error map[string]interface{} `json:"error"`
}

func CommitsHandler(writer http.ResponseWriter, request *http.Request) {
	dataPayload := ResponsePayload{
		Data: map[string]interface{}{
			"yolo": "swag",
		},
	}
	err := json.NewEncoder(writer).Encode(&dataPayload)
	if err != nil {
		errorPayload := ResponsePayload{
			Error: map[string]interface{}{
				"error": err,
			},
		}
		err = json.NewEncoder(writer).Encode(&errorPayload)
		if err != nil {
			panic(err)
		}
	}
}