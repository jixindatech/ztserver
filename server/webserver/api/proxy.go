package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"zt-server/pkg/app"
	"zt-server/pkg/core/golog"
	"zt-server/pkg/e"
	"zt-server/webserver/service"
)

type proxies struct {
	IP   string `json:"ip" valid:"Required;IP;MaxSize(254)"`
	Port int    `json:"port" valid:"Required;Min(1);Max(65535)"`
}
type proxyForm struct {
	Server  string    `json:"server" valid:"Required;MaxSize(254)"`
	Lb      string    `json:"lb" valid:"Required;"`
	Backend []proxies `json:"backend" valid:"Required"`
	Remark  string    `json:"remark" valid:"MaxSize(254)"`
}

func AddProxy(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     proxyForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	backends, err := json.Marshal(form.Backend)
	if err != nil {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	proxy := service.Proxy{
		Server:  form.Server,
		Lb:      form.Lb,
		Backend: []byte(backends),
		Remark:  form.Remark,
	}

	err = proxy.Save()
	if err != nil {
		golog.Error("proxy", zap.String("add", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.AddProxyFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetProxy(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if id == 0 || err != nil {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	proxy := service.Proxy{
		ID: id,
	}
	idProxy, err := proxy.Get()
	if err != nil {
		golog.Error("proxy", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetProxyFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["record"] = idProxy
	appG.Response(httpCode, errCode, data)
}

type queryProxyForm struct {
	Server   string `form:"server" valid:"MaxSize(254)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetProxys(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryProxyForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	proxy := service.Proxy{
		Server:   form.Server,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	proxys, count, err := proxy.GetList()
	if err != nil {
		golog.Error("proxy", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetProxyFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = proxys
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

func PutProxy(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     proxyForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if id == 0 || err != nil {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	backends, err := json.Marshal(form.Backend)
	if err != nil {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	proxy := service.Proxy{
		ID:      id,
		Server:  form.Server,
		Lb:      form.Lb,
		Backend: backends,
		Remark:  form.Remark,
	}
	err = proxy.Save()
	if err != nil {
		golog.Error("proxy", zap.String("put", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutProxyFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func DeleteProxy(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if id == 0 || err != nil {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	proxy := service.Proxy{
		ID: id,
	}
	err = proxy.Delete()
	if err != nil {
		golog.Error("proxy", zap.String("delete", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.DeleteProxyFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}
