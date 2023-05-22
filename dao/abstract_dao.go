package dao

import (
	"encoding/json"
	"io"
	"os"
)

// AbstractDao currently we are using file storage as data persistence
type AbstractDao struct {
}

func (dao *AbstractDao) GetAllItems(fileLocation string, entity interface{}) error {
	jsonFile, err := os.Open(fileLocation)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, entity); err != nil {
		return err
	}
	return nil
}

func (dao *AbstractDao) SaveAllItems(fileLocation string, entity interface{}) error {
	data, err := json.MarshalIndent(entity, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileLocation, data, 0644)
}
