package repository

import (
	"fmt"
	"github.com/alekseikovrigin/yandexbackendschool2023/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CourierRepository interface {
	Add(couriers []models.Courier) ([]models.Courier, error)
	Save(courier models.Courier) (uint, error)
	Update(models.Courier) error
	Delete(models.Courier) error
	FindAll(offset int, limit int) []*models.Courier
	FindByID(courierID uint) (*models.Courier, error)
	DeleteByID(courierID uint) error
	FindByName(name string) (*models.Courier, error)
	FindByField(fieldName, fieldValue string) (*models.Courier, error)
	UpdateSingleField(courier models.Courier, fieldName, fieldValue string) error
}
type courierDatabase struct {
	connection *gorm.DB
}

func NewCourierRepository() CourierRepository {
	if DB == nil {
		_, err = Connect()
		if err != nil {
			log.Error(err)
		}
	}
	return &courierDatabase{
		connection: DB,
	}
}

func (db courierDatabase) Add(couriers []models.Courier) ([]models.Courier, error) {
	result := db.connection.Preload(clause.Associations).Create(&couriers)
	if result.Error != nil {
		return couriers, result.Error
	}
	return couriers, nil
}

func (db courierDatabase) DeleteByID(courierID uint) error {
	courier := models.Courier{}
	courier.ID = courierID
	result := db.connection.Delete(&courier)
	return result.Error
}

func (db courierDatabase) Save(courier models.Courier) (uint, error) {
	result := db.connection.Create(&courier)
	if result.Error != nil {
		return 0, result.Error
	}
	return courier.ID, nil
}

func (db courierDatabase) Update(courier models.Courier) error {
	result := db.connection.Save(&courier)
	return result.Error
}

func (db courierDatabase) Delete(courier models.Courier) error {
	result := db.connection.Delete(&courier)
	return result.Error
}

func (db courierDatabase) FindAll(offset int, limit int) []*models.Courier {
	var couriers []*models.Courier
	db.connection.Preload(clause.Associations).Limit(limit).Offset(offset).Find(&couriers)
	return couriers
}

func (db courierDatabase) FindByID(courierID uint) (*models.Courier, error) {
	var courier models.Courier
	result := db.connection.Preload(clause.Associations).Find(&courier, "id = ?", courierID)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &courier, nil
	}
	return nil, nil
}

func (db courierDatabase) FindByName(name string) (*models.Courier, error) {
	var courier models.Courier
	result := db.connection.Preload(clause.Associations).Find(&courier, "name = ?", name)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &courier, nil
	}
	return nil, nil
}

func (db courierDatabase) FindByField(fieldName, fieldValue string) (*models.Courier, error) {
	var courier models.Courier
	result := db.connection.Preload(clause.Associations).Find(&courier, fmt.Sprintf("%s = ?", fieldName), fieldValue)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &courier, nil
	}
	return nil, nil
}

func (db courierDatabase) UpdateSingleField(courier models.Courier, fieldName, fieldValue string) error {
	result := db.connection.Model(&courier).Update(fieldName, fieldValue)
	return result.Error
}
