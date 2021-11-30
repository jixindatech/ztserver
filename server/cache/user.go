package cache

type Routes struct {
	Host string `json:"host"`
	Path string `json:"path"`
	Methods []string `json:"methods"`
}

type UserCache struct {
	ID        uint64         `json:"id"`
	Name      string         `json:"name"`
	Secret    string         `json:"secret"`
	Resources []Routes       `json:"resources"`
}

type GwCache struct {
	Timestamp int64 `json:"timestamp"`
	Values []*UserCache `json:"values"`
}