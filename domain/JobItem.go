package domain

type JobItem struct {
	Id          int64  `json:"id"`
	DisplayName string `json:"name"`
}
type JobItemResponse struct {
	Id          int64  `json:"id"`
	DisplayName string `json:"displayName"`
}
