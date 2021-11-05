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

type userForm struct {
	Name   string `json:"name"   valid:"Required;MaxSize(254)"`
	Email  string `json:"email" valid:"Required;Email;MaxSize(254)"`
	Phone  string `json:"phone"  valid:"Required;Phone"`
	Status int    `json:"status" valid:"Range(0,1)"`
	Remark string `json:"remark" valid:"MaxSize(254)"`
}

func AddUser(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     userForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	user := service.User{
		Name:   form.Name,
		Email:  form.Email,
		Phone:  form.Phone,
		Status: form.Status,
		Remark: form.Remark,
	}
	err := user.Save()
	if err != nil {
		golog.Error("user", zap.String("add", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.AddUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetUser(c *gin.Context) {
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

	user := service.User{
		ID: id,
	}
	idUser, err := user.Get()
	if err != nil {
		golog.Error("user", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["user"] = idUser
	appG.Response(httpCode, errCode, data)
}

func GetUserCount(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	user := service.User{
		Page:     1,
		PageSize: 1,
	}
	_, count, err := user.GetList()
	if err != nil {
		golog.Error("user", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["count"] = count
	appG.Response(httpCode, errCode, data)
}

type queryUserForm struct {
	Name     string `form:"username" valid:"MaxSize(254)"`
	Page     int    `form:"current" valid:"Required;Range(1,50)"`
	PageSize int    `form:"size" valid:"Required;Min(1)"`
}

func GetUsers(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryUserForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = e.InvalidParams
		appG.Response(httpCode, errCode, nil)
		return
	}

	user := service.User{
		Name:     form.Name,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	users, count, err := user.GetList()
	if err != nil {
		golog.Error("user", zap.String("get", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = users
	data["total"] = count
	appG.Response(httpCode, errCode, data)
}

func PutUser(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     userForm
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

	user := service.User{
		ID:     id,
		Name:   form.Name,
		Email:  form.Email,
		Phone:  form.Phone,
		Status: form.Status,
		Remark: form.Remark,
	}
	err = user.Save()
	if err != nil {
		golog.Error("user", zap.String("put", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.PutUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func DeleteUser(c *gin.Context) {
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

	user := service.User{
		ID: id,
	}
	err = user.Delete()
	if err != nil {
		golog.Error("user", zap.String("delete", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.DeleteUserFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

type userResourceForm struct {
	Ids []uint64 `json:"ids" valid:"Required;"`
}

func SaveUserResource(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     userResourceForm
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

	user := service.User{
		ID:          id,
		ResourceIds: form.Ids,
	}

	err = user.SaveUserResource()
	if err != nil {
		golog.Error("user", zap.String("resource", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.SaveUserResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	appG.Response(httpCode, errCode, nil)
}

func GetUserResource(c *gin.Context) {
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

	user := service.User{
		ID: id,
	}

	res, err := user.GetUserResource()
	if err != nil {
		golog.Error("user", zap.String("resource", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = res

	appG.Response(httpCode, errCode, data)
}

func SendUserMail(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     userForm
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

	user := service.User{
		ID:     id,
		Name:   form.Name,
		Email:  form.Email,
		Phone:  form.Phone,
		Status: form.Status,
		Remark: form.Remark,
	}
	err = user.SendMail()
	if err != nil {
		golog.Error("mail", zap.String("send", err.Error()))
		httpCode = http.StatusInternalServerError
		errCode = e.GetUserResourceFailed
		appG.Response(httpCode, errCode, nil)
		return
	}
	appG.Response(httpCode, errCode, nil)
}
