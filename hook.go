package tor

const (
	HookAfterInit           = "AfterInit"
	HookBeforeMethodGet     = "BeforeMethodGet"
	HookAfterMethodGet      = "AfterMethodGet"
	HookBeforeMethodPost    = "BeforeMethodPost"
	HookAfterMethodPost     = "AfterMethodPost"
	HookBeforeMethodHead    = "BeforeMethodHead"
	HookAfterMethodHead     = "AfterMethodHead"
	HookBeforeMethodDelete  = "BeforeMethodDelete"
	HookAfterMethodDelete   = "AfterMethodDelete"
	HookBeforeMethodPut     = "BeforeMethodPut"
	HookAfterMethodPut      = "AfterMethodPut"
	HookBeforeMethodPatch   = "BeforeMethodPatch"
	HookAfterMethodPatch    = "AfterMethodPatch"
	HookBeforeMethodOptions = "BeforeMethodOptions"
	HookAfterMethodOptions  = "AfterMethodOptions"
	HookBeforeRender        = "BeforeRender"
	HookAfterRender         = "AfterRender"
	HookBeforeOutput        = "BeforeOutput"
	HookAfterOutput         = "AfterOutput"
)

type torHook struct {
	app             *torApp
	controllerHooks map[string][]HookControllerFunc
}

type HookController struct {
	Context  *torContext
	Template *torTemplate
	Session  *torSession
}

type HookControllerFunc func(*HookController)

func (this *torHook) AddControllerHook(event string, hookFunc HookControllerFunc) {
	if this.controllerHooks == nil {
		this.controllerHooks = make(map[string][]HookControllerFunc)
	}
	if _, ok := this.controllerHooks[event]; !ok {
		this.controllerHooks[event] = []HookControllerFunc{}
	}
	this.controllerHooks[event] = append(this.controllerHooks[event], hookFunc)
}

func (this *torHook) CallControllerHook(event string, hc *HookController) {
	if funcList, ok := this.controllerHooks[event]; ok {
		for _, hookFunc := range funcList {
			hookFunc(hc)
			if hc.Context.Response.Finished {
				return
			}
		}
	}
}
