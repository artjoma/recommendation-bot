package services

import (
	"recommendation-bot/dao"
)

type TriggerService struct {
	triggerDao *dao.TriggerDao
	triggers   *dao.TriggersModel // cache
}

func NewTriggerService(fileLocation string) (*TriggerService, error) {
	triggerDao := dao.NewTriggerDao(fileLocation)

	triggers, err := triggerDao.GetAll()
	if err != nil {
		return nil, err
	}

	return &TriggerService{triggerDao: triggerDao, triggers: triggers}, nil
}

// GetTriggers get phrases from cache
func (service *TriggerService) GetTriggers() *dao.TriggersModel {
	return service.triggers
}

func (service *TriggerService) AddTrigger(templateId string, trigger string) string {
	id := templateId + ":" + RandStringRunes(3)
	service.triggers.Triggers[id] = trigger
	// save to DB
	service.triggerDao.SaveAll(service.triggers)
	return id
}

func (service *TriggerService) DeleteTrigger(triggerId string) {
	delete(service.triggers.Triggers, triggerId)
	service.triggerDao.SaveAll(service.triggers)
}
