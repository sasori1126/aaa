package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/categories"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
)

type CategoryService interface {
	GetCategories(limit, offset int) ([]dto.CategoryResponse, *entities.ApiError)
	CreateDiagramCat(request *dto.CategoryRequest) *entities.ApiError
	UpdateDiagramCat(hid string, request *dto.CategoryRequest) *entities.ApiError
	DeleteDiagramCat(hid string) *entities.ApiError
	CreateDiagramSubCat(request *dto.SubCategoryRequest) *entities.ApiError
	GetCategoryById(hid string) (*dto.CategoryResponse, *entities.ApiError)
}

type DefaultCatService struct {
	repo categories.CategoryRepo
}

func (d DefaultCatService) UpdateDiagramCat(hid string, req *dto.CategoryRequest) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	data := models.DiagramCat{
		Name:        req.Name,
		Description: req.Description,
	}
	data.ID = id

	err = d.repo.UpdateDiagramCat(&data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"cat_id": "failed to update"},
		)
	}
	return nil
}

func (d DefaultCatService) DeleteDiagramCat(hid string) *entities.ApiError {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	err = d.repo.DeleteCategoryById(id)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"cat_id": "failed to delete"},
		)
	}

	return nil
}

func (d DefaultCatService) CreateDiagramCat(request *dto.CategoryRequest) *entities.ApiError {
	data := models.DiagramCat{
		Name:        request.Name,
		Description: request.Description,
	}

	err := d.repo.CreateDiagramCat(data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create diagram category",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultCatService) CreateDiagramSubCat(request *dto.SubCategoryRequest) *entities.ApiError {
	catId, err := models.DecodeHashId(request.CategoryId)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error, could not decode id"},
		)
	}

	data := models.DiagramSubCat{
		DiagramCatId: catId,
		Name:         request.Name,
		Description:  request.Description,
	}

	err = d.repo.CreateDiagramSubCat(data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create diagram sub category",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultCatService) GetCategoryById(hid string) (*dto.CategoryResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(hid)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	cat, err := d.repo.GetCategoryById(id)
	if err != nil {
		bugsnag.Notify(err)
		return nil, nil
	}

	var catResponse dto.CategoryResponse = cat.ToResponse()

	return &catResponse, nil
}

func (d DefaultCatService) GetCategories(limit, offset int) ([]dto.CategoryResponse, *entities.ApiError) {
	cats, err := d.repo.GetCategories(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to get categoryResponses",
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	var categoryResponses []dto.CategoryResponse
	categoryResponses = []dto.CategoryResponse{}
	for _, cat := range cats {
		ct := cat.ToResponse()
		categoryResponses = append(categoryResponses, ct)
	}

	return categoryResponses, nil
}

func NewDefaultCatService(repo categories.CategoryRepo) CategoryService {
	return &DefaultCatService{repo: repo}
}
