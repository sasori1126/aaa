package controllers

import (
	"axis/ecommerce-backend/configs"
	head_types "axis/ecommerce-backend/internal/domain/head-types"
	"axis/ecommerce-backend/internal/domain/heads"
	"axis/ecommerce-backend/internal/domain/sales"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"mime/multipart"
	"net/http"
)

func (s *Serve) GetHeadTypes(c *gin.Context) {
	hs := services.NewDefaultHeadTypeService(head_types.NewHeadTypeRepoDb(s.Db))
	hts, apiErr := hs.GetHeadTypes(-1, -1)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of head types",
		Data:    hts,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetHeadByHeadType(c *gin.Context) {
	hid := c.Param("headTypeId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	hs := services.NewDefaultHeadService(heads.NewHeadRepoDb(s.Db))
	getHeadsByType, apiErr := hs.GetHeadsType(hid, -1, -1)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "getHeadsByType by head type ",
		Data:    getHeadsByType,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) CreateHead(c *gin.Context) {
	headRequest := &dto.HeadRequest{}
	err := c.ShouldBind(headRequest)
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

	files := headRequest.Images
	var images []models.Image
	images = []models.Image{
		{
			Name:        "image1",
			Description: "image1",
			Ext:         "png",
			Path:        "https://picsum.photos/500/300",
		},
		{
			Name:        "image2",
			Description: "image2",
			Ext:         "png",
			Path:        "https://picsum.photos/500/300",
		},
	}
	//https://stackoverflow.com/questions/57085778/how-to-upload-images-to-s3-via-react-axios-and-golang-gin
	for _, file := range files {
		err := func(file *multipart.FileHeader) error {
			f, err := file.Open()
			if err != nil {
				return err
			}
			defer f.Close()
			size := file.Size
			buffer := make([]byte, size)
			_, err = f.Read(buffer)
			if err != nil {
				return err
			}
			_ = bytes.NewReader(buffer)
			fileType := http.DetectContentType(buffer)
			images = append(images, models.Image{
				Name:        file.Filename,
				Description: file.Filename,
				Ext:         fileType,
				Path:        "https://picsum.photos/500/300",
			})
			return nil
		}(file)
		if err != nil {
			log.Println(err, "failed to open image")
		}
	}

	hs := services.NewDefaultHeadService(heads.NewHeadRepoDb(s.Db))
	apiErr := hs.CreateHead(headRequest, images)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"head created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreateHeadType(c *gin.Context) {
	hr := &dto.HeadTypeRequest{}
	err := c.ShouldBindJSON(hr)
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

	hs := services.NewDefaultHeadTypeService(head_types.NewHeadTypeRepoDb(s.Db))
	apiErr := hs.CreateHeadType(hr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"head type created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetHeadById(c *gin.Context) {
	hid := c.Param("headId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	hs := services.NewDefaultHeadService(heads.NewHeadRepoDb(s.Db))
	head, apiErr := hs.GetHeadById(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "head detail",
		Data:    head,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) ControllerOrder(c *gin.Context) {
	rq := &dto.ControllerOrderRequest{}
	err := c.ShouldBind(rq)
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

	hs := services.NewDefaultGeneralService(sales.NewSaleRepoDb(s.Db))
	apiErr := hs.ControllerSaleOrder(rq)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"head created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) AxisHeadOrder(c *gin.Context) {
	rq := &dto.AxisHeadRequest{}
	err := c.ShouldBind(rq)
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
	hs := services.NewDefaultGeneralService(sales.NewSaleRepoDb(s.Db))
	apiErr := hs.AxisHeadSaleOrder(rq)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"head created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) KeslaHeadOrder(c *gin.Context) {
	rq := &dto.KeslaOrderRequest{}
	err := c.ShouldBind(rq)
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

	hs := services.NewDefaultGeneralService(sales.NewSaleRepoDb(s.Db))
	apiErr := hs.KeslaSaleOrder(rq)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"head created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetHeads(c *gin.Context) {
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
	hs := services.NewDefaultHeadService(heads.NewHeadRepoDb(s.Db))
	getHeads, apiErr := hs.GetHeads(limit, offset)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of getHeads",
		Data:    getHeads,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}
