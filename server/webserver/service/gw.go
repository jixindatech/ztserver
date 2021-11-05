package service

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"zt-server/storage"
)

const TotalStaticsLength = 11

type Gw struct {
	User     string
	Method   string
	Resource string
	Start    int64
	End      int64

	Page     int
	PageSize int
}

func (g *Gw) GetList() (interface{}, int64, error) {
	query := make(map[string]interface{})
	query["user"] = g.User
	query["method"] = g.Method
	query["resource"] = g.Resource
	query["start"] = g.Start
	query["end"] = g.End

	data, err := storage.Query(storage.GW_INDEX, query, g.Page, g.PageSize)
	if err != nil {
		return nil, 0, err
	}

	res := make(map[string]interface{})
	res["count"] = data["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"]
	res["data"] = data["hits"].(map[string]interface{})["hits"]

	return res["data"], int64(res["count"].(float64)), nil
}

type intBucketItem struct {
	Key      float64 `json:"key"`
	DocCount int     `json:"doc_count"`
}

type intBucketArray struct {
	Buckets []intBucketItem `json:"buckets"`
}

type userBucket struct {
	Key          string         `json:"key"`
	DocCount     int            `json:"doc_count"`
	IntervalData intBucketArray `json:"interval_data"`
}

func (g *Gw) GetStatics() (*[TotalStaticsLength]int, error) {
	query := make(map[string]interface{})
	query["start"] = g.Start
	query["end"] = g.End
	data, err := storage.QueryInfo(storage.GW_INDEX, query)
	if err != nil {
		return nil, err
	}

	statics := new([TotalStaticsLength]int) // golang array
	users := gjson.GetBytes(data, "aggregations.group_user.buckets")
	if users.Exists() {
		var items []userBucket
		err := json.Unmarshal([]byte(users.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			for index, item1 := range item.IntervalData.Buckets {
				if item1.DocCount > 0 {
					statics[index]++
				}
			}
		}
	}

	return statics, nil
}
