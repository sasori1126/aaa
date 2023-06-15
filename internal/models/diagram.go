package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
	"encoding/json"
	"github.com/bugsnag/bugsnag-go/v2"
	"gorm.io/gorm"
)

type Diagram struct {
	configs.GormModel
	Name         string
	Description  string
	Images       []Image         `gorm:"many2many:diagram_images;"`
	Models       []Model         `gorm:"many2many:diagram_models;"`
	Parts        []Part          `gorm:"many2many:diagram_parts;"`
	Cats         []DiagramCat    `gorm:"many2many:diagram_diagram_cats;"`
	SubCats      []DiagramSubCat `gorm:"many2many:diagram_diagram_sub_cats;"`
	Status       string
	ControllerId *uint `gorm:"default:null"`
	Controller   Controller
	BgImage      string
	Thumbnail    string
	Draft        string
	Series       string
	IsAdmin      bool `gorm:"-:migration;->"`
}

type FigureImage struct {
	configs.GormModel
	Title        string
	SelectedType string
	File         string
}

func (i FigureImage) ToResponse() dto.FigureImageResponse {
	hid, _ := EncodeHashId(i.ID)
	return dto.FigureImageResponse{
		Id:    hid,
		Title: i.Title,
		File:  i.File,
	}
}

func (receiver Diagram) ToResponse() dto.DiagramResponse {
	hid, _ := EncodeHashId(receiver.ID)
	parts := make([]dto.EmbedPart, 0, len(receiver.Parts))
	for _, part := range receiver.Parts {
		id, _ := EncodeHashId(part.ID)
		images, types := formatPartData(part)

		d := dto.EmbedPart{
			Id:                    id,
			PartNumber:            part.PartNumber,
			PartDiagramNumber:     part.PartDiagramNumber,
			Name:                  part.Name,
			Description:           part.Description,
			Detail:                part.Detail,
			Price:                 part.Price,
			PricePerMeter:         part.PricePerMeter,
			DealerPrice:           part.DealerPrice,
			DealerPricePercentage: part.DealerPricePercentage,
			DealerPricePerMeter:   part.DealerPricePerMeter,
			SalePrice:             part.SalePrice,
			SalePricePerMeter:     part.SalePricePerMeter,
			QuantityOnHand:        part.QuantityOnHand,
			QuantityOnOrder:       part.QuantityOnOrder,
			QuantityOnSaleOrder:   part.QuantityOnSaleOrder,
			QuantityRecommended:   part.QuantityRecommended,
			Weight:                part.Weight,
			Length:                part.Length,
			Width:                 part.Width,
			Height:                part.Height,
			Status:                part.Status,
			VideoUrl:              part.VideoUrl,
			Seo:                   part.Seo,
			MetaKeywords:          part.MetaKeywords,
			MetaDescription:       part.MetaDescription,
			OemNumber:             part.OemNumber,
			OemCompatible:         part.OemCompatible,
			Material:              part.Material,
			CountryOfOrigin:       part.CountryOfOrigin,
			Code:                  part.Code,
			GuideLink:             part.GuideLink,
			Featured:              part.Featured,
			CylinderType:          part.CylinderType,
			WhereUsed:             part.WhereUsed,
			Drivers:               part.Drivers,
			Images:                images,
			Types:                 types,
		}

		parts = append(parts, d)
	}

	diagram := dto.DiagramResponse{
		Id:          hid,
		Name:        receiver.Name,
		Description: receiver.Description,
		Status:      receiver.Status,
		BgImage:     receiver.BgImage,
		Parts:       parts,
	}

	if receiver.Controller.ID != 0 {
		diagram.Controller = receiver.Controller.ToResponse()
	}

	if receiver.IsAdmin {
		series := make(map[string][]string)
		if receiver.Series != "" {
			err := json.Unmarshal([]byte(receiver.Series), &series)
			bugsnag.Notify(err)
		}

		models := make([]dto.EmbedModel, 0, len(receiver.Models))
		for _, model := range receiver.Models {
			hid, _ := EncodeHashId(receiver.ID)
			m := dto.EmbedModel{
				Id:    hid,
				Name:  model.Name,
				Image: model.ImagePath,
			}

			getSeries, ok := series[model.Name]
			if ok {
				m.Series = getSeries
			}

			models = append(models, m)
		}

		cats := make([]dto.CategoryResponse, 0, len(receiver.Cats))
		for _, cat := range receiver.Cats {
			c := cat.ToResponse()
			cats = append(cats, c)
		}

		diagram.Draft = receiver.Draft
		diagram.Cats = cats
		diagram.Models = models
	}

	return diagram
}

type DiagramPart struct {
	DiagramId         uint `gorm:"primaryKey"`
	PartId            uint `gorm:"primaryKey"`
	PartDiagramNumber int
	PartNumber        string
	PartDescription   string
	RecommendedQty    int
	CatNote           string
	Series            string
}

func (DiagramPart) BeforeCreate(db *gorm.DB) error {
	return nil
}

func formatPartData(part Part) ([]dto.ImageResponse, []dto.EmbedPartType) {
	var images []dto.ImageResponse
	images = []dto.ImageResponse{}
	for _, image := range part.Images {
		id, _ := EncodeHashId(image.ID)
		img := dto.ImageResponse{
			Id:   id,
			Name: image.Name,
			Path: image.Path,
		}
		images = append(images, img)
	}

	var types []dto.EmbedPartType
	types = []dto.EmbedPartType{}
	for _, tp := range part.Types {
		id, _ := EncodeHashId(tp.ID)
		typ := dto.EmbedPartType{
			Id:          id,
			Name:        tp.Name,
			Description: tp.Description,
		}
		types = append(types, typ)
	}

	return images, types
}
