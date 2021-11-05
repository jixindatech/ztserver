package webserver

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"zt-server/webserver/api"
)

func registerRouter(r *gin.Engine, jwt *jwt.GinJWTMiddleware) {
	r.GET("/ws", WsHandler)

	r.NoRoute(jwt.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/admin/login", jwt.LoginHandler)

	adminRouter := r.Group("/admin")
	adminRouter.Use(jwt.MiddlewareFunc())
	adminRouter.GET("/info", getInfo)
	adminRouter.POST("/logout", logout)

	routers := r.Group("/api/v1")
	routers.POST("/user/", api.AddUser)
	routers.GET("/user/:id", api.GetUser)
	routers.GET("/user/", api.GetUsers)
	routers.PUT("/user/:id", api.PutUser)
	routers.DELETE("/user/:id", api.DeleteUser)
	routers.POST("/user/:id/resource/", api.SaveUserResource)
	routers.GET("/user/:id/resource/", api.GetUserResource)
	routers.POST("/user/:id/email/", api.SendUserMail)

	routers.POST("/resource/", api.AddResource)
	routers.GET("/resource/:id", api.GetResource)
	routers.GET("/resource/", api.GetResources)
	routers.PUT("/resource/:id", api.PutResource)
	routers.DELETE("/resource/:id", api.DeleteResource)

	routers.POST("/email/", api.AddEmail)
	routers.GET("/email/", api.GetEmail)
	routers.PUT("/email/:id", api.PutEmail)

	routers.POST("/cert/", api.AddCert)
	routers.GET("/cert/:id", api.GetCert)
	routers.GET("/cert/", api.GetCerts)
	routers.PUT("/cert/:id", api.PutCert)
	routers.DELETE("/cert/:id", api.DeleteCert)

	routers.POST("/proxy/", api.AddProxy)
	routers.GET("/proxy/:id", api.GetProxy)
	routers.GET("/proxy/", api.GetProxys)
	routers.PUT("/proxy/:id", api.PutProxy)
	routers.DELETE("/proxy/:id", api.DeleteProxy)

	routers.GET("/event/gw/", api.GetGwEvents)
	routers.GET("/event/ws/", api.GetWsEvents)

	routers.GET("/info/user/", api.GetUserCount)
	routers.GET("/info/resource/", api.GetResourceCount)
	routers.GET("/info/user/online/", api.GetUserOnline)

	routers.GET("/info/user/gw/", api.GetUserGwStatics)
}
