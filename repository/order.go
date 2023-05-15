package repository

import (
	"fmt"
	"github.com/alekseikovrigin/yandexbackendschool2023/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	FindAllBy(data map[uint][]uint) []*models.Order
	FindCompleteByCourierId(courierID uint) []*models.Order
	Add(orders []models.Order) ([]models.Order, error)
	Save(order models.Order) (uint, error)
	Update(models.Order) error
	Delete(models.Order) error
	FindAll(offset int, limit int) []*models.Order
	FindAllAssign(offset int, limit int) []models.Order
	FindByID(orderID uint) (*models.Order, error)
	DeleteByID(orderID uint) error
	FindByName(name string) (*models.Order, error)
	FindByField(fieldName, fieldValue string) (*models.Order, error)
	UpdateSingleField(order models.Order, fieldName, fieldValue string) error
}
type orderDatabase struct {
	connection *gorm.DB
}

func NewOrderRepository() OrderRepository {
	if DB == nil {
		_, err = Connect()
		if err != nil {
			log.Error(err)
		}
	}
	return &orderDatabase{
		connection: DB,
	}
}

func (db orderDatabase) FindAllBy(data map[uint][]uint) []*models.Order {
	var orders []*models.Order

	clauses := make([]clause.Expression, 0)

	for _, val := range data {
		clauses = append(clauses, clause.Or(clause.And(
			gorm.Expr("id = ?", val[0]),
			gorm.Expr("courier_id = ?", val[1]),
			gorm.Expr("\"completed_time\" IS NULL"))))
	}
	//clauses = append(clauses, clause.And(gorm.Expr("completed_time IS NULL")))
	db.connection.Clauses(clauses...).Find(&orders)

	return orders
}

func (db orderDatabase) Add(orders []models.Order) ([]models.Order, error) {
	result := db.connection.Omit("CompletedTime").Create(&orders)
	log.Println(orders)
	if result.Error != nil {
		return orders, result.Error
	}
	return orders, nil
}

func (db orderDatabase) DeleteByID(orderID uint) error {
	order := models.Order{}
	order.ID = orderID
	result := db.connection.Delete(&order)
	return result.Error
}

func (db orderDatabase) Save(order models.Order) (uint, error) {
	result := db.connection.Create(&order)
	if result.Error != nil {
		return 0, result.Error
	}
	return order.ID, nil
}

func (db orderDatabase) Update(order models.Order) error {
	result := db.connection.Save(&order)
	return result.Error
}

func (db orderDatabase) Delete(order models.Order) error {
	result := db.connection.Delete(&order)
	return result.Error
}

func (db orderDatabase) FindAll(offset int, limit int) []*models.Order {
	var orders []*models.Order
	db.connection.Preload(clause.Associations).Limit(limit).Offset(offset).Find(&orders)
	return orders
}

func (db orderDatabase) FindAllAssign(offset int, limit int) []models.Order {
	var orders []models.Order
	db.connection.Preload(clause.Associations).Limit(limit).Offset(offset).Find(&orders)
	return orders
}

func (db orderDatabase) FindByID(orderID uint) (*models.Order, error) {
	var order models.Order
	result := db.connection.Preload(clause.Associations).Find(&order, "id = ?", orderID)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &order, nil
	}
	return nil, nil
}

func (db orderDatabase) FindCompleteByCourierId(courierID uint) []*models.Order { //TODO
	var orders []*models.Order
	db.connection.Preload(clause.Associations).Where("\"completed_time\" IS NOT NULL AND courier_id=?", courierID).Find(&orders)
	return orders
}

func (db orderDatabase) FindByName(name string) (*models.Order, error) {
	var order models.Order
	result := db.connection.Preload(clause.Associations).Find(&order, "name = ?", name)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &order, nil
	}
	return nil, nil
}

func (db orderDatabase) FindByField(fieldName, fieldValue string) (*models.Order, error) {
	var order models.Order
	result := db.connection.Preload(clause.Associations).Find(&order, fmt.Sprintf("%s = ?", fieldName), fieldValue)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &order, nil
	}
	return nil, nil
}

func (db orderDatabase) UpdateSingleField(order models.Order, fieldName, fieldValue string) error {
	result := db.connection.Model(&order).Update(fieldName, fieldValue)
	return result.Error
}
