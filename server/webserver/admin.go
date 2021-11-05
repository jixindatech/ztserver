package webserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zt-server/pkg/app"
	"zt-server/pkg/e"
)

func getInfo(c *gin.Context) {
	appG := app.Gin{C: c}

	data := make(map[string]interface{})
	data["introduction"] = "Administrator"
	data["name"] = "admin"
	data["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	data["roles"] = []string{"admin"}

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func logout(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
