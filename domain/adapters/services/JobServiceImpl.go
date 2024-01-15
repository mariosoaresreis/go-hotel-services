package services

import (
	"bytes"
	"fmt"
	"github.com/mariosoaresreis/go-hotel/domain"
	"github.com/mariosoaresreis/go-hotel/domain/dto"
	"github.com/mariosoaresreis/go-hotel/domain/ports/interfaces"
	"github.com/mariosoaresreis/go-hotel/errors"
	"sync"
)

type DefaultJobService struct {
	supplierService interfaces.SupplierService
}

const departmentHouseKeeping int64 = 7
const departmentEngineering int64 = 3

const departmentRoomService int64 = 13
const mattress = 81
const blanket = 3

var houseKeepingItems = map[int64]string{
	mattress: "Mattress",
	blanket:  "Blanket",
}

const locationTypeFloor = 2
const locationTypeRoom = 5

func (service DefaultJobService) AddJob(jobRequest dto.JobRequestDTO) (*domain.JobResponse, *errors.ApplicationError) {
	var department = new(domain.Department)
	var jobItem = new(domain.JobItemResponse)
	var location = new(domain.Location)
	var buffer bytes.Buffer
	var errorCount = 0
	var wg sync.WaitGroup

	var errorMessageMap = interfaces.MessageMap{
		M:     make(map[int64]string),
		Mutex: sync.RWMutex{},
	}

	var locationMap = interfaces.LocationMap{
		M:     make(map[int64]domain.Location),
		Mutex: sync.RWMutex{},
	}

	//basic validations
	if jobRequest.Department != 0 {
		wg.Add(1)
		go service.supplierService.GetDepartment(jobRequest.Department, department, &wg)
	}

	if jobRequest.JobItem != 0 {
		wg.Add(1)
		go service.supplierService.GetJobItem(jobRequest.JobItem, jobItem, &wg)
	}

	for _, v := range jobRequest.Locations {
		wg.Add(1)
		go service.supplierService.GetLocation(v, location, &wg, errorMessageMap, locationMap)
	}

	wg.Wait()

	//1. If the given department is not null and does not exist in the property then return a bad request
	if department.ID == 0 {
		buffer.WriteString(fmt.Sprintf(`Department with id %d not found;`, jobRequest.Department))
		errorCount++
	}

	//2. If the given job_item is not null and does not exist in the property then return a bad request
	if jobItem.Id == 0 {
		buffer.WriteString(fmt.Sprintf(`JobItem with id %d not found;`, jobRequest.JobItem))
		errorCount++
	}

	//3. If the given locations are not null and do not exist in the property then return a bad request
	if len(errorMessageMap.M) > 0 {
		for _, v := range errorMessageMap.M {
			buffer.WriteString(v)
		}
	}

	//if basic validation failed, then return error messages
	if errorCount > 0 || len(errorMessageMap.M) >= 1 {
		return nil, errors.NewBadRequestError(buffer.String())
	}
	//end of basic validations
	var locations = make([]domain.Location, 0)

	for _, l := range locationMap.M {
		locations = append(locations, l)
	}

	//4. If the given department is ‘Housekeeping’, the given job_item is one of ['Blanket','Sheets','Mattress'],
	//and a location with a location type of 'Room' is given then create a job to clean the bed(s) in the given room
	//* If a location with a location type of 'Floor' is given then create a job to clean the beds in all rooms with a
	//location type of ‘Room’ on that floor.
	//* If the given location is not of location type ‘Room’ or ‘Floor’ then return a bad request.
	//ATTENTION - JobItem with displayName = 'Sheets' not found!
	var jobCleanBed, errCleanBed = service.cleanBedTaskCreateJob(jobRequest, jobItem.DisplayName, locations)

	if errCleanBed != nil {
		return nil, errCleanBed
	}

	if jobCleanBed != nil {
		return jobCleanBed, nil
	}

	// 5. If the given department is ‘Engineering’ and a job_item and at least one location is given
	// then create a job to repair the given job item at the given location(s).
	var jobRepairTask, errRepair = service.repairTaskCreateJob(jobRequest, jobItem.DisplayName, locations)

	if errRepair != nil {
		return nil, errRepair
	}

	if jobRepairTask != nil {
		return jobRepairTask, nil
	}

	// 6. If the department is ‘Room Service’ and a job item and at least 1 location is given then create a job to deliver
	// that job item to the given locations.
	//* If only one location is given and is of location type ‘Floor’ then create a job to deliver the given job item in
	//all locations with a location type of 'Room' on that floor.
	var jobDeliver, errDeliver = service.deliverTaskCreateJob(jobRequest, jobItem.DisplayName, locations)
	if errDeliver != nil {
		return nil, errDeliver
	}

	if jobDeliver != nil {
		return jobDeliver, nil
	}

	//7. If none of the above criteria is met return a bad request.
	return nil, errors.NewBadRequestError("Request hasn't matched any expected criteria. Please, review your request")
}

// Get all children locations from the parentLocation and add it to the job
func addChildrenFromLocation(addToJob *domain.Job, allLocations []domain.Location,
	parentLocation domain.Location) (*domain.Job, *errors.ApplicationError) {
	children, errChildren := interfaces.GetLocationsChildrenFrom(allLocations, parentLocation, locationTypeRoom)

	if len(children) == 0 {
		return nil, errors.NewNotFoundError("Location without rooms inside")
	}

	for _, c := range children {
		addToJob.Location = append(addToJob.Location, domain.LocationRequest{Id: c.Id})
	}

	if errChildren != nil {
		return nil, errChildren
	}

	return addToJob, nil
}

