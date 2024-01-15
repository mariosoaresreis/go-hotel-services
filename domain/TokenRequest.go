package domain

type TokenRequest struct {
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}
