package dto

type JobRequestDTO struct {
	Department int64   `json:"department"`
	JobItem    int64   `json:"job_item"`
	Locations  []int64 `json:"locations"`
}
