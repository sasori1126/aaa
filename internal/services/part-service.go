package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/parts"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/platform/awsS3"
	otel2 "axis/ecommerce-backend/internal/platform/otel"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"context"
	"errors"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/xuri/excelize/v2"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"mime/multipart"
	"strconv"
)

type PartService interface {
	CreatePart(request *dto.PartRequest) *entities.ApiError
	UpdatePart(ctx context.Context, partHid string, request *dto.UpdatePartRequest, files []*multipart.FileHeader) *entities.ApiError
	UpdatePrice(ctx context.Context, request *dto.UpdatePriceRequest) *entities.ApiError
	MergeParts(request *dto.MergePartRequest) *entities.ApiError
	GetPartById(ctx context.Context, id string) (*dto.PartResponse, *entities.ApiError)
	GetParts(limit, offset int, returnZeroPrice bool) ([]dto.PartResponse, *entities.ApiError)
	GetDuplicateParts(limit, offset int, returnZeroPrice bool) ([]dto.DuplicatePartResponse, bool, *entities.ApiError)
	Search(limit, offset int, q string, returnZeroPrice bool) ([]dto.PartResponse, *entities.ApiError)
	LoadPriceDifference(limit, offset int) ([]dto.PartPriceDifferenceUpdate, *entities.ApiError)
}

type DefaultPartService struct {
	repo parts.PartRepo
}

func (d DefaultPartService) UpdatePart(ctx context.Context, partHid string, request *dto.UpdatePartRequest, files []*multipart.FileHeader) *entities.ApiError {
	ctx, span := otel.Tracer("").Start(ctx, "partServiceUpdatePart")
	defer span.End()
	id, err := models.DecodeHashId(partHid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	part, err := d.repo.GetPartById(ctx, id)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError(
				"Part not found",
				configs.NotFound,
				entities.E{},
			)
		}

		return utils.FormatApiError(
			"failed to retrieve model",
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	pr := models.Part{
		PartNumber:            request.PartNumber,
		Name:                  request.Name,
		Description:           request.Description,
		Detail:                request.Detail,
		Price:                 request.Price,
		PricePerMeter:         request.PricePerMeter,
		DealerPrice:           request.DealerPrice,
		DealerPricePercentage: request.DealerPricePercentage,
		DealerPricePerMeter:   request.DealerPricePerMeter,
		SalePrice:             request.SalePrice,
		SalePricePerMeter:     request.SalePricePerMeter,
		QuantityOnHand:        request.QuantityOnHand,
		QuantityOnOrder:       request.QuantityOnOrder,
		QuantityOnSaleOrder:   request.QuantityOnSaleOrder,
		QuantityRecommended:   request.QuantityRecommended,
		Weight:                request.Weight,
		Length:                request.Length,
		Width:                 request.Width,
		Height:                request.Height,
		Status:                request.Status,
		VideoUrl:              request.VideoUrl,
		Seo:                   request.Seo,
		MetaKeywords:          request.MetaKeywords,
		MetaDescription:       request.MetaDescription,
		OemNumber:             request.OemNumber,
		OemCompatible:         request.OemCompatible,
		Material:              request.Material,
		CountryOfOrigin:       request.CountryOfOrigin,
		Code:                  request.Code,
		GuideLink:             request.GuideLink,
		Featured:              request.Featured,
		CylinderType:          request.CylinderType,
		WhereUsed:             request.WhereUsed,
		Drivers:               request.Drivers,
	}
	pr.ID = part.ID

	s3, err := awsS3.NewS3Client()
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	images := part.Images
	for _, file := range files {
		iconImage, name, err := s3.CreateImage(file, 118, 118, "icon", part.ID)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed create icon image: "+err.Error(),
				configs.ServerError,
				entities.E{"server": "server error"},
			)
		}
		iconImagePath, err := s3.UploadImage(name, iconImage)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed create icon image: "+err.Error(),
				configs.ServerError,
				entities.E{"server": "server error"},
			)
		}

		previewImage, pName, err := s3.CreateImage(file, 300, 300, "preview", part.ID)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed create icon image: "+err.Error(),
				configs.ServerError,
				entities.E{"server": "server error"},
			)
		}
		previewImagePath, err := s3.UploadImage(pName, previewImage)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed create icon image: "+err.Error(),
				configs.ServerError,
				entities.E{"server": "server error"},
			)
		}

		listingImage, lName, err := s3.CreateImage(file, 600, 600, "listing", part.ID)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed create icon image: "+err.Error(),
				configs.ServerError,
				entities.E{"server": "server error"},
			)
		}
		listingImagePath, err := s3.UploadImage(lName, listingImage)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed create icon image: "+err.Error(),
				configs.ServerError,
				entities.E{"server": "server error"},
			)
		}

		images = append(images, models.Image{
			Name:        "icon_" + file.Filename,
			Description: "icon",
			Ext:         "png",
			Path:        iconImagePath,
		},
			models.Image{
				Name:        "preview_" + file.Filename,
				Description: "preview",
				Ext:         "png",
				Path:        previewImagePath,
			},
			models.Image{
				Name:        "listing_" + file.Filename,
				Description: "listing",
				Ext:         "png",
				Path:        listingImagePath,
			},
		)
	}

	pr.Images = images
	err = d.repo.UpdatePart(pr)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to update part",
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	return nil
}

