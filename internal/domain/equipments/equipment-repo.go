package equipments

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type EquipmentRepoDb struct {
	client storage.Storage
}

func (e EquipmentRepoDb) AddEquipment(data models.Equipment) error {
	return e.client.AddEquipment(data)
}

func (e EquipmentRepoDb) GetEquipmentByField(query models.QueryByField) ([]models.Equipment, error) {
	return e.client.GetEquipmentByQuery(query)
}

func NewEquipmentRepoDb(db storage.Storage) EquipmentRepo {
	return &EquipmentRepoDb{client: db}
}
