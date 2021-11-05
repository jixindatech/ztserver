package service

import (
	"encoding/json"
	"go.uber.org/zap"
	"time"
	"zt-server/cache"
	"zt-server/pkg/core/golog"
	"zt-server/webserver/model"
)

const cacheProxyName string = "balancer"

type Proxy struct {
	ID uint64

	Server  string
	Lb      string
	Backend []byte
	Remark  string

	Page     int
	PageSize int
}

func (p *Proxy) Save() (err error) {
	data := make(map[string]interface{})
	data["server"] = p.Server
	data["lb"] = p.Lb
	data["backend"] = p.Backend
	data["remark"] = p.Remark

	if p.ID > 0 {
		err = model.PutProxy(p.ID, data)
	} else {
		err = model.AddProxy(data)
	}

	if err != nil {
		return err
	}

	return SetupProxys()
}

func (p *Proxy) Get() (*model.Proxy, error) {
	return model.GetProxy(p.ID)
}

func (p *Proxy) GetList() ([]*model.Proxy, int, error) {
	data := make(map[string]interface{})
	data["server"] = p.Server
	data["page"] = p.Page
	data["pagesize"] = p.PageSize

	return model.GetProxys(data)
}

func (p *Proxy) Delete() error {
	err := model.DeleteProxy(p.ID)
	if err != nil {
		return err
	}

	return SetupProxys()
}

func SetupProxys() error {
	proxy := Proxy{}
	proxys, count, err := proxy.GetList()
	if err != nil {
		return err
	}
	if count == 0 {
		data := make(map[string]interface{})
		data["data"] = struct{}{}
		data["timestamp"] = time.Now().Unix()

		proxyStr, err := json.Marshal(data)
		if err != nil {
			golog.Error("proxy", zap.String("err", err.Error()))
			return err
		}

		err = cache.Set(cacheProxyName, string(proxyStr))
		if err != nil {
			golog.Error("proxy", zap.String("err", err.Error()))
			return err
		}
		return nil
	}
	backends := make(map[string]interface{})
	for _, item := range proxys {
		backends[item.Server] =
			map[string]interface{}{
				"upstreams": item.Backend,
				"schedule":  item.Lb,
			}
	}

	data := make(map[string]interface{})
	data["data"] = backends
	data["timestamp"] = time.Now().Unix()

	proxyStr, err := json.Marshal(data)
	if err != nil {
		golog.Error("proxy", zap.String("err", err.Error()))
		return err
	}

	err = cache.Set(cacheProxyName, string(proxyStr))
	if err != nil {
		golog.Error("proxy", zap.String("err", err.Error()))
		return err
	}

	return nil
}
