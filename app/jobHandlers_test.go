package app

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mariosoaresreis/go-hotel/domain"
	services "github.com/mariosoaresreis/go-hotel/domain/adapters/mocks"
	"github.com/mariosoaresreis/go-hotel/domain/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_should_return_job_with_status_code_200(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := services.NewMockJobService(ctrl)

	job := domain.Job{
		Item: domain.JobItem{Id: 1},
	}

	jobRequest := dto.JobRequestDTO{
		JobItem:    1,
		Department: 1,
		Locations:  []int64{1, 2, 3},
	}

	var marshalled, err = json.Marshal(jobRequest)

	if err != nil {
		panic("impossible to marshall object: %s")
	}

	mockService.EXPECT().AddJob(jobRequest).Return(&job, nil)
	jh := JobHandlers{service: mockService}
	router := mux.NewRouter()
	router.HandleFunc("/addJob", jh.addJob)
	request, _ := http.NewRequest(http.MethodPost, "/addJob", bytes.NewReader(marshalled))
	//Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Test_should_return_job_with_status_code_200 failed")
	}
}
