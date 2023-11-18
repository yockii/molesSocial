package domain

type OauthAuthorizeRequest struct {
	ResponseType string `json:"response_type" form:"response_type" query:"response_type"`
	ClientID     string `json:"client_id" form:"client_id" query:"client_id"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri" query:"redirect_uri"`
	Scope        string `json:"scope" form:"scope" query:"scope"`
	State        string `json:"state" form:"state" query:"state"`
}
