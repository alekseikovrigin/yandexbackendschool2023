package repository

import (
	"fmt"
	"github.com/alekseikovrigin/yandexbackendschool2023/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TypeRepository interface {
	Save(Type models.Type) (string, error)
	Update(models.Type) error
	Delete(models.Type) error
	FindAll() []*models.Type
	FindByID(TypeID uint) (*models.Type, error)
	DeleteByID(TypeID string) error
	FindByName(name string) (*models.Type, error)
	FindByField(fieldName, fieldValue string) (*models.Type, error)
	UpdateSingleField(Type models.Type, fieldName, fieldValue string) error
}
type TypeDatabase struct {
	connection *gorm.DB
}

func NewTypeRepository() TypeRepository {
	if DB == nil {
		_, err = Connect()
		if err != nil {
			log.Error(err)
		}
	}
	return &TypeDatabase{
		connection: DB,
	}
}

func (db TypeDatabase) DeleteByID(TypeID string) error {
	Type := models.Type{}
	Type.ID = TypeID
	result := db.connection.Delete(&Type)
	return result.Error
}

func (db TypeDatabase) Save(Type models.Type) (string, error) {
	result := db.connection.Create(&Type)
	if result.Error != nil {
		return "", result.Error
	}
	return Type.ID, nil
}

func (db TypeDatabase) Update(Type models.Type) error {
	result := db.connection.Save(&Type)
	return result.Error
}

func (db TypeDatabase) Delete(Type models.Type) error {
	result := db.connection.Delete(&Type)
	return result.Error
}

func (db TypeDatabase) FindAll() []*models.Type {
	var Types []*models.Type
	db.connection.Preload(clause.Associations).Find(&Types)
	return Types
}

func (db TypeDatabase) FindByID(TypeID uint) (*models.Type, error) {
	var Type models.Type
	result := db.connection.Preload(clause.Associations).Find(&Type, "id = ?", TypeID)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &Type, nil
	}
	return nil, nil
}

func (db TypeDatabase) FindByName(name string) (*models.Type, error) {
	var Type models.Type
	result := db.connection.Preload(clause.Associations).Find(&Type, "name = ?", name)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &Type, nil
	}
	return nil, nil
}

func (db TypeDatabase) FindByField(fieldName, fieldValue string) (*models.Type, error) {
	var Type models.Type
	result := db.connection.Preload(clause.Associations).Find(&Type, fmt.Sprintf("%s = ?", fieldName), fieldValue)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &Type, nil
	}
	return nil, nil
}

func (db TypeDatabase) UpdateSingleField(Type models.Type, fieldName, fieldValue string) error {
	result := db.connection.Model(&Type).Update(fieldName, fieldValue)
	return result.Error
}
