package services

import (
	"recommendation-bot/dao"
)

type TemplateService struct {
	templateDao *dao.TemplateDao
	templates   *dao.TemplatesModel // cache
}

func NewTemplateService(fileLocation string) (*TemplateService, error) {
	templateDao := dao.NewTemplateDao(fileLocation)
	templates, err := templateDao.GetAll()
	if err != nil {
		return nil, err
	}

	return &TemplateService{templateDao: templateDao, templates: templates}, nil
}

func (service *TemplateService) GetTemplates() *dao.TemplatesModel {
	return service.templates
}

func (service *TemplateService) AddTemplate(templateId string, trigger string) {
	service.templates.Templates[templateId] = trigger
	service.templateDao.SaveAll(service.templates)
}

func (service *TemplateService) DeleteTemplate(templateId string) {
	delete(service.templates.Templates, templateId)
	service.templateDao.SaveAll(service.templates)
}
