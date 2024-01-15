package services

import (
	"encoding/json"
	"fmt"
	"github.com/mariosoaresreis/go-hotel/domain"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type DefaultAuthService struct {
}

func (DefaultAuthService) GetToken() (*domain.TokeResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	const ClientSecret = "client_secret"
	const ClientId = "client_id"
	const GrantType = "grant_type"
	const Scope = "scope"

	data := url.Values{}
	data.Add(ClientSecret, os.Getenv(ClientSecret))
	data.Add(ClientId, os.Getenv(ClientId))
	data.Add(GrantType, os.Getenv(GrantType))
	data.Add(Scope, os.Getenv(Scope))
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", os.Getenv("AUTH_URL"), strings.NewReader(encodedData))

	if err != nil {
		return nil, fmt.Errorf("error %s", err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(req)
	var target *domain.TokeResponse

	if err != nil {
		return nil,
			fmt.Errorf("error %s", err.Error())
	}

	defer response.Body.Close()
	var e = json.NewDecoder(response.Body).Decode(&target)

	if e != nil {
		panic(e)
	}

	return target, nil
}
