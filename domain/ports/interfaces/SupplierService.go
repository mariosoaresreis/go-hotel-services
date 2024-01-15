package interfaces

import (
	"github.com/mariosoaresreis/go-hotel/domain"
	"github.com/mariosoaresreis/go-hotel/errors"
	"sync"
)

type MessageMap struct {
	Mutex sync.RWMutex
	M     map[int64]string
}

type LocationMap struct {
	Mutex sync.RWMutex
	M     map[int64]domain.Location
}

//go:generate mockgen -destination=../../../domain/adapters/mocks/mockSupplierService.go -package=services github.com/mariosoaresreis/go-hotel/domain/ports/interfaces SupplierService
type SupplierService interface {
	GetDepartment(DepartmentID int64, target *domain.Department, wg *sync.WaitGroup) (*domain.Department, *errors.ApplicationError)
	GetLocation(LocationID int64, target *domain.Location, wg *sync.WaitGroup, errorMessageMap MessageMap,
		locationMap LocationMap) (*domain.Location, *errors.ApplicationError)
	GetJobItem(JobItem int64, target *domain.JobItemResponse, wg *sync.WaitGroup) (*domain.JobItemResponse, *errors.ApplicationError)
	CreateJob(job *domain.Job, wg *sync.WaitGroup) (*domain.JobResponse, *errors.ApplicationError)
	GetAllLocations() ([]domain.Location, *errors.ApplicationError)
}

func GetLocationsChildrenFrom(locations []domain.Location, fromLocation domain.Location,
	locationType int64) ([]domain.Location, *errors.ApplicationError) {

	childrenLocations := make([]domain.Location, 0)

	for _, l := range locations {
		if l.ParentLocation.Id == fromLocation.Id && l.LocationType.Id == locationType {
			childrenLocations = append(childrenLocations, l)
		}
	}

	return childrenLocations, nil
}
