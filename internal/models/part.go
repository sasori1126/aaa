package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
	"time"
)

type Part struct {
	configs.GormModel
	PartNumber            string
	Name                  string
	Description           string
	Detail                string
	Price                 float64
	PricePerMeter         float64
	DealerPrice           float64
	DealerPricePercentage int
	DealerPricePerMeter   float64
	SalePrice             float64
	SalePricePerMeter     float64
	QuantityOnHand        int
	QuantityOnOrder       int
	QuantityOnSaleOrder   int
	QuantityRecommended   int
	Weight                float64
	Length                float64
	Width                 float64
	Height                float64
	Status                string
	VideoUrl              string
	Seo                   string
	MetaKeywords          string
	MetaDescription       string
	OemNumber             string
	OemCompatible         string
	Material              string
	CountryOfOrigin       string
	ListId                string
	Code                  string
	GuideLink             string
	Featured              bool
	CylinderType          string
	WhereUsed             string
	Drivers               string
	Images                []Image    `gorm:"many2many:part_images;"`
	Models                []Model    `gorm:"many2many:part_models;"`
	Types                 []PartType `gorm:"many2many:part_part_types;"`
	SyncDate              time.Time
	PartDiagramNumber     int    `gorm:"-:migration;->"`
	CatNote               string `gorm:"-:migration;->"`
	RecommendedQty        int    `gorm:"-:migration;->"`
}

func (p Part) ToResponse() dto.PartResponse {
	hid, _ := EncodeHashId(p.ID)
	var images []dto.ImageResponse = []dto.ImageResponse{}
	for _, image := range p.Images {
		id, _ := EncodeHashId(image.ID)
		img := dto.ImageResponse{
			Id:   id,
			Name: image.Name,
			Path: image.Path,
		}
		images = append(images, img)
	}
	return dto.PartResponse{
		Id:                    hid,
		PartNumber:            p.PartNumber,
		Name:                  p.Name,
		Description:           p.Description,
		Detail:                p.Detail,
		Price:                 p.Price,
		PricePerMeter:         p.PricePerMeter,
		DealerPrice:           p.DealerPrice,
		DealerPricePercentage: p.DealerPricePercentage,
		DealerPricePerMeter:   p.DealerPricePerMeter,
		SalePrice:             p.SalePrice,
		SalePricePerMeter:     p.SalePricePerMeter,
		QuantityOnHand:        p.QuantityOnHand,
		QuantityOnOrder:       p.QuantityOnOrder,
		QuantityOnSaleOrder:   p.QuantityOnSaleOrder,
		QuantityRecommended:   p.QuantityRecommended,
		Weight:                p.Weight,
		Length:                p.Length,
		Width:                 p.Width,
		Height:                p.Height,
		Status:                p.Status,
		VideoUrl:              p.VideoUrl,
		Seo:                   p.Seo,
		MetaKeywords:          p.MetaKeywords,
		MetaDescription:       p.MetaDescription,
		OemNumber:             p.OemNumber,
		OemCompatible:         p.OemCompatible,
		Material:              p.Material,
		CountryOfOrigin:       p.CountryOfOrigin,
		ListId:                p.ListId,
		Code:                  p.Code,
		GuideLink:             p.GuideLink,
		Featured:              p.Featured,
		CylinderType:          p.CylinderType,
		WhereUsed:             p.WhereUsed,
		Drivers:               p.Drivers,
		Images:                images,
	}
}
