package domain

type LocationPageInfo struct {
	EndCursor   int  `json:"endCursor"`
	HasNextPage bool `json:"hasNextPage"`
	TotalCount  int  `json:"totalCount"`
}
