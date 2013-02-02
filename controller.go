package tor

import (
	"net/http"
)

type torControllerInterface interface {
	Init(*torApp, *torContext, *torTemplate, *torSession, string)
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Render()
	Output()
}

type Controller struct {
	app      *torApp
	Context  *torContext
	Template *torTemplate
	Session  *torSession
}

func (this *Controller) Init(app *torApp, ctx *torContext, tpl *torTemplate, sess *torSession, cn string) {
	this.app = app
	this.Context = ctx
	this.Context.ctlr = this
	this.Template = tpl
	this.Template.ctlr = this
	this.Session = sess
	this.Session.ctlr = this
}

func (this *Controller) Get() {
	http.Error(this.Context.Response, "Method Not Allowed", 405)
}

func (this *Controller) Post() {
	http.Error(this.Context.Response, "Method Not Allowed", 405)
}

func (this *Controller) Delete() {
	http.Error(this.Context.Response, "Method Not Allowed", 405)
}

func (this *Controller) Put() {
	http.Error(this.Context.Response, "Method Not Allowed", 405)
}

func (this *Controller) Head() {
	http.Error(this.Context.Response, "Method Not Allowed", 405)
}

func (this *Controller) Patch() {
	http.Error(this.Context.Response, "Method Not Allowed", 405)
}

func (this *Controller) Options() {
	http.Error(this.Context.Response, "Method Not Allowed", 405)
}

func (this *Controller) Render() {
	this.Template.Parse()
}

func (this *Controller) Output() {
	content := this.Template.GetResult()
	if len(content) > 0 {
		this.Context.WriteBytes(content)
	}
}

func (this *Controller) getHookController() *HookController {
	return &HookController{
		Context:  this.Context,
		Template: this.Template,
		Session:  this.Session,
	}
}
