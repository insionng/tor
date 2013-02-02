package tor

import (
	"os"
)

var (
	app          *torApp
	util         torUtil
	cfg          *torConfig
	cfgFile      string = "app.conf"
	ListenAddr   string = ""
	ListenPort   int    = 80
	RunMode      string = "http"
	EnableDaemon bool   = false
	EnableStats  bool   = true
	CookieSecret string = "foobar"
	SessionName  string = "torSESSID"
	SessionTTL   int64  = 60 * 15
	EnablePprof  bool   = true
	EnableGzip   bool   = true
)

func init() {
	// Check the first argument of cmd line,
	// if it is not a flag (begin with '-'),
	// try to use it as the config file path.
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg[0] != '-' {
			cfgFile = arg
		}
	}

	cfg = &torConfig{}
	LoadConfig()

	app = NewApp()
	util = torUtil{}
}

func NewApp() *torApp {
	return new(torApp).init()
}

func GetMainApp() *torApp {
	return app
}

func RegisterController(pattern string, c torControllerInterface) {
	app.RegisterController(pattern, c)
}

func RegisterControllerHook(event string, hookFunc HookControllerFunc) {
	app.RegisterControllerHook(event, hookFunc)
}

func SetStaticPath(sPath, fPath string) {
	app.SetStaticPath(sPath, fPath)
}

func RegisterSessionStorage(storage SessionStorageInterface) {
	app.RegisterSessionStorage(storage)
}

func RegisterCustomHttpStatus(code int, filePath string) {
	app.RegisterCustomHttpStatus(code, filePath)
}

func Run() {
	if EnableDaemon {
		util.CallMethod(&util, "SetDaemonMode", 1, 0)
	}
	app.Run(RunMode, ListenAddr, ListenPort)
}

func LoadConfig() {
	err := cfg.LoadConfig("app.conf")
	if err != nil {
		return
	}
	if v, ok := cfg.GetConfig("ListenAddr").String(); ok {
		ListenAddr = v
	}
	if v, ok := cfg.GetConfig("ListenPort").Int(); ok {
		ListenPort = v
	}
	if v, ok := cfg.GetConfig("RunMode").String(); ok {
		RunMode = v
	}
	if v, ok := cfg.GetConfig("EnableDaemon").Bool(); ok {
		EnableDaemon = v
	}
	if v, ok := cfg.GetConfig("EnableStats").Bool(); ok {
		EnableStats = v
	}
	if v, ok := cfg.GetConfig("SessionName").String(); ok {
		SessionName = v
	}
	if v, ok := cfg.GetConfig("SessionTTL").Int(); ok {
		SessionTTL = int64(v)
	}
	if v, ok := cfg.GetConfig("EnablePprof").Bool(); ok {
		EnablePprof = v
	}
}

func GetConfig(key string) *torConfigValue {
	return cfg.GetConfig(key)
}
