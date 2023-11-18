package domain

type ApplicationCreateRequest struct {
	ClientName   string `json:"client_name"`
	RedirectURIs string `json:"redirect_uris"`
	Scopes       string `json:"scopes"`
	Website      string `json:"website"`
}