func (d DefaultPartService) MergeParts(request *dto.MergePartRequest) *entities.ApiError {
	getParts, err := d.repo.GetDuplicates(request.PartNumber)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"Failed to retrieve parts",
			configs.ServerError,
			entities.E{},
		)
	}

	mainItemId, err := models.DecodeHashId(request.MainItem)
	if err != nil || mainItemId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"invalid item",
			configs.BadRequest,
			entities.E{"server": "invalid item"},
		)
	}

	var delIds []uint
	for _, part := range getParts {
		if part.ID != mainItemId {
			id := part.ID
			delIds = append(delIds, id)
		}
	}

	if len(delIds) < 1 {
		return utils.FormatApiError(
			"No duplicates found",
			configs.BadRequest,
			entities.E{},
		)
	}

	err = d.repo.MergeParts(mainItemId, delIds)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"Failed to merge items",
			configs.ServerError,
			entities.E{},
		)
	}
	return nil
}

func (d DefaultPartService) GetDuplicateParts(limit, offset int, returnZeroPrice bool) ([]dto.DuplicatePartResponse, bool, *entities.ApiError) {
	hasNext := false
	var duplicates []dto.DuplicatePartResponse
	getParts, err := d.repo.GetParts(limit+1, offset, returnZeroPrice)
	if err != nil {
		bugsnag.Notify(err)
		return nil, false, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	if len(getParts) > limit {
		hasNext = true
	}

	temp := map[string]bool{}
	for _, getPart := range getParts[:limit] {
		_, ok := temp[getPart.PartNumber]
		if !ok {
			pts, err := d.repo.GetDuplicates(getPart.PartNumber)
			if err != nil {
				bugsnag.Notify(err)
				return nil, false, utils.FormatApiError(
					"internal server error",
					configs.ServerError,
					entities.E{"server": "internal server error"},
				)
			}

			if len(pts) > 1 {
				var eParts []dto.EmbedPart
				for _, part := range pts {
					pid, _ := models.EncodeHashId(part.ID)

					var images []dto.ImageResponse
					for _, image := range part.Images {
						imHid, _ := models.EncodeHashId(image.ID)
						img := dto.ImageResponse{
							Id:   imHid,
							Name: image.Name,
							Path: image.Path,
						}

						images = append(images, img)
					}

					p := dto.EmbedPart{
						Id:                    pid,
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
						NoOfImages:            len(images),
					}
					eParts = append(eParts, p)
				}
				dub := dto.DuplicatePartResponse{
					PartNumber: getPart.PartNumber,
					Name:       getPart.Name,
					Parts:      eParts,
				}
				duplicates = append(duplicates, dub)
			}
			temp[getPart.PartNumber] = true
		}
	}

	if len(duplicates) == 0 {
		duplicates = []dto.DuplicatePartResponse{}
	}
	return duplicates, hasNext, nil
}

func (d DefaultPartService) UpdatePrice(ctx context.Context, request *dto.UpdatePriceRequest) *entities.ApiError {
	ctx, span := otel.Tracer("").Start(ctx, "userAddItem")
	defer span.End()
	part, err := d.repo.GetPartById(ctx, request.ID)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError(
				"Part not found",
				configs.NotFound,
				entities.E{},
			)
		}
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}
	err = d.repo.UpdateByField(part, models.FindByField{
		Field: "price",
		Value: request.NewPrice,
	})
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	return nil
}

func (d DefaultPartService) LoadPriceDifference(limit, offset int) ([]dto.PartPriceDifferenceUpdate, *entities.ApiError) {
	excel, err := excelize.OpenFile("./docs/price_list.xlsx")
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.NotFound,
			entities.E{"server": "internal server error"},
		)
	}
	defer excel.Close()

	if offset == 0 {
		return nil, utils.FormatApiError(
			"offset cannot equal to zero",
			configs.BadRequest,
			entities.E{"server": "offset cannot equal to zero"},
		)
	}

	lmt := offset + limit
	var partsPrice []dto.PartPriceDifferenceUpdate
	for index := offset; index <= lmt; index++ {
		n := strconv.Itoa(index + 1)
		price, err := excel.GetCellValue("Sheet1", "D"+n)
		if err != nil {
			bugsnag.Notify(err)
			return nil, utils.FormatApiError(
				"failed to load price on excel",
				configs.ServerError,
				entities.E{"server": err.Error()},
			)
		}

		code, err := excel.GetCellValue("Sheet1", "A"+n)
		if err != nil {
			bugsnag.Notify(err)
			return nil, utils.FormatApiError(
				"failed to load part code on excel",
				configs.ServerError,
				entities.E{"server": err.Error()},
			)
		}

		if price == "Price" || price == "" {
			continue
		}

		convertedPrice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			bugsnag.Notify(err)
			return nil, utils.FormatApiError(
				"internal server error",
				configs.ServerError,
				entities.E{"server": err.Error()},
			)
		}

		if convertedPrice > 0.0 {
			getPartByField, err := d.repo.GetPartByField(models.FindByField{
				Field: "part_number",
				Value: code,
			})
			if err != nil {
				bugsnag.Notify(err)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				} else {
					return nil, utils.FormatApiError(
						"internal server error",
						configs.ServerError,
						entities.E{"server": err.Error()},
					)
				}
			}
			for _, part := range getPartByField {
				if part.Price != convertedPrice {
					p := dto.PartPriceDifferenceUpdate{
						Id:                    part.ID,
						PartNumber:            part.PartNumber,
						Name:                  part.Name,
						Description:           part.Description,
						Detail:                part.Detail,
						PricePerMeter:         part.PricePerMeter,
						DealerPrice:           part.DealerPrice,
						DealerPricePercentage: part.DealerPricePercentage,
						DealerPricePerMeter:   part.DealerPricePerMeter,
						SalePrice:             part.SalePrice,
						SalePricePerMeter:     part.SalePricePerMeter,
						OemNumber:             part.OemNumber,
						OemCompatible:         part.OemCompatible,
						Code:                  part.Code,
						OldPrice:              part.Price,
						NewPrice:              convertedPrice,
					}
					partsPrice = append(partsPrice, p)
				}
			}
		}
	}

	if partsPrice == nil {
		partsPrice = []dto.PartPriceDifferenceUpdate{}
	}

	return partsPrice, nil
}

