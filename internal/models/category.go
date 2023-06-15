package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type DiagramCat struct {
	configs.GormModel
	Name        string `gorm:"unique"`
	Description string
	//ImagePath      string
	DiagramSubCats []DiagramSubCat
	Diagrams       []Diagram `gorm:"many2many:diagram_diagram_cats;"`
}

type DiagramSubCat struct {
	configs.GormModel
	Name         string `gorm:"unique"`
	Description  string
	DiagramCatId uint
}

func (ds DiagramSubCat) ToSubCategoryResponse() dto.SubCategoryResponse {
	hid, _ := EncodeHashId(ds.ID)
	return dto.SubCategoryResponse{
		Id:          hid,
		Name:        ds.Name,
		Description: ds.Description,
	}
}

func (c DiagramCat) ToResponse() dto.CategoryResponse {
	hid, _ := EncodeHashId(c.ID)
	var subCats []dto.SubCategoryResponse
	subCats = []dto.SubCategoryResponse{}
	for _, subCat := range c.DiagramSubCats {
		sc := subCat.ToSubCategoryResponse()
		subCats = append(subCats, sc)
	}

	//var dgs []dto.EmbedDiagram
	//dgs = []dto.EmbedDiagram{}
	//for _, d := range c.Diagrams {
	//	id, _ := EncodeHashId(d.ID)
	//	ds := dto.EmbedDiagram{
	//		Id:          id,
	//		Name:        d.Name,
	//		Description: d.Description,
	//		BgImage:     d.BgImage,
	//	}
	//	dgs = append(dgs, ds)
	//}

	return dto.CategoryResponse{
		Id:            hid,
		Name:          c.Name,
		Description:   c.Description,
		SubCategories: subCats,
		//Diagrams:      dgs,
	}
}
