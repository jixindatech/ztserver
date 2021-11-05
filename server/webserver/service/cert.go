package service

import (
	"encoding/json"
	"go.uber.org/zap"
	"time"
	"zt-server/cache"
	"zt-server/pkg/core/golog"
	"zt-server/webserver/model"
)

const cacheSSLName string = "cert"

type Cert struct {
	ID uint64

	Server string
	Pub    string
	Pri    string
	Remark string

	Page     int
	PageSize int
}

func (c *Cert) Save() (err error) {
	data := make(map[string]interface{})
	data["server"] = c.Server
	data["pub"] = c.Pub
	data["pri"] = c.Pri
	data["remark"] = c.Remark

	if c.ID > 0 {
		err = model.PutCert(c.ID, data)
	} else {
		err = model.AddCert(data)
	}
	if err != nil {
		return err
	}

	return SetupCerts()
}

func (c *Cert) Get() (*model.Cert, error) {
	return model.GetCert(c.ID)
}

func (c *Cert) GetList() ([]*model.Cert, int, error) {
	data := make(map[string]interface{})
	data["server"] = c.Server
	data["page"] = c.Page
	data["pagesize"] = c.PageSize

	return model.GetCerts(data)
}

func (c *Cert) Delete() error {
	err := model.DeleteCert(c.ID)
	if err != nil {
		return err
	}
	return SetupCerts()
}

func SetupCerts() error {
	cert := Cert{
		Page: 0,
	}
	certs, count, err := cert.GetList()
	if err != nil {
		return err
	}
	if count == 0 {
		serverData := make(map[string]interface{})
		serverData["timestamp"] = time.Now().Unix()
		serverData["data"] = struct{}{}

		certData, err := json.Marshal(serverData)
		if err != nil {
			golog.Error("cert", zap.String("err", err.Error()))
			return err
		}

		err = cache.Set(cacheSSLName, string(certData))
		if err != nil {
			golog.Error("cert", zap.String("err", err.Error()))
			return err
		}

		return nil
	}

	servers := make(map[string]interface{})
	for _, item := range certs {
		certs := make(map[string]interface{})
		certs["pub"] = item.Pub
		certs["pri"] = item.Pri
		servers[item.Server] = certs
	}
	serverData := make(map[string]interface{})
	serverData["timestamp"] = time.Now().Unix()
	serverData["data"] = servers

	certData, err := json.Marshal(serverData)
	if err != nil {
		golog.Error("cert", zap.String("err", err.Error()))
		return err
	}

	err = cache.Set(cacheSSLName, string(certData))
	if err != nil {
		golog.Error("cert", zap.String("err", err.Error()))
		return err
	}

	return nil
}