func (d DefaultPartService) Search(limit, offset int, q string, returnZeroPrice bool) ([]dto.PartResponse, *entities.ApiError) {
	res, err := d.repo.SearchParts(limit, offset, q, returnZeroPrice)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	response := make([]dto.PartResponse, 0, len(res))
	for _, re := range res {
		p := re.ToResponse()
		response = append(response, p)
	}

	return response, nil
}

func (d DefaultPartService) GetPartById(ctx context.Context, hid string) (*dto.PartResponse, *entities.ApiError) {
	ctx, span := otel.Tracer("").Start(ctx, "getPartById")
	defer span.End()
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		otel2.RecordSpanError(span, err)
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	part, err := d.repo.GetPartById(ctx, id)
	if err != nil {
		bugsnag.Notify(err)
		otel2.RecordSpanError(span, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.FormatApiError(
				"Part not found",
				configs.NotFound,
				entities.E{},
			)
		}
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	pr := part.ToResponse()

	return &pr, nil
}

func (d DefaultPartService) GetParts(limit, offset int, returnZeroPrice bool) ([]dto.PartResponse, *entities.ApiError) {
	getParts, err := d.repo.GetParts(limit, offset, returnZeroPrice)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var response []dto.PartResponse
	for _, getPart := range getParts {
		gp := getPart.ToResponse()
		response = append(response, gp)
	}

	return response, nil
}

func (d DefaultPartService) CreatePart(request *dto.PartRequest) *entities.ApiError {
	pr := models.Part{
		PartNumber:            request.PartNumber,
		Name:                  request.Name,
		Description:           request.Description,
		Detail:                request.Detail,
		Price:                 request.Price,
		PricePerMeter:         request.PricePerMeter,
		DealerPrice:           request.DealerPrice,
		DealerPricePercentage: request.DealerPricePercentage,
		DealerPricePerMeter:   request.DealerPricePerMeter,
		SalePrice:             request.SalePrice,
		SalePricePerMeter:     request.SalePricePerMeter,
		QuantityOnHand:        request.QuantityOnHand,
		QuantityOnOrder:       request.QuantityOnOrder,
		QuantityOnSaleOrder:   request.QuantityOnSaleOrder,
		QuantityRecommended:   request.QuantityRecommended,
		Weight:                request.Weight,
		Length:                request.Length,
		Width:                 request.Width,
		Height:                request.Height,
		Status:                request.Status,
		VideoUrl:              request.VideoUrl,
		Seo:                   request.Seo,
		MetaKeywords:          request.MetaKeywords,
		MetaDescription:       request.MetaDescription,
		OemNumber:             request.OemNumber,
		OemCompatible:         request.OemCompatible,
		Material:              request.Material,
		CountryOfOrigin:       request.CountryOfOrigin,
		Code:                  request.Code,
		GuideLink:             request.GuideLink,
		Featured:              request.Featured,
		CylinderType:          request.CylinderType,
		WhereUsed:             request.WhereUsed,
		Drivers:               request.Drivers,
	}

	images := []models.Image{
		{
			Name:        "demo1",
			Description: "demo1",
			Ext:         "jpg",
			Path:        "https://picsum.photos/500/300",
		},
		{
			Name:        "image2",
			Description: "image3",
			Ext:         "jpg",
			Path:        "https://picsum.photos/500/300",
		},
		{
			Name:        "image3",
			Description: "image3",
			Ext:         "jpg",
			Path:        "https://picsum.photos/500/300",
		},
	}
	pr.Images = images

	err := d.repo.CreatePart(pr)
	if err != nil {
		return utils.FormatApiError(
			"failed to part",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func NewDefaultPartService(repo parts.PartRepo) PartService {
	return &DefaultPartService{repo: repo}
}
