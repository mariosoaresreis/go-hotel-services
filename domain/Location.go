package domain

type Location struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	ParentLocation struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"parentLocation"`

	LocationType LocationType `json:"locationType"`
}

type LocationRequest struct {
	Id int64 `json:"id"`
}
