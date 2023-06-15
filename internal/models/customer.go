package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
	"encoding/json"
)

type CustomerTempRequest struct {
	configs.GormModel
	CustomerName string
	CustomerEmail string
	CustomerContact string
	CustomerAddress string
	ItemsJson string
}

func (c *CustomerTempRequest) ConvertItemsToJson(items []dto.Item) error {
	convertedItems, err := json.Marshal(items)
	if err != nil {
		return err
	}
	c.ItemsJson = string(convertedItems)
	return nil
}
