package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/diagrams"
	man_models "axis/ecommerce-backend/internal/domain/man-models"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/platform/awsS3"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (s *Serve) GetFigureImages(c *gin.Context) {
	q, _ := c.GetQuery("q")

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

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	getFigureImages, apiErr := ds.GetFigureImages(limit, offset, q)
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
		Data:    getFigureImages,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) AddDiagramImagesFromCsv(ctx *gin.Context) {
	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	apiErr := ds.AddFigureImagesFromCsv()
	if apiErr != nil {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, entities.ApiResponse{
		Message: "success",
		Data:    nil,
		Links:   nil,
		Meta:    nil,
	})
}

func (s *Serve) DeleteDiagramImage(ctx *gin.Context) {
	hid := ctx.Param("imageId")
	if hid == "" {
		ctx.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"diagramId": "invalid request"},
		))
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	apiErr := ds.DeleteFigureImage(hid)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, entities.ApiResponse{
		Message: "success",
		Data:    nil,
		Links:   nil,
		Meta:    nil,
	})
}

func (s *Serve) UpdateDiagramImages(ctx *gin.Context) {
	request := &dto.UploadFileRequest{}
	err := ctx.Bind(request)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	name := u.ID + strconv.Itoa(int(time.Now().UnixNano()))
	extension := ".png"
	getFileName, err := s.holdUploadImageTemp(ctx, name, extension)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file: " + err.Error(),
		})
		return
	}

	f, err := os.Open("./docs/tempFiles/" + getFileName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file" + err.Error(),
		})
		return
	}
	defer f.Close()
	defer os.Remove("./docs/tempFiles/" + getFileName)

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	res, apiErr := ds.UpdateFigureImage(request.ID, f, request.Title, name, extension)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}

	// File saved successfully. Return proper result
	ctx.JSON(http.StatusOK, res)

}

func (s *Serve) holdUploadImageTemp(ctx *gin.Context, name, ext string) (string, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return "", err
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	if extension != ext {
		return "", err
	}

	if file.Size > 800000 {
		return "", errors.New("file size not allowed")
	}

	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := name + extension

	// The file is received, so let's save it
	if err := ctx.SaveUploadedFile(file, "./docs/tempFiles/"+newFileName); err != nil {
		return "", errors.New("failed to save image")
	}

	return newFileName, nil
}

func (s *Serve) SaveDiagramImages(ctx *gin.Context) {
	request := &dto.UploadFileRequest{}
	err := ctx.Bind(request)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	name := u.ID + strconv.Itoa(int(time.Now().UnixNano()))
	extension := ".png"
	newFileName, err := s.holdUploadImageTemp(ctx, name, extension)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file: " + err.Error(),
		})
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	f, err := os.Open("./docs/tempFiles/" + newFileName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file" + err.Error(),
		})
		return
	}
	defer f.Close()
	defer os.Remove("./docs/tempFiles/" + newFileName)

	res, apiErr := ds.UploadFigureImage(f, request.Title, name, extension)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}

	// File saved successfully. Return proper result
	ctx.JSON(http.StatusOK, res)
}

func (s *Serve) DeleteDiagram(ctx *gin.Context) {
	hid := ctx.Param("diagramId")
	if hid == "" {
		ctx.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"diagramId": "invalid request"},
		))
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	apiErr := ds.DeleteDiagram(hid)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, entities.ApiResponse{
		Message: "success",
		Data:    nil,
		Links:   nil,
		Meta:    nil,
	})
}

func (s *Serve) UpdateDiagram(ctx *gin.Context) {
	request := &dto.DiagramImage{}
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	hid := ctx.Param("diagramId")
	if hid == "" {
		ctx.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	base64String := request.DiagramImg[strings.IndexByte(request.DiagramImg, ',')+1:]
	// decode image
	decodedImg, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entities.ApiError{
			Code:    http.StatusBadRequest,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	reader := bytes.NewReader(decodedImg)
	img, _, err := image.Decode(reader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// get diagram id
	diagramId, err := models.DecodeHashId(hid)
	if err != nil || diagramId == 0 {
		bugsnag.Notify(err)
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	// get diagram
	diagram, apiErr := ds.GetRawDiagramById(hid, true)
	if apiErr != nil {
		bugsnag.Notify(err)
		ctx.JSON(http.StatusInternalServerError, apiErr)
		return
	}

	parts := strings.Split(diagram.BgImage, "/")

	var nameIndex int
	for i, part := range parts {
		if part == "figures" {
			nameIndex = i + 1
			break
		}
	}

	imageName := parts[nameIndex]

	out, err := os.Create(imageName + ".png")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	s3, err := awsS3.NewS3Client()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get s3 session",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	f, err := os.Open(imageName + ".png")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer f.Close()

	res, err := s3.UploadImage(fmt.Sprintf("/figure_images/figures/%s/file.%s.png", imageName, "large"), f)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer func(name string) {
		log.Println(name)
		err := os.Remove(name + ".png")
		if err != nil {
			log.Println(err)
		}
	}(imageName)

	buff := new(bytes.Buffer)
	resized := imaging.Resize(img, 256, 256, imaging.Lanczos)
	err = png.Encode(buff, resized)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	resS, err := s3.UploadImage(fmt.Sprintf("/figure_images/figures/%s/file.%s.png", imageName, "small"), buff)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	apiErr = ds.UpdateFigDiagram(diagram, res, resS, request)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, entities.ApiResponse{
		Message: "success",
		Data:    nil,
		Links:   nil,
		Meta:    nil,
	})
}
func (s *Serve) SaveDiagram(ctx *gin.Context) {
	request := &dto.DiagramImage{}
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			ctx.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	base64String := request.DiagramImg[strings.IndexByte(request.DiagramImg, ',')+1:]
	// decode image
	decodedImg, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entities.ApiError{
			Code:    http.StatusBadRequest,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	reader := bytes.NewReader(decodedImg)
	img, _, err := image.Decode(reader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	t := time.Now()
	imageName := fmt.Sprintf("%d%02d%02dT%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	out, err := os.Create(imageName + ".png")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	s3, err := awsS3.NewS3Client()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	f, err := os.Open(imageName + ".png")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer f.Close()

	res, err := s3.UploadImage(fmt.Sprintf("/figure_images/figures/%s/file.%s.png", imageName, "large"), f)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer func(name string) {
		err := os.Remove(name + ".png")
		if err != nil {
			log.Println(err)
		}
	}(imageName)

	buff := new(bytes.Buffer)
	resized := imaging.Resize(img, 256, 256, imaging.Lanczos)
	err = png.Encode(buff, resized)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	resS, err := s3.UploadImage(fmt.Sprintf("/figure_images/figures/%s/file.%s.png", imageName, "small"), buff)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get image",
			Errors:  map[string]interface{}{"error": err.Error()},
		})
		return
	}

	ds := services.NewDefaultDiagramService(diagrams.NewDiagramRepoDb(s.Db), man_models.NewModelRepoDb(s.Db))
	apiErr := ds.CreateFigDiagram(res, resS, request)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, entities.ApiResponse{
		Message: "success",
		Data:    nil,
		Links:   nil,
		Meta:    nil,
	})
}
