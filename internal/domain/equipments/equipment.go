package equipments

import "axis/ecommerce-backend/internal/models"

type EquipmentRepo interface {
	AddEquipment(data models.Equipment) error
	GetEquipmentByField(query models.QueryByField) ([]models.Equipment, error)
}
