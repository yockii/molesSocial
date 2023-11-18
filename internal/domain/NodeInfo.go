package domain

type NodeInfo struct {
	Version           string                 `json:"version"`
	Software          *NodeInfoSoftware      `json:"software"`
	Protocols         []string               `json:"protocols"`
	Services          *NodeInfoServices      `json:"services"`
	OpenRegistrations bool                   `json:"openRegistrations"`
	Usage             *NodeInfoUsage         `json:"usage"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type NodeInfoSoftware struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type NodeInfoServices struct {
	Inbound  []string `json:"inbound"`
	Outbound []string `json:"outbound"`
}

type NodeInfoUsage struct {
	Users      *NodeInfoUsers `json:"users"`
	LocalPosts int64          `json:"localPosts"`
}

type NodeInfoUsers struct {
	Total          int64 `json:"total"`
	ActiveMonth    int64 `json:"activeMonth"`
	ActiveHalfyear int64 `json:"activeHalfyear"`
}

type InstanceInfo struct {
	URI              string            `json:"uri"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Email            string            `json:"email"`
	Version          string            `json:"version"`
	Urls             map[string]string `json:"urls"`
	Stats            map[string]int64  `json:"stats"`
	ShortDescription string            `json:"short_description"`
	Thumbnail        string            `json:"thumbnail"`
	Languages        []string          `json:"languages"`
	ContactAccount   *AccountInfo      `json:"contact_account"`
}
