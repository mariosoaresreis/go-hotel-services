package domain

import (
	"time"
)

type Job struct {
	Item       JobItem           `json:"item"`
	Action     string            `json:"action"`
	Department Department        `json:"department"`
	Location   []LocationRequest `json:"location"`
}

type JobResponse struct {
	Id          int               `json:"id"`
	Item        JobItemResponse   `json:"item"`
	DisplayName string            `json:"displayName"`
	Type        string            `json:"type"`
	Priority    string            `json:"priority"`
	Action      string            `json:"action"`
	Department  []Department      `json:"departments"`
	Location    []LocationRequest `json:"locations"`
	Roles       []Role            `json:"role"`
	Notes       []Note            `json:"notes"`
	Assignee    Employee          `json:"assignee"`
	DueBy       time.Time         `json:"dueBy"`
}
