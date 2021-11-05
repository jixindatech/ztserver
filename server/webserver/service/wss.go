package service

import (
	"sync/atomic"
	"zt-server/ssh"
	"zt-server/storage"
	"zt-server/webserver/model"
)

var userOnopenCount int64 = 0

type WssEvent struct {
	ID uint

	User      string
	Email     string
	Client    string
	Event     string
	Time      int64
	Dev       string
	Resources []*model.Resource
	Gw        []string

	Page     int
	PageSize int
}

func (w *WssEvent) Add(index string) error {
	event := map[string]interface{}{
		"user":   w.User,
		"email":  w.Email,
		"client": w.Client,
		"event":  w.Event,
		"dev":    w.Dev,
		"time":   w.Time,
		"gw":     w.Gw,
	}

	return storage.Save(index, event)
}

func (w *WssEvent) OpenDoor(servers map[string]int) ([]string, error) {
	atomic.AddInt64(&userOnopenCount, 1)

	ips := ssh.OpenConnection(w.Client, servers)
	return ips, nil
}

func (w *WssEvent) CloseDoor(servers map[string]int) ([]string, error) {
	atomic.AddInt64(&userOnopenCount, -1)

	ips := ssh.CloseConnection(w.Client, servers)
	return ips, nil
}

func GetUserOnopenCount() int64 {
	pv := atomic.LoadInt64(&userOnopenCount)
	return pv
}
