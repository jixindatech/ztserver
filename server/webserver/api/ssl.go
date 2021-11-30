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

type sslForm struct {
	Name string `json:"name" valid:"Required;MaxSize(254)"`
	Pub    string `json:"pub" valid:"Required;MaxSize(5120)"`
	Pri    string `json:"pri"  valid:"Required;MaxSize(5120)"`
	Remark string `json:"remark" valid:"MaxSize(254)"`
}

func AddSSL(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     sslForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	ssl := service.SSL{
		Name:   form.Name,
		Pub:    strings.TrimSpace(form.Pub),
		Pri:    strings.TrimSpace(form.Pri),
		Remark: form.Remark,
	}
	err := ssl.Save()
	if err != nil {
		golog.Error("ssl", zap.String("add", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.AddSSLFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetSSL(c *gin.Context) {
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

	ssl := service.SSL{
		ID: id,
	}
	idSSL, err := ssl.Get()
	if err != nil {
		golog.Error("user", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["record"] = idSSL
	appG.Response(httpCode, errCode, data)
}

type querySSLForm struct {
	Name     string `form:"name" valid:"MaxSize(254)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetSSLs(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     querySSLForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	ssl := service.SSL{
		Name:     form.Name,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	ssls, count, err := ssl.GetList()
	if err != nil {
		golog.Error("ssl", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetSSLFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = ssls
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

func PutSSL(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     sslForm
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

	ssl := service.SSL{
		ID:     id,
		Name:   form.Name,
		Pub:    form.Pub,
		Pri:    form.Pri,
		Remark: form.Remark,
	}
	err = ssl.Save()
	if err != nil {
		golog.Error("ssl", zap.String("put", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutSSLFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func DeleteSSL(c *gin.Context) {
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

	ssl := service.SSL{
		ID: id,
	}
	err = ssl.Delete()
	if err != nil {
		golog.Error("ssl", zap.String("delete", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.DeleteSSLFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}
