package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"zt-server/pkg/app"
	"zt-server/pkg/core/golog"
	"zt-server/pkg/e"
	"zt-server/webserver/service"
)

type queryGwForm struct {
	User     string `form:"user" valid:"MaxSize(254)"`
	Method   string `form:"method" valid:"MaxSize(16)"`
	Resource string `form:"resource" valid:"MaxSize(254)"`
	Start    int64  `form:"start" valid:"Required;Min(1)"`
	End      int64  `form:"end" valid:"Required;Min(1)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetGwEvents(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryGwForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	gw := service.Gw{
		User:     form.User,
		Method:   form.Method,
		Resource: form.Resource,
		Start:    form.Start,
		End:      form.End,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	gwEvents, count, err := gw.GetList()
	if err != nil {
		golog.Error("gw", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetGwEventsFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = gwEvents
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

type queryUserGwForm struct {
	Start int64 `form:"start" valid:"Required;Min(1)"`
	End   int64 `form:"end" valid:"Required;Min(1)"`
}

func GetUserGwStatics(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		form     queryUserGwForm
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	gw := service.Gw{
		Start: form.Start,
		End:   form.End,
	}

	res, err := gw.GetStatics()
	if err != nil {
		golog.Error("gw", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetGwEventsFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["statics"] = res

	appG.Response(httpCode, errCode, data)
}
