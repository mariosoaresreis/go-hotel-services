package app

import (
	"encoding/json"
	"fmt"
	"github.com/mariosoaresreis/go-hotel/domain/dto"
	"github.com/mariosoaresreis/go-hotel/domain/ports/interfaces"
	"net/http"
)

type JobHandlers struct {
	service interfaces.JobService
}

func hello(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello World!")
}
func (ch *JobHandlers) addJob(writer http.ResponseWriter, request *http.Request) {
	var jobRequest dto.JobRequestDTO
	err := json.NewDecoder(request.Body).Decode(&jobRequest)

	if err != nil {
		writeResponse(writer, http.StatusBadRequest, err.Error())
	} else {
		job, appError := ch.service.AddJob(jobRequest)

		if appError != nil {
			writeResponse(writer, appError.Code, appError.Message)
		} else {
			writeResponse(writer, http.StatusOK, job)
		}
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if error := json.NewEncoder(w).Encode(data); error != nil {
		panic(error)
	}
}
