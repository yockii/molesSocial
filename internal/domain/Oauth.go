package domain

type OauthAuthorizeRequest struct {
	ResponseType string `json:"response_type" form:"response_type" query:"response_type"`
	ClientID     string `json:"client_id" form:"client_id" query:"client_id"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri" query:"redirect_uri"`
	Scope        string `json:"scope" form:"scope" query:"scope"`
	State        string `json:"state" form:"state" query:"state"`
}

type OauthTokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	Code         string `json:"code" form:"code" query:"code"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri" query:"redirect_uri"`
	ClientID     string `json:"client_id" form:"client_id" query:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
}

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	CreatedAt   int64  `json:"created_at"`
}
