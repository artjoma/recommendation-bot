package dao

import (
	"time"
)

type TemplateDao struct {
	fileLocation string
	AbstractDao
}

func NewTemplateDao(fileLocation string) *TemplateDao {
	return &TemplateDao{
		fileLocation: fileLocation,
	}
}

// GetAll get all phrases from json file
func (dao *TemplateDao) GetAll() (*TemplatesModel, error) {
	var templates TemplatesModel
	if err := dao.GetAllItems(dao.fileLocation, &templates); err != nil {
		return nil, err
	}
	return &templates, nil
}

// SaveAll get all phrases from json file
func (dao *TemplateDao) SaveAll(templatesModel *TemplatesModel) error {
	templatesModel.Modon = time.Now().UTC()
	return dao.SaveAllItems(dao.fileLocation, templatesModel)
}
