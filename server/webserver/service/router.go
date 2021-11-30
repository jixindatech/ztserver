package service

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
	"zt-server/cache"
	"zt-server/pkg/core/golog"
	"zt-server/webserver/model"
)

const cacheRouterName string = "router"

type Router struct {
	ID uint64

	Name string
	Host string
	Path string
	UpstreamRef uint64
	Remark  string

	Page     int
	PageSize int
}

func (r *Router) Save() (err error) {
	data := make(map[string]interface{})
	data["host"] = r.Host
	data["path"] = r.Path
	data["upstreamRef"] = r.UpstreamRef
	data["remark"] = r.Remark

	if r.ID > 0 {
		err = model.PutRouter(r.ID, data)
	} else {
		data["name"] = r.Name
		err = model.AddRouter(data)
	}

	if err != nil {
		return err
	}

	return SetupRouters()
}

func (r *Router) Get() (*model.Router, error) {
	return model.GetRouter(r.ID)
}

func (r *Router) GetList() ([]*model.Router, int, error) {
	data := make(map[string]interface{})
	data["name"] = r.Name
	data["page"] = r.Page
	data["pagesize"] = r.PageSize

	return model.GetRouters(data)
}

func (r *Router) Delete() error {
	err := model.DeleteRouter(r.ID)
	if err != nil {
		return err
	}

	return SetupRouters()
}

func SetupRouters() error {
	router := Router{}
	routers, count, err := router.GetList()
	if err != nil {
		return err
	}
	if count == 0 {
		data := make(map[string]interface{})
		data["values"] = [][]struct{}{}
		data["timestamp"] = time.Now().Unix()

		routerStr, err := json.Marshal(data)
		if err != nil {
			golog.Error("router", zap.String("err", err.Error()))
			return err
		}

		err = cache.Set(cacheRouterName, string(routerStr))
		if err != nil {
			golog.Error("router", zap.String("err", err.Error()))
			return err
		}
		return nil
	}

	routesInfos := []map[string]interface{}{}
	for _, item := range routers {
		route := make(map[string]interface{})
		route["id"] = item.ID
		route["host"] = item.Host
		route["uri"] = item.Path
		if len(item.Upstreams) != 1 {
			return fmt.Errorf("%s", "invalid router upstream")
		}

		route["upstream_id"] = item.Upstreams[0].ID
		routesInfos = append(routesInfos, route)
	}

	data := make(map[string]interface{})
	data["values"] = routesInfos
	data["timestamp"] = time.Now().Unix()

	routerStr, err := json.Marshal(data)
	if err != nil {
		golog.Error("router", zap.String("err", err.Error()))
		return err
	}

	err = cache.Set(cacheRouterName, string(routerStr))
	if err != nil {
		golog.Error("router", zap.String("err", err.Error()))
		return err
	}

	return nil
}
