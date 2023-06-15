package utils

import (
	"axis/ecommerce-backend/internal/models"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"time"
)

type ItemOrder struct {
	Title string `json:"title"`
}

func DownloadOrdersCsv(data []models.CustomerTempRequest) ([]byte, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	err := w.Write([]string{"Customer Name",  "Customer Email", "Customer Contact", "Customer Address", "Order", "Order Date"})
	if err != nil {
		return nil,err
	}

	for _, datum := range data {
		var record []string
		record = append(record, datum.CustomerName)
		record = append(record, datum.CustomerEmail)
		record = append(record, datum.CustomerContact)
		record = append(record, datum.CustomerAddress)
		order := datum.ItemsJson
		var itemData []ItemOrder
		err := json.Unmarshal([]byte(order), &itemData)
		if err != nil {
			return nil,err
		}
		orderDescription := ""
		for _, itemDatum := range itemData {
			orderDescription = orderDescription + itemDatum.Title + ", "
		}
		record = append(record, orderDescription)
		record = append(record, datum.CreatedAt.Format(time.RFC1123))

		err = w.Write(record)
		if err != nil {
			return nil,err
		}
	}
	w.Flush()

	return b.Bytes(), nil
}
