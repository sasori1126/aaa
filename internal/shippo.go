package internal

import (
	"fmt"
	"strconv"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"

	shippo "github.com/coldbrewcloud/go-shippo"
	"github.com/coldbrewcloud/go-shippo/client"
	shippoModels "github.com/coldbrewcloud/go-shippo/models"
)

type Shippo struct {
	client *client.Client
}

func (s *Shippo) CreateShipment(
	fromAddressInput shippoModels.AddressInput,
	toAddressInput shippoModels.AddressInput,
	weight float32,
	parcel dto.ShippoParcel,
) (*shippoModels.Shipment, error) {
	addressFrom, err := s.client.CreateAddress(&fromAddressInput)
	if err != nil {
		return nil, err
	}

	addressTo, err := s.client.CreateAddress(&toAddressInput)
	if err != nil {
		return nil, err
	}

	p, err := s.client.CreateParcel(&shippoModels.ParcelInput{
		DistanceUnit: parcel.DistanceUnit,
		Height:       parcel.Height,
		Length:       parcel.Length,
		MassUnit:     parcel.MassUnit,
		Weight:       fmt.Sprintf("%.4f", weight),
		Width:        parcel.Width,
	})
	if err != nil {
		return nil, err
	}

	shipment, err := s.client.CreateShipment(&shippoModels.ShipmentInput{
		AddressFrom: addressFrom.ObjectID,
		AddressTo:   addressTo.ObjectID,
		Async:       false,
		Parcels:     []string{p.ObjectID},
	})
	if err != nil {
		return nil, err
	}

	for _, rate := range shipment.Rates {
		amount, err := strconv.ParseFloat(rate.Amount, 32)
		if err != nil {
			return nil, err
		}
		// Add extra 20% based on the Shippo returned shipping amount. This value need to be adjustable in the admin portal.
		rate.Amount = fmt.Sprintf("%.2f", amount*1.2)
	}
	return shipment, nil
}

func (s *Shippo) GetRate(rateObjectId string) (*shippoModels.Rate, error) {
	rate, err := s.client.RetrieveRate(rateObjectId)
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func NewShippoClient() (*Shippo, error) {
	token, err := configs.GetShippoToken()
	if err != nil {
		return nil, err
	}

	c := shippo.NewClient(token)
	return &Shippo{client: c}, nil
}
