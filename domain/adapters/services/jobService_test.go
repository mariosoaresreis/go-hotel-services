package services

import (
	"github.com/golang/mock/gomock"
	"github.com/mariosoaresreis/go-hotel/domain"
	services "github.com/mariosoaresreis/go-hotel/domain/adapters/mocks"
	"github.com/mariosoaresreis/go-hotel/domain/dto"
	"github.com/mariosoaresreis/go-hotel/domain/ports/interfaces"
	"net/http"
	"sync"
	"testing"
)

func GetLocationsForTest() []domain.Location {
	parent1 := struct {
		Id          int64
		Name        string
		DisplayName string
	}{
		Id:          0,
		Name:        "",
		DisplayName: "",
	}

	parent2 := struct {
		Id          int64
		Name        string
		DisplayName string
	}{
		Id:          2,
		Name:        "Floor 1",
		DisplayName: "Floor 1",
	}

	var location1 = new(domain.Location)
	location1.ParentLocation = struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}(parent1)
	location1.LocationType = domain.LocationType{Id: locationTypeFloor, DisplayName: "Floor"}
	location1.Id = 2
	location1.DisplayName = "Floor 1"
	location1.Name = "Floor 1"

	var room1 = new(domain.Location)
	room1.ParentLocation = struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}(parent2)
	room1.LocationType = domain.LocationType{Id: locationTypeRoom, DisplayName: "Room"}
	room1.Id = 3
	room1.DisplayName = "Room 1"
	room1.Name = "Room 1"

	var room2 = new(domain.Location)
	room2.ParentLocation = struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}(parent2)
	room2.LocationType = domain.LocationType{Id: locationTypeRoom, DisplayName: "Room"}
	room2.Id = 3
	room2.DisplayName = "Room 2"
	room2.Name = "Room 2"

	locations := make([]domain.Location, 0)
	locations = append(locations, *location1)
	locations = append(locations, *room1)
	locations = append(locations, *room2)
	return locations
}
func Test_addJob_should_return_code_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := services.NewMockSupplierService(ctrl)

	var errorMessageMap = interfaces.MessageMap{
		M:     make(map[int64]string),
		Mutex: sync.RWMutex{},
	}

	var locationMap = interfaces.LocationMap{
		M:     make(map[int64]domain.Location),
		Mutex: sync.RWMutex{},
	}

	mockService.EXPECT().GetJobItem(int64(1), gomock.Any(), gomock.Any()).Return(&domain.JobItemResponse{Id: 0, DisplayName: ""}, nil)
	mockService.EXPECT().GetLocation(int64(2), gomock.Any(), gomock.Any(), errorMessageMap, locationMap).Return(&domain.Location{Id: 0}, nil)
	mockService.EXPECT().GetDepartment(int64(1), gomock.Any(), gomock.Any()).Return(&domain.Department{ID: 0, Name: ""}, nil)

	jobService := DefaultJobService{mockService}
	jobRequest := dto.JobRequestDTO{
		JobItem:    1,
		Department: 1,
		Locations:  []int64{2},
	}

	var result, err = jobService.AddJob(jobRequest)

	if err.Code != http.StatusBadRequest {
		t.Error("Test_addJob_should_return_code_400 failed - Code is not 400")
	}

	if result != nil {
		t.Error("Test_addJob_should_return_code_400 failed - Result is not nil")
	}

}

func Test_clean_Bed_Created_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := services.NewMockSupplierService(ctrl)

	var errorMessageMap = interfaces.MessageMap{
		M:     make(map[int64]string),
		Mutex: sync.RWMutex{},
	}

	var locationMap = interfaces.LocationMap{
		M:     make(map[int64]domain.Location),
		Mutex: sync.RWMutex{},
	}

	locations := GetLocationsForTest()

	mockService.EXPECT().GetJobItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.JobItemResponse{Id: mattress,
		DisplayName: "item"}, nil)
	mockService.EXPECT().GetLocation(int64(1), gomock.Any(), gomock.Any(), errorMessageMap, locationMap).Return(&domain.Location{Id: 1,
		DisplayName: "A"}, nil)
	mockService.EXPECT().GetDepartment(departmentHouseKeeping, gomock.Any(), gomock.Any()).Return(&domain.Department{ID: departmentHouseKeeping,
		Name: "dep"}, nil)
	mockService.EXPECT().CreateJob(gomock.Any(), gomock.Any()).Return(&domain.JobResponse{Id: 1}, nil)
	mockService.EXPECT().GetAllLocations().Return(locations, nil)
	jobService := DefaultJobService{mockService}
	jobRequest := dto.JobRequestDTO{
		JobItem:    mattress,
		Department: departmentHouseKeeping,
		Locations:  []int64{1},
	}

	var result, err = jobService.AddJob(jobRequest)

	if err != nil {
		t.Error("Test_clean_Bed_Created_200 - Code is not 200")
	}

	if result == nil || result.Id == 0 {
		t.Error("Test_clean_Bed_Created_200 - Id is 0")
	}
}

