package tor

import (
	"html/template"
	"io/ioutil"
)

var tplFuncMap template.FuncMap

func init() {
	tplFuncMap = make(template.FuncMap)
}

func AddTemplateFunc(name string, tplFunc interface{}) {
	tplFuncMap[name] = tplFunc
}

type torTemplate struct {
	ctlr      *Controller
	tpl       *template.Template
	tplVars   map[string]interface{}
	tplResult *torTemplateResult
}

func (this *torTemplate) SetVar(name string, value interface{}) {
	this.tplVars[name] = value
}

func (this *torTemplate) GetVar(name string) interface{} {
	if val, ok := this.tplVars[name]; ok {
		return val
	}
	return nil
}

func (this *torTemplate) SetTemplateString(str string) bool {
	this.tpl = template.New("")
	this.tpl.Parse(str)
	return true
}

func (this *torTemplate) SetTemplateFile(filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	return this.SetTemplateString(string(content))
}

func (this *torTemplate) SetSubTemplateString(name, str string) bool {
	if this.tpl == nil {
		return false
	}
	tpl := this.tpl.New(name)
	tpl.Parse(`{{define "` + name + `"}}` + str + `{{end}}`)
	return true
}

func (this *torTemplate) SetSubTemplateFile(name, filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	return this.SetSubTemplateString(name, string(content))
}

func (this *torTemplate) Parse() bool {
	if this.tpl == nil {
		return false
	}
	if this.tplResult != nil {
		return false
	}

	hc := this.ctlr.getHookController()
	this.ctlr.app.callControllerHook("BeforeRender", hc)
	if this.ctlr.Context.Response.Finished {
		return true
	}

	this.tplResult = &torTemplateResult{data: []byte{}}
	this.tpl.Funcs(tplFuncMap)
	err := this.tpl.Execute(this.tplResult, this.tplVars)
	if err != nil {
		return false
	}

	this.ctlr.app.callControllerHook("AfterRender", hc)
	return true
}

func (this *torTemplate) GetResult() []byte {
	if this.tplResult == nil {
		return []byte{}
	}
	return this.tplResult.Bytes()
}

func (this torTemplate) GetResultString() string {
	return string(this.GetResult())
}

func (this *torTemplate) SetResult(p []byte) {
	if this.tplResult == nil {
		return
	}
	this.tplResult.SetBytes(p)
}

func (this *torTemplate) SetResultString(s string) {
	this.SetResult([]byte(s))
}

type torTemplateResult struct {
	data []byte
}

func (this *torTemplateResult) Write(p []byte) (n int, err error) {
	this.data = append(this.data, p...)
	return len(p), nil
}

func (this *torTemplateResult) SetBytes(p []byte) {
	this.data = p
}

func (this *torTemplateResult) String() string {
	return string(this.data)
}

func (this *torTemplateResult) Bytes() []byte {
	return this.data
}
