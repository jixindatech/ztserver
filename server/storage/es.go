package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
	"go.uber.org/zap"
	"zt-server/pkg/core/golog"
	"zt-server/settings"
)

const maxItemGroup = 10

type EsStorage struct {
	config  *settings.Es
	client  *elasticsearch.Client
	GwIndex string
	WsIndex string
}

var es EsStorage
var WS_INDEX string = "wsindex"
var GW_INDEX string = "gwindex"

func InitStorage(cfg *settings.Es) error {
	var err error

	if es.client != nil {
		return nil
	}
	if len(cfg.GwIndex) == 0 {
		return fmt.Errorf("%s", "invalid gw index")
	}
	if len(cfg.WsIndex) == 0 {
		return fmt.Errorf("%s", "invalid ws index")
	}

	es.GwIndex = cfg.GwIndex
	es.WsIndex = cfg.WsIndex

	es.config = cfg
	es.client, err = elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: es.config.Host,
			Username:  es.config.User,
			Password:  es.config.Password,
		},
	)
	if err != nil {
		return err
	}

	res, err := es.client.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func Save(index string, body interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}
	var esIndex string
	if index == GW_INDEX {
		esIndex = es.GwIndex
	} else if index == WS_INDEX {
		esIndex = es.WsIndex
	} else {
		return fmt.Errorf("%s", "unkonwn es index")
	}

	res, err := es.client.Index(esIndex, &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func Query(index string, _query map[string]interface{}, page, pageSize int) (map[string]interface{}, error) {
	var queries []esquery.Mappable

	var esIndex string
	if index == GW_INDEX {
		esIndex = es.GwIndex
	} else if index == WS_INDEX {
		esIndex = es.WsIndex
	} else {
		return nil, fmt.Errorf("%s", "unkonwn es index")
	}

	if _query["user"] != nil {
		user := _query["user"].(string)
		if len(user) > 0 {
			queries = append(queries, esquery.Term("user", user))
		}
	}

	if _query["method"] != nil {
		method := _query["method"].(string)
		if len(method) > 0 {
			queries = append(queries, esquery.Match("method", method))
		}
	}

	if _query["resource"] != nil {
		resource := _query["resource"].(string)
		if len(resource) > 0 {
			queries = append(queries, esquery.Term("resource", resource))
		}
	}

	start := _query["start"].(int64) / 1000
	end := _query["end"].(int64) / 1000
	if start > 0 && end > 0 && end > start {
		queries = append(queries, esquery.Range("time").Gte(start).Lte(end))
	} else {
		return nil, errors.New("query time has error")
	}

	query := esquery.Search().Sort("time", "desc").Query(
		esquery.Bool().Must(queries...),
	)
	res, err := query.Run(
		es.client,
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex(esIndex),
		es.client.Search.WithFrom((page-1)*pageSize),
		es.client.Search.WithSize(pageSize),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		golog.Error("es", zap.String("err", err.Error()))
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("elasticsearch body has error")
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}

func QueryInfo(index string, _query map[string]interface{}) ([]byte, error) {
	var queries []esquery.Mappable

	var esIndex string
	if index == GW_INDEX {
		esIndex = es.GwIndex
	} else if index == WS_INDEX {
		esIndex = es.WsIndex
	} else {
		return nil, fmt.Errorf("%s", "unkonwn es index")
	}

	start := _query["start"].(int64) / 1000
	end := _query["end"].(int64) / 1000
	if start > 0 && end > 0 && end > start {
		queries = append(queries, esquery.Range("time").Gte(start).Lte(end))
	} else {
		return nil, errors.New("query time has error")
	}

	interval := (end - start) / maxItemGroup

	res := esquery.Search().Query(esquery.Bool().Must(queries...))
	res = res.Aggs(
		esquery.TermsAgg("group_user", "user").Aggs(
			esquery.CustomAgg("interval_data", map[string]interface{}{
				"histogram": map[string]interface{}{
					"field":    "time",
					"interval": interval,
					"extended_bounds": map[string]interface{}{
						"min": start,
						"max": end},
				},
			})),
	)

	response, err := res.Size(0).Run(
		es.client,
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex(esIndex),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, errors.New("elasticsearch body has error")
	}

	var b bytes.Buffer
	_, err = b.ReadFrom(response.Body)

	return b.Bytes(), err
}

func QueryStatics(index string, _query map[string]interface{}) ([]byte, error) {

	return nil, nil
}
