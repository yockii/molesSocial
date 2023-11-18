package domain

type WellKnownResponse struct {
	Subject string          `json:"subject,omitempty"`
	Aliases []string        `json:"aliases,omitempty"`
	Links   []WellKnownLink `json:"links,omitempty"`
}

type WellKnownLink struct {
	Rel  string `json:"rel,omitempty"`
	Type string `json:"type,omitempty"`
	Href string `json:"href,omitempty"`
}
