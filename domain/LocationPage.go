package domain

type LocationPage struct {
	PageInfo LocationPageInfo `json:"pageInfo"`
	Items    []Location       `json:"items"`
}
