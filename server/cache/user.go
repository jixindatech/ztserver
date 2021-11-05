package cache

type UserCache struct {
	ID        uint64         `json:"id"`
	Timestamp int64          `json:"timestamp"`
	Resource  map[string]int `json:"resource"`
}
