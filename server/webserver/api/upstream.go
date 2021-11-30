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

type upstreamForm struct {
	Name  string    `json:"name" valid:"Required;MaxSize(254)"`
	Lb      string  `json:"lb" valid:"Required;Match(/chash|roundrobin/)"`
	Key     string  `json:"key"`
	Backend []service.Node `json:"backend" valid:"Required"`
	Retry          int `json:"retry" valid:"Required;Min(1)"`
	TimeoutConnect int `json:"timeoutConnect" valid:"Required;Min(1)"`
	TimeoutSend    int `json:"timeoutSend" valid:"Required;Min(1)"`
	TimeoutReceive int `json:"timeoutReceive" valid:"Required;Min(1)"`

	Remark         string    `json:"remark" valid:"MaxSize(254)"`
}

func AddUpstream(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     upstreamForm
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

	if form.Lb == "chash" && len(form.Key) == 0{
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	if len(form.Key) > 0 {
		if form.Key != "remote-addr" {
			httpCode = e.InvalidParams
			appG.Response(httpCode, errCode, nil)
			return
		}
	}

	upstream := service.Upstream{
		Name:  form.Name,
		Lb:      form.Lb,
		Key: form.Key,
		Backend: []byte(backends),
		Retry: form.Retry,
		TimeoutConnect: form.TimeoutConnect,
		TimeoutSend: form.TimeoutSend,
		TimeoutReceive: form.TimeoutReceive,

		Remark:  form.Remark,
	}

	err = upstream.Save()
	if err != nil {
		golog.Error("upstream", zap.String("add", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.AddUpstreamFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetUpstream(c *gin.Context) {
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

	upstream := service.Upstream{
		ID: id,
	}
	idUpstream, err := upstream.Get()
	if err != nil {
		golog.Error("upstream", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUpstreamFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["record"] = idUpstream
	appG.Response(httpCode, errCode, data)
}

type queryUpstreamForm struct {
	Name   string `form:"name" valid:"MaxSize(254)"`
	Page     int    `form:"current" valid:"Min(0)"`
	PageSize int    `form:"size" valid:"Min(0);Max(50)"`
}

func GetUpstreams(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryUpstreamForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	upstream := service.Upstream{
		Name:    form.Name,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	upstreams, count, err := upstream.GetList()
	if err != nil {
		golog.Error("upstream", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUpstreamFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = upstreams
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

func PutUpstream(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     upstreamForm
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

	upstream := service.Upstream{
		ID:      id,
		Name:    form.Name,
		Lb:      form.Lb,
		Backend: backends,
		Retry: form.Retry,
		TimeoutConnect: form.TimeoutConnect,
		TimeoutSend: form.TimeoutSend,
		TimeoutReceive: form.TimeoutReceive,

		Remark:  form.Remark,
	}
	err = upstream.Save()
	if err != nil {
		golog.Error("upstream", zap.String("put", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutUpstreamFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func DeleteUpstream(c *gin.Context) {
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

	upstream := service.Upstream{
		ID: id,
	}
	err = upstream.Delete()
	if err != nil {
		golog.Error("upstream", zap.String("delete", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.DeleteUpstreamFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}
