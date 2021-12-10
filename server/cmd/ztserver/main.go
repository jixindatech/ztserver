package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"zt-server/cache"
	"zt-server/pkg/core/golog"
	"zt-server/settings"
	"zt-server/ssh"
	"zt-server/storage"
	"zt-server/webserver"
)

var configFile *string = flag.String("config", "./etc/config.yaml", "zt-server config file")
var version *bool = flag.Bool("v", false, "the version ")

var (
	BuildDate    string
	BuildVersion string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if *version {
		fmt.Printf("Git commit:%s\n", BuildVersion)
		fmt.Printf("Build time:%s\n", BuildDate)
		return
	}

	_ = golog.SetDefaultZapLog()

	var err error
	cfg, err := settings.ParseConfigFile(*configFile)
	if err != nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}

	if len(cfg.LogPath) != 0 {
		_ = golog.InitZapLog(cfg.LogLevel, cfg.LogPath)
	} else {
		_ = golog.InitZapLog(cfg.LogLevel, "stdout")
	}

	if cfg.WebAdmin == nil {
		golog.Fatal("main", zap.String("web", fmt.Errorf("%s", "web config missed").Error()))
	}

	err = ssh.Init(cfg.Ssh)
	if err != nil {
		golog.Fatal("main", zap.String("ssh", err.Error()))
	}

	err = cache.Init(cfg.Redis)
	if err != nil {
		golog.Fatal("main", zap.String("cache", err.Error()))
	}

	err = storage.InitStorage(cfg.Es)
	if err != nil {
		golog.Fatal("main", zap.String("es", err.Error()))
	}

	var web *webserver.Server
	web = &webserver.Server{
		Addr:        cfg.WebAdmin.Addr,
		User:        cfg.WebAdmin.User,
		Password:    cfg.WebAdmin.Password,
		JwtKey:      []byte(cfg.WebAdmin.JwtKey),
		IdentityKey: cfg.WebAdmin.IdentityKey,
		Dbname:      cfg.DbPath,
	}

	gin.SetMode(cfg.Mode)
	gin.DefaultWriter = golog.GetLogger()

	err = web.Init(cfg.Secret, cfg.Database)
	if err != nil {
		golog.Fatal("main", zap.String("web", err.Error()))
	}

	golog.Info("main", zap.String("web", "running"))
	go func() {
		err := web.Run()
		if err != nil {
			golog.Fatal("main", zap.String("web", err.Error()))
		}
	}()

	golog.Info("main", zap.String("status", "running"))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
		//syscall.SIGUSR1,
	)

	for {
		sig := <-sc
		if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
			golog.Info("main", zap.String("signal", fmt.Sprintf("%d", sig)))
			if web != nil {
				web.Close()
			}
			golog.Close()
			break
		} else if sig == syscall.SIGPIPE {
			//golog.Info("main", "main", "Ignore broken pipe signal", 0)
			//skip
		}
	}

	golog.Info("main", zap.String("status", "quit"))
}
