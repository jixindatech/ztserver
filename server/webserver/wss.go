package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
	"zt-server/cache"
	"zt-server/pkg/core/golog"
	"zt-server/storage"
	"zt-server/webserver/service"
	"zt-server/webserver/utils"
)

const (
	OPEN    = "OPEN"
	RUNNGIN = "RUNNGIN"
	CLOSE   = "CLOSE"
	WsIndex = "ztserver"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getUserResoruce(token string, userMail string) (map[string]int, error) {
	userJs, err := simplejson.NewJson([]byte(token))
	if err != nil {
		return nil, err
	}
	email, err := userJs.Get("email").String()
	if err != nil {
		return nil, err
	}
	if len(email) == 0 {
		return nil, fmt.Errorf("%s", "invalid token email")
	}

	if email != userMail {
		return nil, fmt.Errorf("%s", "invalid user email")
	}

	timestamp, err := userJs.Get("timestamp").Int64()
	if err != nil {
		return nil, err
	}
	if timestamp == 0 {
		return nil, fmt.Errorf("%s", "invalid token timestamp")
	}

	id, err := userJs.Get("id").Uint64()
	if err != nil {
		return nil, err
	}
	if timestamp == 0 {
		return nil, fmt.Errorf("%s", "invalid user id")
	}

	user, err := cache.Get(email)
	if err != nil {
		return nil, err
	}
	userBytes := user.([]uint8)
	var userCache cache.UserCache
	err = json.Unmarshal([]byte(userBytes), &userCache)
	if err != nil {
		return nil, err
	}
	if userCache.ID != id {
		return nil, fmt.Errorf("%s", "invalid id")
	}
	if userCache.Timestamp != timestamp {
		return nil, fmt.Errorf("%s", "invalid timestamp")
	}

	return userCache.Resource, nil
}

func WsHandler(c *gin.Context) {
	var err error
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		golog.Error("ws", zap.String("upgrade", err.Error()))
		return
	}
	defer ws.Close()

	header := c.GetHeader("x-forwarded-for")
	var ip []string
	if len(header) > 0{
		ip = strings.Split(header, ",")
		if (len(ip)) == 0{
			golog.Error("ws",
				zap.String("ip", fmt.Errorf("%s", "invalid x-forwarded-for header").Error()))
			return
		}
	} else {
		return
	}

	// ip := strings.Split(ws.RemoteAddr().String(), ":")
	wsService := service.WssEvent{
		Client: ip[0],
	}

	var resources map[string]int

	login := false
	for {
		var mt int
		var message []byte
		mt, message, err = ws.ReadMessage()
		if err != nil {
			break
		}

		if len(message) > 0 {
			var js *simplejson.Json
			js, err = simplejson.NewJson(message)
			if err != nil {
				golog.Error("ws", zap.String("json", err.Error()))
				break
			}

			var token string
			token, err = js.Get("token").String()
			if err != nil {
				golog.Error("ws", zap.String("token", err.Error()))
				break
			}

			var email string
			email, err = js.Get("email").String()
			if err != nil {
				golog.Error("ws", zap.String("email", err.Error()))
				break
			}

			var info string
			info, err = utils.CBCDecryptWithPKCS7(token)
			if err != nil {
				golog.Error("ws", zap.String("err", err.Error()))
				break
			}

			resources, err = getUserResoruce(info, email)
			if err != nil {
				golog.Error("ws", zap.String("err", err.Error()))
				break
			}

			if !login {
				var dev []byte
				dev, err = js.Get("dev").MarshalJSON()
				if err != nil {
					golog.Error("ws", zap.String("err", err.Error()))
					break
				}

				login = true
				wsService.Event = OPEN
				wsService.User = email
				wsService.Email = email
				wsService.Dev = string(dev)
				wsService.Time = time.Now().Unix()

				var ips []string
				ips, err = wsService.OpenDoor(resources)
				if err != nil {
					golog.Error("ws", zap.String("opendoor", err.Error()))
					break
				}

				wsService.Gw = ips
				err = wsService.Add(storage.WS_INDEX)
				if err != nil {
					golog.Error("ws", zap.String("database", err.Error()))
				}
			}
		}

		time.Sleep(time.Duration(5) * time.Second)

		err = ws.WriteMessage(mt, []byte("ok"))
		if err != nil {
			golog.Error("ws", zap.String("send", err.Error()))
			break
		}
	}

	if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		golog.Error("ws", zap.String("break", err.Error()))
	}

	if login {
		_, err = wsService.CloseDoor(resources)
		if err != nil {
			golog.Error("ws", zap.String("closedoor", err.Error()))
		}

		wsService.Event = CLOSE
		wsService.Time = time.Now().Unix()
		err = wsService.Add(storage.WS_INDEX)
		if err != nil {
			golog.Error("ws", zap.String("index", err.Error()))
		}
	}
}
