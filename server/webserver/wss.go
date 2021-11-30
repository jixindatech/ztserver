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
	"zt-server/webserver/model"
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

func getUserResoruce(data string, user *model.User) (map[string]int, error) {
	/* id, timestamp, email */
	userJs, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return nil, err
	}

	email, err := userJs.Get("email").String()
	if err != nil {
		return nil, err
	}

	if email != user.Email {
		return nil, fmt.Errorf("%s", "invalid user email")
	}

	timestamp, err := userJs.Get("timestamp").Int64()
	if err != nil {
		return nil, err
	}
	if timestamp != user.Timestamp {
		return nil, fmt.Errorf("%s", "invalid user timestamp")
	}

	id, err := userJs.Get("id").Uint64()
	if err != nil {
		return nil, err
	}

	if id != user.ID {
		return nil, fmt.Errorf("%s", "invalid user id")
	}

	users, err := cache.Get(service.GW)
	if err != nil {
		return nil, err
	}

	var gwCache cache.GwCache
	err = json.Unmarshal(users.([]byte), &gwCache)
	if err != nil {
		return nil, err
	}
	/*
	userBytes := users.([]uint8)
	js, err := simplejson.NewJson(userBytes)
	if err != nil {
		return nil, err
	}

	userData := js.Get("values").Interface()
	usersCache := userData.([]cache.UserCache)
	*/
	var userResouces *cache.UserCache
	for _, userCache := range gwCache.Values {
		if userCache.ID == id {
			userResouces = userCache
			break
		}
	}

	res := make(map[string]int)
	if userResouces != nil  {
		for _, resouce := range userResouces.Resources {
			res[resouce.Host] = 1
		}
	} else {
		return nil, fmt.Errorf("%s", "invalid user")
	}

	return res, nil
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

			var data string
			data, err = js.Get("data").String()
			if err != nil {
				golog.Error("ws", zap.String("data 111", err.Error()))
				break
			}

			id, err := js.Get("id").Uint64()
			if err != nil {
				golog.Error("ws", zap.String("id", err.Error()))
				break
			}

			userSrv := service.User{
				ID:       id,
			}
			user, err := userSrv.Get()
			if err != nil {
				golog.Error("ws", zap.String("user", err.Error()))
				break
			}

			var info string
			info, err = utils.CBCDecryptWithPKCS7(user.Secret, data)
			if err != nil {
				golog.Error("ws", zap.String("err", err.Error()))
				break
			}

			resources, err = getUserResoruce(info, user)
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
				wsService.User = user.Name
				wsService.Email = user.Email
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
