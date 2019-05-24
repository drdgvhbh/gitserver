package internal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
)

type ResponsePayload struct {
	Data  map[string]interface{} `json:"data"`
	Error map[string]interface{} `json:"error"`
}

func CommitsHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	repositoryPath := vars["directory"]
	logrus.Println(vars["directory"])
	_, err := git.PlainOpen(repositoryPath)
	if err != nil {

	}
	logrus.Println(err.Error())

	dataPayload := ResponsePayload{
		Data: map[string]interface{}{
			"yolo": "swag",
		},
	}
	err = json.NewEncoder(writer).Encode(&dataPayload)
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
