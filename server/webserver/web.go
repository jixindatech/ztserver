package webserver

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"net/http"
	"zt-server/webserver/model"
	"zt-server/webserver/service"
	// "github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"time"
)

type Server struct {
	Addr        string
	User        string
	Password    string
	JwtKey      []byte
	IdentityKey string
	Dbname      string

	g *gin.Engine
}

type User struct {
	UserName string
	UserId   uint
}

type userLoginForm struct {
	Username string `json:"username" form:"username" validate:"required,max=254"`
	Password string `json:"password" form:"password" validate:"required,max=254"`
}

func (web *Server) Init(key string) error {
	if len(web.Addr) == 0 {
		return fmt.Errorf("%s", "invalid webaddr")
	}
	if len(web.User) == 0 || len(web.Password) == 0 {
		return fmt.Errorf("%s", "user or password is empty")
	}
	if len(web.JwtKey) == 0 {
		return fmt.Errorf("%s", "jwtkey is empty")
	}
	if len(web.IdentityKey) == 0 {
		return fmt.Errorf("%s", "identitykey is empty")
	}

	var err error
	err = model.OpenDatabase(web.Dbname)
	if err != nil {
		return err
	}

	err = service.SetupSSLs()
	if err != nil {
		return err
	}

	err = service.SetupUpstreams()
	if err != nil {
		return err
	}

	err = service.SetupRouters()
	if err != nil {
		return err
	}

	err = service.SetupUserResource()
	if err != nil {
		return err
	}

	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "zero-trust",
		Key:         web.JwtKey,
		Timeout:     24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: web.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					web.IdentityKey: v.UserId,
					"username":      v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserId:   uint(claims[web.IdentityKey].(float64)),
				UserName: claims["username"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login userLoginForm
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := login.Username
			password := login.Password

			if username == web.User && password == web.Password {
				return &User{
					UserName: username,
					UserId:   1,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == web.User {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, cookie: Admin-Token",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	err = authMiddleware.MiddlewareInit()
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	registerRouter(r, authMiddleware)

	r.StaticFS("/static", http.Dir("dashboard/dist/static"))
	r.StaticFile("/", "dashboard/dist/index.html")
	r.StaticFile("/favicon.ico", "dashboard/dist/favicon.ico")

	web.g = r
	return nil
}

func (web *Server) Run() error {
	return web.g.Run(web.Addr)
}

func (web *Server) Close() {
	if web.g != nil {

	}
}
