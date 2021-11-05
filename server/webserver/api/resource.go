package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"zt-server/pkg/app"
	"zt-server/pkg/core/golog"
	"zt-server/pkg/e"
	"zt-server/webserver/service"
)

type resourceForm struct {
	Name   string `form:"name" valid:"Required;MaxSize(254)"`
	Server string `json:"server" valid:"Required;MaxSize(254)"`
	Remark string `json:"remark" valid:"MaxSize(254)"`
}

func AddResource(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     resourceForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	resource := service.Resource{
		Name:   form.Name,
		Server: form.Server,
		Remark: form.Remark,
	}
	err := resource.Save()
	if err != nil {
		golog.Error("resource", zap.String("add", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.AddResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetResourceCount(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	resource := service.Resource{
		Page:     1,
		PageSize: 1,
	}
	_, count, err := resource.GetList()
	if err != nil {
		golog.Error("resource", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["count"] = count
	appG.Response(httpCode, errCode, data)
}

type queryResourceForm struct {
	Name     string `form:"name" valid:"MaxSize(254)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetResources(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryResourceForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	resource := service.Resource{
		Name:     form.Name,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	resources, count, err := resource.GetList()
	if err != nil {
		golog.Error("resource", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = resources
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

func GetResource(c *gin.Context) {
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

	resource := service.Resource{
		ID: id,
	}
	idResource, err := resource.Get()
	if err != nil {
		golog.Error("resource", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["resource"] = idResource
	appG.Response(httpCode, errCode, data)
}

func PutResource(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     resourceForm
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

	resource := service.Resource{
		ID:     id,
		Name:   form.Name,
		Server: form.Server,
		Remark: form.Remark,
	}
	err = resource.Save()
	if err != nil {
		golog.Error("resource", zap.String("put", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func DeleteResource(c *gin.Context) {
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

	resource := service.Resource{
		ID: id,
	}
	err = resource.Delete()
	if err != nil {
		golog.Error("resource", zap.String("delete", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}
