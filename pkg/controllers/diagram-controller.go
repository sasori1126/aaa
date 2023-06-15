package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/categories"
	"axis/ecommerce-backend/internal/domain/diagrams"
	man_models "axis/ecommerce-backend/internal/domain/man-models"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (s *Serve) CreateDiagramSubCat(c *gin.Context) {
	dscRequest := &dto.SubCategoryRequest{}
	err := c.ShouldBindJSON(dscRequest)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	apiErr := cs.CreateDiagramSubCat(dscRequest)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"diagram sub category created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreateDiagram(c *gin.Context) {
	dr := &dto.DiagramRequest{}
	err := c.ShouldBindJSON(dr)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	apiErr := ds.CreateDiagram(dr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"diagram created successfully",
		utils.Empty(), utils.Empty()))

}

func (s *Serve) UpdateDiagramCat(c *gin.Context) {
	hid := c.Param("catId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	request := &dto.CategoryRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}
	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	apiErr := cs.UpdateDiagramCat(hid, request)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"cat updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) DeleteDiagramCat(c *gin.Context) {
	hid := c.Param("catId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	apiErr := cs.DeleteDiagramCat(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"category deleted successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreateDiagramCat(c *gin.Context) {
	dcRequest := &dto.CategoryRequest{}
	err := c.ShouldBindJSON(dcRequest)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	apiErr := cs.CreateDiagramCat(dcRequest)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"diagram category created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetDiagramCat(c *gin.Context) {
	hid := c.Param("categoryId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"category_id": "invalid request"},
		))
		return
	}

	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	getCategoryById, apiErr := cs.GetCategoryById(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "category detail",
		Data:    getCategoryById,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetDiagramById(c *gin.Context) {
	hid := c.Param("diagramId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	isAdmin := false
	_, ok := c.GetQuery("isAdmin")
	if ok {
		isAdmin = true
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	diagram, apiErr := ds.GetDiagramById(hid, isAdmin)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "diagram detail",
		Data:    diagram,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetDiagrams(c *gin.Context) {
	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	cats, apiErr := cs.GetCategories(-1, -1)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	//get getModels
	ms := services.NewDefaultModelService(man_models.NewModelRepoDb(s.Db))
	getModels, apiErr := ms.GetModels(-1, -1)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	getDiagrams, apiErr := ds.GetDiagrams(10, 0, nil, nil, false)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "diagram index",
		Data: entities.E{
			"categories": cats,
			"models":     getModels,
			"diagrams":   getDiagrams,
		},
	})
}

func (s *Serve) GetAdminDiagrams(c *gin.Context) {
	qp := entities.QueryPathParam{
		Limit:  10,
		Offset: 0,
	}
	limit, offset, err := queryParams(c, qp)
	if err != nil {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

	f := false
	_, ok := c.GetQuery("figures")
	if ok {
		f = true
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	getDiagrams, apiErr := ds.Diagrams(limit, offset, f)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of diagrams",
		Data:    getDiagrams,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) GetDiagramsSearch(c *gin.Context) {
	var modelIds []string
	var catIds []string
	modelId, ok := c.GetQuery("model_id")
	if ok {
		modelIds = append(modelIds, modelId)
	}
	catId, ok := c.GetQuery("cat_id")
	if ok {
		catIds = append(catIds, catId)
	}

	f := false
	_, ok = c.GetQuery("figures")
	if ok {
		f = true
	}

	defValues := entities.QueryPathParam{
		Limit:  10,
		Offset: 0,
	}

	limit, offset, err := queryParams(c, defValues)
	if err != nil {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	getDiagrams, apiErr := ds.GetDiagrams(limit, offset, modelIds, catIds, f)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, utils.GetApiResponse(getDiagrams, "list of diagrams", entities.E{}, entities.E{}))
}

func (s *Serve) GetDiagramCats(c *gin.Context) {
	qp := entities.QueryPathParam{
		Limit:  -1,
		Offset: -1,
	}
	limit, offset, err := queryParams(c, qp)
	if err != nil {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}
	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	cats, apiErr := cs.GetCategories(limit, offset)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of diagram categories",
		Data:    cats,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}
