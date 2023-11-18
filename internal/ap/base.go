package ap

const (
	// ActivityPub context
	ActivityPubContext = "https://www.w3.org/ns/activitystreams"
	// ActivityPub context public
	ActivityPubContextPublic = "https://www.w3.org/ns/activitystreams#Public"
)

type AcitivityPubContextLanguage struct {
	AtLanguage string `json:"@language"`
}

type Base struct {
	AtContext []interface{} `json:"@context"`
}
