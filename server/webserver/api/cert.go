package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"zt-server/pkg/app"
	"zt-server/pkg/core/golog"
	"zt-server/pkg/e"
	"zt-server/webserver/service"
)

type certForm struct {
	Server string `json:"server" valid:"Required;MaxSize(254)"`
	Pub    string `json:"pub" valid:"Required;MaxSize(5120)"`
	Pri    string `json:"pri"  valid:"Required;MaxSize(5120)"`
	Remark string `json:"remark" valid:"MaxSize(254)"`
}

func AddCert(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     certForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	cert := service.Cert{
		Server: form.Server,
		Pub:    strings.TrimSpace(form.Pub),
		Pri:    strings.TrimSpace(form.Pri),
		Remark: form.Remark,
	}
	err := cert.Save()
	if err != nil {
		golog.Error("cert", zap.String("add", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.AddCertFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetCert(c *gin.Context) {
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

	cert := service.Cert{
		ID: id,
	}
	idCert, err := cert.Get()
	if err != nil {
		golog.Error("user", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["record"] = idCert
	appG.Response(httpCode, errCode, data)
}

type queryCertForm struct {
	Server   string `form:"server" valid:"MaxSize(254)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetCerts(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryCertForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	cert := service.Cert{
		Server:   form.Server,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	certs, count, err := cert.GetList()
	if err != nil {
		golog.Error("cert", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetCertFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = certs
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

func PutCert(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     certForm
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

	cert := service.Cert{
		ID:     id,
		Server: form.Server,
		Pub:    form.Pub,
		Pri:    form.Pri,
		Remark: form.Remark,
	}
	err = cert.Save()
	if err != nil {
		golog.Error("cert", zap.String("put", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutCertFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func DeleteCert(c *gin.Context) {
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

	cert := service.Cert{
		ID: id,
	}
	err = cert.Delete()
	if err != nil {
		golog.Error("cert", zap.String("delete", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.DeleteCertFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}
