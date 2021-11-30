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

	routers.POST("/ssl/", api.AddSSL)
	routers.GET("/ssl/:id", api.GetSSL)
	routers.GET("/ssl/", api.GetSSLs)
	routers.DELETE("/ssl/:id", api.DeleteSSL)

	routers.POST("/upstream/", api.AddUpstream)
	routers.GET("/upstream/:id", api.GetUpstream)
	routers.GET("/upstream/", api.GetUpstreams)
	routers.PUT("/upstream/:id", api.PutUpstream)
	routers.DELETE("/upstream/:id", api.DeleteUpstream)

	routers.POST("/router/", api.AddRouter)
	routers.GET("/router/:id", api.GetRouter)
	routers.GET("/router/", api.GetRouters)
	routers.PUT("/router/:id", api.PutRouter)
	routers.DELETE("/router/:id", api.DeleteRouter)

	routers.GET("/event/gw/", api.GetGwEvents)
	routers.GET("/event/ws/", api.GetWsEvents)

	routers.GET("/info/user/", api.GetUserCount)
	routers.GET("/info/resource/", api.GetResourceCount)
	routers.GET("/info/user/online/", api.GetUserOnline)

	routers.GET("/info/user/gw/", api.GetUserGwStatics)
}
