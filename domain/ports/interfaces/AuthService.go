package interfaces

import "github.com/mariosoaresreis/go-hotel/domain"

type AuthService interface {
	GetToken() (*domain.TokeResponse, error)
}
