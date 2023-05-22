package services

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type AppCfg struct {
	TokenId string  `yaml:"tokenId"`
	Admins  []int64 `yaml:"admins"`
}

type AppCtx struct {
	active          bool
	appCfg          *AppCfg
	triggerService  *TriggerService
	templateService *TemplateService
}

func NewAppCtx(appCfg *AppCfg, triggerService *TriggerService, templateService *TemplateService) *AppCtx {
	return &AppCtx{
		appCfg:          appCfg,
		triggerService:  triggerService,
		templateService: templateService,
	}
}

func (appCtx *AppCtx) GetTriggerService() *TriggerService {
	return appCtx.triggerService
}

func (appCtx *AppCtx) SetActiveState() {
	appCtx.active = true
}
func (appCtx *AppCtx) Destroy() {
	appCtx.active = false
}

func (appCtx AppCtx) GetAppCfg() *AppCfg {
	return appCtx.appCfg
}

const (
	AddTemplateCmd    string = "addtemplate"
	RemoveTemplateCmd        = "rmtemplate"
	ListTemplatesCmd         = "lstemplates"
	AddTriggerCmd            = "addtrigger"
	RemoveTriggerCmd         = "rmtrigger"
	ListTiggersCmd           = "lstriggers"
)