func (service DefaultJobService) deliverTaskCreateJob(jobRequest dto.JobRequestDTO, jobItemName string,
	locations []domain.Location) (*domain.JobResponse, *errors.ApplicationError) {
	if isDeliverTask(jobRequest) {
		var job = getJobToDeliver(jobRequest, jobItemName, locations)

		//* If only one location is given and is of location type ‘Floor’ then create a job to deliver the given job item
		//in all locations with a location type of 'Room' on that floor.
		if len(locations) == 1 && locations[0].LocationType.Id == locationTypeFloor {
			//clean all
			job.Location = make([]domain.LocationRequest, 0)
			allLocations, errLocations := service.supplierService.GetAllLocations()

			if errLocations != nil {
				return nil, errLocations
			}

			var newJob, errJob = addChildrenFromLocation(job, allLocations, locations[0])

			if errJob != nil {
				return nil, errJob
			}

			job = newJob
		}

		var returnJob, err = service.supplierService.CreateJob(job, nil)

		if err != nil {
			return nil, errors.NewUnexpectedError(err.Message)
		}

		return returnJob, nil
	}

	return nil, nil
}

func (service DefaultJobService) repairTaskCreateJob(jobRequest dto.JobRequestDTO, jobItemName string,
	locations []domain.Location) (*domain.JobResponse, *errors.ApplicationError) {

	if isRepairTask(jobRequest, locations) {
		var job = getJobToRepair(jobRequest, jobItemName, locations)

		//* If only one location is given and is of location type ‘Floor’ then create a job to repair the given
		//job item in all locations on that floor.
		if len(locations) == 1 && locations[0].LocationType.Id == locationTypeFloor {
			allLocations, errLocations := service.supplierService.GetAllLocations()

			if errLocations != nil {
				return nil, errLocations
			}

			var newJob, errJob = addChildrenFromLocation(job, allLocations, locations[0])

			if errJob != nil {
				return nil, errJob
			}

			job = newJob
		}

		var returnJob, err = service.supplierService.CreateJob(job, nil)

		if err != nil {
			return nil, errors.NewUnexpectedError(err.Message)
		}

		return returnJob, nil
	}

	return nil, nil
}
func (service DefaultJobService) cleanBedTaskCreateJob(jobRequest dto.JobRequestDTO, jobItemName string, locations []domain.Location) (*domain.JobResponse,
	*errors.ApplicationError) {
	if isCleanBedTask(jobRequest, locations) {
		for _, l := range locations {
			//* If the given location is not of location type ‘Room’ or ‘Floor’ then return a bad request.
			if l.LocationType.Id != locationTypeFloor && l.LocationType.Id != locationTypeRoom {
				return nil, errors.NewBadRequestError("With a department HouseKeeping, locationType must be Room or Floor")
			}
		}

		var job = getJobToCleanBed(jobRequest, jobItemName, locations)

		for _, l := range locations {
			if len(locations) == 1 && l.LocationType.Id == locationTypeFloor {
				//repair all
				allLocations, errLocations := service.supplierService.GetAllLocations()

				if errLocations != nil {
					return nil, errLocations
				}

				var newJob, errJob = addChildrenFromLocation(job, allLocations, locations[0])

				if errJob != nil {
					return nil, errJob
				}

				job = newJob
			}
		}

		var jobReturn, err = service.supplierService.CreateJob(job, nil)

		if err != nil {
			return nil, err
		}

		return jobReturn, nil
	}

	return nil, nil
}

func getJobToCleanBed(jobRequest dto.JobRequestDTO, jobItemName string, locations []domain.Location) *domain.Job {
	return getJobFromDTO(jobRequest, jobItemName, locations, "clean")
}

func getJobToRepair(jobRequest dto.JobRequestDTO, jobItemName string, locations []domain.Location) *domain.Job {
	return getJobFromDTO(jobRequest, jobItemName, locations, "repair")
}

func getJobToDeliver(jobRequest dto.JobRequestDTO, jobItemName string, locations []domain.Location) *domain.Job {
	return getJobFromDTO(jobRequest, jobItemName, locations, "deliver")

}

func getJobFromDTO(jobRequest dto.JobRequestDTO, jobItemName string, locations []domain.Location, action string) *domain.Job {
	var department = domain.Department{
		ID: jobRequest.Department,
	}

	var locationRequestList = make([]domain.LocationRequest, 0)

	for _, v := range locations {
		var location = domain.LocationRequest{
			Id: v.Id,
		}
		locationRequestList = append(locationRequestList, location)
	}

	var job = domain.Job{
		Item: domain.JobItem{
			Id:          jobRequest.JobItem,
			DisplayName: jobItemName,
		},
		Location:   locationRequestList,
		Department: department,
		Action:     action,
	}

	return &job
}
func isCleanBedTask(job dto.JobRequestDTO, locations []domain.Location) bool {
	var isHouseKeeping = job.Department == departmentHouseKeeping
	var hasJobItems = len(houseKeepingItems[job.JobItem])

	if isHouseKeeping && hasJobItems > 0 {
		for _, l := range locations {
			if l.LocationType.Id == locationTypeFloor || l.LocationType.Id == locationTypeRoom {
				return true
			}
		}
	}

	return false
}

func isRepairTask(job dto.JobRequestDTO, locations []domain.Location) bool {
	if job.Department == departmentEngineering && job.JobItem != 0 && len(locations) >= 1 {
		return true
	}

	return false
}

// 6. If the department is ‘Room Service’ and a job item and at least 1 location is given then create a job to deliver
// that job item to the given locations.
func isDeliverTask(job dto.JobRequestDTO) bool {
	if job.Department == departmentRoomService && job.JobItem != 0 && len(job.Locations) >= 1 {
		return true
	}

	return false
}
func NewJobService() DefaultJobService {
	return DefaultJobService{supplierService: DefaultSupplierService{authService: DefaultAuthService{}}}
}
