package interfaces

import (
	"github.com/mariosoaresreis/go-hotel/domain"
	"github.com/mariosoaresreis/go-hotel/domain/dto"
	"github.com/mariosoaresreis/go-hotel/errors"
)

//go:generate mockgen -destination=../../../domain/adapters/mocks/mockJobService.go -package=services github.com/mariosoaresreis/go-hotel/domain/ports/interfaces JobService
type JobService interface {
	AddJob(jobRequest dto.JobRequestDTO) (*domain.JobResponse, *errors.ApplicationError)
}
