package constant

const (
	NodeInfoVersion      = "2.0"
	NodeInfoRel          = "http://nodeinfo.diaspora.software/ns/schema/" + NodeInfoVersion
	NodeInfoSoftwareName = "molesocial"
)

var (
	Version = "0.0.1"

	NodeInfoProtocols = []string{"activitypub"}
	NodeInfoInbound   = []string{}
	NodeInfoOutbound  = []string{}
	NodeInfoMetadata  = make(map[string]interface{})
)
