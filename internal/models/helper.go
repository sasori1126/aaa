package models

import (
	"axis/ecommerce-backend/configs"
	"errors"
	"github.com/speps/go-hashids/v2"
)

func DecodeHashId(hashedId string) (uint, error) {
	if hashedId == "" {
		return 0, errors.New("failed to get id")
	}

	hi, err := createHash()
	if err != nil {
		return 0, err
	}

	id, err := hi.DecodeWithError(hashedId)
	if err != nil {
		return 0, err
	}
	originalId := uint(id[0])
	if originalId == 0 {
		return 0, errors.New("failed to get id")
	}
	return originalId, nil
}

func DecodeMultipleHash(hashids ...string) ([]uint, error) {
	var ids []uint
	hi, err := createHash()
	if err != nil {
		return nil, err
	}

	for _, hid := range hashids {
		id, err := hi.DecodeWithError(hid)
		if err != nil {
			return nil, err
		}
		getId := uint(id[0])
		if getId == 0 {
			return nil, errors.New("failed to get id")
		}
		ids = append(ids, getId)
	}

	return ids, nil
}

func createHash() (*hashids.HashID, error) {
	hd := hashids.NewData()
	hd.Salt = configs.AppConfig.AppKey
	hd.MinLength = configs.HashIdMinLength
	hashed, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	return hashed, nil
}

func EncodeHashId(id uint) (string, error) {
	hi, err := createHash()
	if err != nil {
		return "", err
	}

	hs, err := hi.Encode([]int{int(id)})
	if err != nil {
		return "", err
	}

	return hs, nil
}
