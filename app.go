package tor

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
)

type torApp struct {
	router           *torRouter
	hook             *torHook
	session          *torSessionManager
	customHttpStatus map[int]string
	// extHook *torHook
}

func (this *torApp) init() *torApp {
	this.router = &torRouter{
		app:         this,
		Rules:       []*torRoutingRule{},
		StaticRules: []*torRoutingRule{},
		StaticDir:   make(map[string]string),
	}
	this.hook = &torHook{app: this}
	// this.extHook = &torHook{app: this}
	this.session = new(torSessionManager)
	this.session.RegisterStorage(new(torDefaultSessionStorage))
	this.customHttpStatus = make(map[int]string)
	return this
}

func (this *torApp) RegisterController(pattern string, c torControllerInterface) {
	this.router.AddRule(pattern, c)
}

func (this *torApp) RegisterControllerHook(event string, hookFunc HookControllerFunc) {
	this.hook.AddControllerHook(event, hookFunc)
}

func (this *torApp) callControllerHook(event string, hc *HookController) {
	this.hook.CallControllerHook(event, hc)
	// this.extHook.CallControllerHook(event, hc)
}

// func (this *torApp) registerAddonControllerHook(event string, hookFunc HookControllerFunc) {
// 	this.extHook.AddControllerHook(event, hookFunc)
// }

// func (this *torApp) clearExtHook(event string, hookFunc HookControllerFunc) {
// 	this.extHook = &torHook{app: this}
// }

func (this *torApp) SetStaticPath(sPath, fPath string) {
	this.router.SetStaticPath(sPath, fPath)
}

func (this *torApp) RegisterSessionStorage(storage SessionStorageInterface) {
	this.session.RegisterStorage(storage)
}

func (this *torApp) RegisterCustomHttpStatus(code int, filePath string) {
	this.customHttpStatus[code] = filePath
}

func (this *torApp) Run(mode string, addr string, port int) {
	listenAddr := net.JoinHostPort(addr, fmt.Sprintf("%d", port))
	var err error
	switch mode {
	case "http":
		err = http.ListenAndServe(listenAddr, this.router)
	case "fcgi":
		l, e := net.Listen("tcp", listenAddr)
		if e != nil {
			panic("Fcgi listen error: " + e.Error())
		}
		err = fcgi.Serve(l, this.router)
	default:
		err = http.ListenAndServe(listenAddr, this.router)
	}
	if err != nil {
		panic("ListenAndServe error: " + err.Error())
	}
}

func (this *torApp) AppPath() string {
	path, _ := os.Getwd()
	return path
}