// 5. If the given department is ‘Engineering’ and a job_item and at least one location is given then create a job to repair the given job item at the given location(s).
// * If only one location is given and is of location type ‘Floor’ then create a job to repair the given job item in all locations on that floor.
func Test_repair_Created_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := services.NewMockSupplierService(ctrl)

	var errorMessageMap = interfaces.MessageMap{
		M:     make(map[int64]string),
		Mutex: sync.RWMutex{},
	}

	var locationMap = interfaces.LocationMap{
		M:     make(map[int64]domain.Location),
		Mutex: sync.RWMutex{},
	}

	locations := GetLocationsForTest()

	mockService.EXPECT().GetJobItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.JobItemResponse{Id: 81, DisplayName: "item"}, nil)
	mockService.EXPECT().GetLocation(int64(1), gomock.Any(), gomock.Any(), errorMessageMap, locationMap).Return(&domain.Location{Id: 1, DisplayName: "A"}, nil)
	mockService.EXPECT().GetDepartment(gomock.Any(), gomock.Any(),
		gomock.Any()).Return(&domain.Department{ID: departmentEngineering, Name: "dep"}, nil)
	mockService.EXPECT().CreateJob(gomock.Any(), gomock.Any()).Return(&domain.JobResponse{Id: 1}, nil)
	mockService.EXPECT().GetAllLocations().Return(locations, nil)

	jobService := DefaultJobService{mockService}
	jobRequest := dto.JobRequestDTO{
		JobItem:    81,
		Department: departmentEngineering,
		Locations:  []int64{1},
	}

	var result, err = jobService.AddJob(jobRequest)

	if err != nil {
		t.Error("Test_clean_Bed_Created_200 - error")
	}

	if result.Id == 0 {
		t.Error("Test_clean_Bed_Created_200 - Id is 0")
	}
}

// 6. If the department is ‘Room Service’ and a job item and at least 1 location is given then create a job to deliver that job item to the given locations.
// * If only one location is given and is of location type ‘Floor’ then create a job to deliver the given job item in all locations with a location type of 'Room' on that floor.
func Test_deliver_job_created_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := services.NewMockSupplierService(ctrl)

	var errorMessageMap = interfaces.MessageMap{
		M:     make(map[int64]string),
		Mutex: sync.RWMutex{},
	}

	var locationMap = interfaces.LocationMap{
		M:     make(map[int64]domain.Location),
		Mutex: sync.RWMutex{},
	}

	locations := GetLocationsForTest()

	mockService.EXPECT().GetJobItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.JobItemResponse{Id: 81, DisplayName: "item"}, nil)
	mockService.EXPECT().GetLocation(int64(1), gomock.Any(), gomock.Any(), errorMessageMap, locationMap).Return(&domain.Location{Id: 1, DisplayName: "A"}, nil)
	mockService.EXPECT().GetDepartment(gomock.Any(), gomock.Any(),
		gomock.Any()).Return(&domain.Department{ID: departmentRoomService, Name: "Room Service"}, nil)
	mockService.EXPECT().CreateJob(gomock.Any(), gomock.Any()).Return(&domain.JobResponse{Id: 1}, nil)
	mockService.EXPECT().GetAllLocations().Return(locations, nil)

	jobService := DefaultJobService{mockService}
	jobRequest := dto.JobRequestDTO{
		JobItem:    81,
		Department: departmentRoomService,
		Locations:  []int64{1},
	}

	var result, err = jobService.AddJob(jobRequest)

	if err != nil {
		t.Error("Test_deliver_job_created_200 - error")
	}

	if result.Id == 0 {
		t.Error("Test_deliver_job_created_200 - Id is 0")
	}
}

// 7. If none of the above criteria is met return a bad request.
func Test_no_criteria_matches_code_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := services.NewMockSupplierService(ctrl)

	var errorMessageMap = interfaces.MessageMap{
		M:     make(map[int64]string),
		Mutex: sync.RWMutex{},
	}

	var locationMap = interfaces.LocationMap{
		M:     make(map[int64]domain.Location),
		Mutex: sync.RWMutex{},
	}
	mockService.EXPECT().GetJobItem(int64(1), gomock.Any(), gomock.Any()).Return(&domain.JobItemResponse{Id: 1, DisplayName: ""}, nil)
	mockService.EXPECT().GetLocation(int64(1), gomock.Any(), gomock.Any(), errorMessageMap, locationMap).Return(&domain.Location{Id: 1, DisplayName: "A"}, nil)
	mockService.EXPECT().GetDepartment(int64(1), gomock.Any(), gomock.Any()).Return(&domain.Department{ID: 1, Name: "dep"}, nil)

	jobService := DefaultJobService{mockService}
	jobRequest := dto.JobRequestDTO{
		JobItem:    1,
		Department: 1,
		Locations:  []int64{1},
	}

	var result, err = jobService.AddJob(jobRequest)

	if err.Code != http.StatusBadRequest {
		t.Error("Test_no_criteria_matches_code_400 - Code is not 400")
	}

	if result != nil {
		t.Error("Test_no_criteria_matches_code_400 - Result is not nil")
	}

	if err.Message != "Request hasn't matched any expected criteria. Please, review your request" {
		t.Error("Test_no_criteria_matches_code_400 - didn't reach no criteria matches")
	}

}
