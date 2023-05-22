package dao

import (
	"time"
)

type TriggerDao struct {
	fileLocation string
	AbstractDao
}

func NewTriggerDao(fileLocation string) *TriggerDao {
	return &TriggerDao{
		fileLocation: fileLocation,
	}
}

// GetAll get all triggers from storage
func (dao *TriggerDao) GetAll() (*TriggersModel, error) {
	var triggers TriggersModel
	if err := dao.GetAllItems(dao.fileLocation, &triggers); err != nil {
		return nil, err
	}
	return &triggers, nil
}

// SaveAll save all triggers to storage
func (dao *TriggerDao) SaveAll(triggersModel *TriggersModel) error {
	triggersModel.Modon = time.Now().UTC()
	return dao.SaveAllItems(dao.fileLocation, triggersModel)
}
