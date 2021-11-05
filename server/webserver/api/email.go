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

type emailForm struct {
	Server   string `json:"server"   valid:"Required;MaxSize(254)"`
	Port     int    `json:"port"   valid:"Required;Min(1);Max(65535)"`
	Email    string `json:"email" valid:"Required;Email;MaxSize(254)"`
	Password string `json:"password"  valid:"Required;MaxSize(254)"`
	Remark   string `json:"remark" valid:"MaxSize(254)"`
}

func AddEmail(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     emailForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	user := service.Email{
		Server:   form.Server,
		Port:     form.Port,
		Email:    form.Email,
		Password: form.Password,

		Remark: form.Remark,
	}
	err := user.Save()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.AddEmailFailed
		golog.Error("email", zap.String("add", err.Error()))
	}

	appG.Response(httpCode, errCode, nil)
}

func GetEmail(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	email := service.Email{}
	idEmail, err := email.Get()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserFailed
		golog.Error("user", zap.String("get", err.Error()))
	}

	data := make(map[string]interface{})
	data["email"] = idEmail
	appG.Response(httpCode, errCode, data)
}

func PutEmail(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     emailForm
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

	email := service.Email{
		ID:       id,
		Server:   form.Server,
		Port:     form.Port,
		Email:    form.Email,
		Password: form.Password,

		Remark: form.Remark,
	}
	err = email.Save()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.PutEmailFailed
		golog.Error("email", zap.String("put", err.Error()))
	}

	appG.Response(httpCode, errCode, nil)
}
