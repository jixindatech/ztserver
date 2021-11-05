package service

import (
	"zt-server/storage"
)

const WsIndex = "ztserver"

type Ws struct {
	User     string
	Resource string
	Start    int64
	End      int64

	Page     int
	PageSize int
}

func (w *Ws) GetList() (interface{}, int64, error) {
	query := make(map[string]interface{})
	query["user"] = w.User
	query["resource"] = w.Resource
	query["start"] = w.Start
	query["end"] = w.End

	data, err := storage.Query(storage.WS_INDEX, query, w.Page, w.PageSize)
	if err != nil {
		return nil, 0, err
	}

	res := make(map[string]interface{})
	res["count"] = data["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"]
	res["data"] = data["hits"].(map[string]interface{})["hits"]

	return res["data"], int64(res["count"].(float64)), nil
}
