package tor

import (
	"fmt"
	"time"
)

type SessionStorageInterface interface {
	Init(int64)
	CreateSessionID() string
	Set(string, map[string]string)
	Get(string) map[string]string
	Delete(string)
}

type torSessionManager struct {
	sessionStorage SessionStorageInterface
	inited         bool
}

func (this *torSessionManager) RegisterStorage(storage SessionStorageInterface) {
	if storage == nil {
		return
	}
	this.sessionStorage = storage
	this.inited = false
}

func (this *torSessionManager) checkInit() {
	if !this.inited {
		this.sessionStorage.Init(SessionTTL)
		this.inited = true
	}
}

func (this *torSessionManager) CreateSessionID() string {
	this.checkInit()
	return this.sessionStorage.CreateSessionID()
}

func (this *torSessionManager) Set(sid string, data map[string]string) {
	this.checkInit()
	this.sessionStorage.Set(sid, data)
}

func (this *torSessionManager) Get(sid string) map[string]string {
	this.checkInit()
	return this.sessionStorage.Get(sid)
}

func (this *torSessionManager) Delete(sid string) {
	this.checkInit()
	this.sessionStorage.Delete(sid)
}

type torSession struct {
	ctlr           *Controller
	sessionManager *torSessionManager
	sessionId      string
	ctx            *torContext
	data           map[string]string
}

func (this *torSession) init() {
	if this.sessionId == "" {
		this.sessionId = this.sessionManager.CreateSessionID()
		this.ctx.SetSecureCookie(SessionName, this.sessionId, 0)
	}
	if this.data == nil {
		this.data = this.sessionManager.Get(this.sessionId)
	}
}

func (this *torSession) Get(key string) string {
	this.init()
	if data, exist := this.data[key]; exist {
		return data
	}
	return ""
}

func (this *torSession) Set(key string, data string) {
	this.init()
	this.data[key] = data
	this.sessionManager.Set(this.sessionId, this.data)
}

func (this *torSession) Delete(key string) {
	this.init()
	delete(this.data, key)
	this.sessionManager.Set(this.sessionId, this.data)
}

type torDefaultSessionStorage struct {
	ttl   int64
	datas map[string]torDefaultSessionStorageData
}

type torDefaultSessionStorageData struct {
	expires int64
	data    map[string]string
}

func (this *torDefaultSessionStorage) Init(ttl int64) {
	if this.datas != nil {
		return
	}
	this.ttl = ttl
	this.datas = make(map[string]torDefaultSessionStorageData)
	go this.gc()
}

func (this *torDefaultSessionStorage) gc() {
	for {
		if len(this.datas) > 0 {
			now := time.Now().Unix()
			for sid, data := range this.datas {
				if data.expires <= now {
					delete(this.datas, sid)
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func (this *torDefaultSessionStorage) CreateSessionID() string {
	t := time.Now()
	return "SESS" + fmt.Sprintf("%d%d", t.Unix(), t.Nanosecond())
}

func (this *torDefaultSessionStorage) Set(sid string, data map[string]string) {
	d := torDefaultSessionStorageData{
		expires: time.Now().Unix() + this.ttl,
		data:    data,
	}
	this.datas[sid] = d
}

func (this *torDefaultSessionStorage) Get(sid string) map[string]string {
	if data, exist := this.datas[sid]; exist {
		data.expires = time.Now().Unix() + this.ttl
		return data.data
	}
	return make(map[string]string)
}

func (this *torDefaultSessionStorage) Delete(sid string) {
	delete(this.datas, sid)
}
