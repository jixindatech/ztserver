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

type queryWsForm struct {
	User     string `form:"user" valid:"MaxSize(254)"`
	Resource string `form:"resource" valid:"MaxSize(254)"`
	Start    int64  `form:"start" valid:"Required;Min(1)"`
	End      int64  `form:"end" valid:"Required;Min(1)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetUserOnline(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	count := service.GetUserOnopenCount()
	data := make(map[string]interface{})
	data["count"] = count
	appG.Response(httpCode, errCode, data)
}

func GetWsEvents(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryWsForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	ws := service.Ws{
		User:     form.User,
		Resource: form.Resource,
		Start:    form.Start,
		End:      form.End,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	wsEvents, count, err := ws.GetList()
	if err != nil {
		golog.Error("ws", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetWsEventsFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = wsEvents
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}
