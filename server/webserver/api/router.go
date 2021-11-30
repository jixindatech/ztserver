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

type routerForm struct {
	Name  string      `json:"name" valid:"Required;MaxSize(254)"`
	Host      string  `json:"host" valid:"Required;"`
	Path      string  `json:"path" valid:"Required;"`
	UptreamRef uint64 `json:"upstreamRef" valid:"Required"`
	Remark  string    `json:"remark" valid:"MaxSize(254)"`
}

func AddRouter(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     routerForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	router := service.Router{
		Name:  form.Name,
		Host:  form.Host,
		Path:  form.Path,
		UpstreamRef: form.UptreamRef,

		Remark:  form.Remark,
	}

	err := router.Save()
	if err != nil {
		golog.Error("router", zap.String("add", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.AddRouterFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetRouter(c *gin.Context) {
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

	router := service.Router{
		ID: id,
	}
	idRouter, err := router.Get()
	if err != nil {
		golog.Error("router", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetRouterFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["record"] = idRouter
	appG.Response(httpCode, errCode, data)
}

type queryRouterForm struct {
	Name   string `form:"name" valid:"MaxSize(254)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetRouters(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryRouterForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	router := service.Router{
		Name:   form.Name,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	routers, count, err := router.GetList()
	if err != nil {
		golog.Error("router", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetRouterFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = routers
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

func PutRouter(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     routerForm
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

	router := service.Router{
		ID:      id,
		Host:  form.Host,
		Path:  form.Path,
		UpstreamRef: form.UptreamRef,

		Remark:  form.Remark,
	}
	err = router.Save()
	if err != nil {
		golog.Error("router", zap.String("put", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutRouterFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func DeleteRouter(c *gin.Context) {
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

	router := service.Router{
		ID: id,
	}
	err = router.Delete()
	if err != nil {
		golog.Error("router", zap.String("delete", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.DeleteRouterFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}
