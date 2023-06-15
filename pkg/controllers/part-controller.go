package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/parts"
	"axis/ecommerce-backend/internal/dto"
	otel2 "axis/ecommerce-backend/internal/platform/otel"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"net/http"
	"path/filepath"
)

func (s *Serve) UpdatePart(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "userUpdatePart")
	defer span.End()
	pr := &dto.UpdatePartRequest{}
	err := c.ShouldBind(pr)
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

	hid := c.Param("partId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"order id": "invalid request"},
		))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		otel2.RecordSpanError(span, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "error getting form",
		})
		return
	}

	files := form.File["images[]"]
	if len(files) > 5 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "too many files, files should not exceed five",
		})
		return
	}

	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	apiErr := ps.UpdatePart(ctx, hid, pr, files)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"part updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreatePart(c *gin.Context) {
	pr := &dto.PartRequest{}
	err := c.ShouldBindJSON(pr)
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

	hs := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	apiErr := hs.CreatePart(pr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"part created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetParts(c *gin.Context) {
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

	allowZeroPrice := false
	_, ok := c.GetQuery("zero_priced")
	if ok {
		allowZeroPrice = true
	}

	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	getParts, apiErr := ps.GetParts(limit, offset, allowZeroPrice)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of parts",
		Data:    getParts,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) GetDuplicates(c *gin.Context) {
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

	allowZeroPrice := false
	_, ok := c.GetQuery("zero_priced")
	if ok {
		allowZeroPrice = true
	}

	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	getParts, hasNext, apiErr := ps.GetDuplicateParts(limit, offset, allowZeroPrice)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})

	next := links.Next
	if !hasNext {
		next = ""
	}
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "duplicate parts",
		Data:    getParts,
		Links: entities.E{
			"next":     next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) SearchParts(c *gin.Context) {
	q, ok := c.GetQuery("q")
	if !ok {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

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

	allowZeroPrice := false
	_, okay := c.GetQuery("zero_priced")
	if okay {
		allowZeroPrice = true
	}

	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	getParts, apiErr := ps.Search(limit, offset, q, allowZeroPrice)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of parts",
		Data:    getParts,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) GetPartById(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "userGetPart")
	defer span.End()
	hid := c.Param("partId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	part, apiErr := ps.GetPartById(ctx, hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "part detail",
		Data:    part,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := "price_list" + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, "./docs/"+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file" + err.Error(),
		})
		return
	}

	// File saved successfully. Return proper result
	c.JSON(http.StatusOK, gin.H{
		"message": "Your file has been successfully uploaded.",
	})
}

func (s *Serve) MergeParts(c *gin.Context) {
	pr := &dto.MergePartRequest{}
	err := c.ShouldBindJSON(pr)
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
	//log.Println(pr)
	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	apiErr := ps.MergeParts(pr)

	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) PriceUpdate(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "userUpdatePrice")
	defer span.End()
	pr := &dto.UpdatePriceRequest{}
	err := c.ShouldBindJSON(pr)
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
	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	apiErr := ps.UpdatePrice(ctx, pr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) LoadPrices(c *gin.Context) {
	qp := entities.QueryPathParam{
		Limit:  10,
		Offset: 1,
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

	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	priceDifference, apiErr := ps.LoadPriceDifference(limit, offset)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of parts difference",
		Data:    priceDifference,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}
