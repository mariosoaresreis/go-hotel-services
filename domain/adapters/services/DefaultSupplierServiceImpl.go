package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mariosoaresreis/go-hotel/domain"
	"github.com/mariosoaresreis/go-hotel/domain/ports/interfaces"
	"github.com/mariosoaresreis/go-hotel/errors"
	"github.com/mariosoaresreis/go-hotel/logger"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type DefaultSupplierService struct {
	authService interfaces.AuthService
}

func (supplierService DefaultSupplierService) GetDepartment(DepartmentID int64, target *domain.Department, wg *sync.WaitGroup) (*domain.Department, *errors.ApplicationError) {
	var err = supplierService.getBody(fmt.Sprintf("/departments/%d", DepartmentID), http.StatusOK, target, wg)

	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}

	return target, nil
}
func (supplierService DefaultSupplierService) GetLocation(LocationID int64, target *domain.Location, wg *sync.WaitGroup,
	errorMessageMap interfaces.MessageMap, locationMap interfaces.LocationMap) (*domain.Location, *errors.ApplicationError) {
	var err = supplierService.getLocation(fmt.Sprintf("/locations/%d", LocationID), http.StatusOK, target, wg, errorMessageMap, locationMap)

	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}

	return target, nil
}
func (supplierService DefaultSupplierService) GetJobItem(JobItem int64, target *domain.JobItemResponse, wg *sync.WaitGroup) (*domain.JobItemResponse, *errors.ApplicationError) {
	var err = supplierService.getBody(fmt.Sprintf("/jobitems/%d", JobItem), http.StatusOK, target, wg)

	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}

	return target, nil
}

func (supplierService DefaultSupplierService) getHttpResponse(urlMethod string) (*http.Response, error) {
	url := os.Getenv("OPTII_URL")
	var token, e = supplierService.authService.GetToken()

	if e != nil {
		errorMessage := fmt.Sprintf("Could not retrieve token. Error %s", e.Error())
		logger.Fatal(errorMessage)
		panic(errorMessage)
	}

	client := &http.Client{Timeout: time.Second * 10}
	var bearer = fmt.Sprintf("Bearer %s", token.AccessToken)
	var completeUrl = fmt.Sprintf("%s%s", url, urlMethod)
	req, err := http.NewRequest("GET", completeUrl, nil)
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	return resp, nil
}

func (supplierService DefaultSupplierService) postHttpResponse(urlMethod string, body interface{}) (*http.Response, error) {
	url := os.Getenv("OPTII_URL")
	var token, e = supplierService.authService.GetToken()

	if e != nil {
		logger.Fatal(fmt.Sprintf("Could not retrieve token. Error %s", e))
	}

	var marshalled, err = json.Marshal(body)

	if err != nil {
		log.Fatalf("impossible to marshall object: %s", err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	var bearer = fmt.Sprintf("Bearer %s", token.AccessToken)
	var completeUrl = fmt.Sprintf("%s%s", url, urlMethod)
	req, err := http.NewRequest("POST", completeUrl, bytes.NewReader(marshalled))
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	return resp, nil
}

func (supplierService DefaultSupplierService) getBody(urlMethod string, expectedHttpStatus int, target interface{}, wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}

	var resp, err = supplierService.getHttpResponse(urlMethod)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	// if we want to check for a specific status code, we can do so here
	// for example, a successful request should return the expected status
	if resp.StatusCode != expectedHttpStatus {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBytes := buf.String()
		errorMessage := fmt.Sprintf("Error from supplier service from method %s. Code %d Message %s", urlMethod,
			resp.StatusCode, respBytes)

		if resp.StatusCode >= http.StatusInternalServerError {
			logger.Fatal(errorMessage)
			panic(errorMessage)
		} else {
			//expected by the business rule when object is not found
			if resp.StatusCode == http.StatusNotFound {
				target = nil
				return nil
			}

			logger.Error(errorMessage)
		}

		target = nil
	}

	json.NewDecoder(resp.Body).Decode(&target)
	return nil
}

func (supplierService DefaultSupplierService) postBody(urlMethod string, expectedHttpStatus int, body interface{}, target interface{}, wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}

	var resp, err = supplierService.postHttpResponse(urlMethod, body)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	// if we want to check for a specific status code, we can do so here
	// for example, a successful request should return the expected status
	if resp.StatusCode != expectedHttpStatus {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBytes := buf.String()
		errorMessage := fmt.Sprintf("Error from supplier service from method %s. Code %d Message %s", urlMethod,
			resp.StatusCode, respBytes)

		if resp.StatusCode >= http.StatusInternalServerError {
			logger.Fatal(errorMessage)
			panic(errorMessage)
		} else {
			//expected by the business rule when object is not found
			if resp.StatusCode == http.StatusNotFound {
				target = nil
				return nil
			}

			logger.Error(errorMessage)
		}

		target = nil
	}

	json.NewDecoder(resp.Body).Decode(&target)
	return nil
}
func (supplierService DefaultSupplierService) getLocation(urlMethod string, httpStatus int, target *domain.Location, wg *sync.WaitGroup,
	messageMap interfaces.MessageMap, locationMap interfaces.LocationMap) error {
	defer wg.Done()

	var resp, err = supplierService.getHttpResponse(urlMethod)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	// if we want to check for a specific status code, we can do so here
	// for example, a successful request should return the expected status
	if resp.StatusCode != httpStatus {
		target = nil
	}

	if resp.StatusCode == http.StatusNotFound {
		target = nil
		return nil
	}

	json.NewDecoder(resp.Body).Decode(&target)

	if len(target.Name) == 0 {
		WriteMessage(messageMap, target.Id, fmt.Sprintf("Location with Id %d not found;", target.Id))
	} else {
		StoreLocation(locationMap, target.Id, *target)
	}

	return nil
}

func (supplierService DefaultSupplierService) CreateJob(job *domain.Job, wg *sync.WaitGroup) (*domain.JobResponse,
	*errors.ApplicationError) {
	var result *domain.JobResponse
	var err = supplierService.postBody("/jobs", http.StatusOK, job, &result, nil)

	if err != nil {
		return nil, errors.NewUnexpectedError(err.Error())
	}

	return result, nil
}

func (supplierService DefaultSupplierService) GetAllLocations() ([]domain.Location, *errors.ApplicationError) {
	hasNextPage := true
	first := 0
	locations := make([]domain.Location, 0)

	for hasNextPage {
		var locationPage = new(domain.LocationPage)
		var e = supplierService.getBody(fmt.Sprintf("/locations?next=100&first=%d", first), http.StatusOK,
			&locationPage, nil)

		if e != nil {
			logger.Error(e.Error())
			return nil, errors.NewUnexpectedError(e.Error())
		}

		first += 100
		for _, v := range locationPage.Items {
			locations = append(locations, v)
		}
		hasNextPage = locationPage.PageInfo.HasNextPage
	}
	return locations, nil
}

func WriteMessage(o interfaces.MessageMap, key int64, value string) {
	o.Mutex.Lock()         // lock for writing, blocks until the Mutex is ready
	defer o.Mutex.Unlock() // again, make SURE you do this, else it will be locked permanently
	o.M[key] = value
}

func StoreLocation(o interfaces.LocationMap, key int64, value domain.Location) {
	o.Mutex.Lock()         // lock for writing, blocks until the Mutex is ready
	defer o.Mutex.Unlock() // again, make SURE you do this, else it will be locked permanently
	o.M[key] = value
}
