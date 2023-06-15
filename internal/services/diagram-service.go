package services

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/diagrams"
	man_models "axis/ecommerce-backend/internal/domain/man-models"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/platform/awsS3"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/bugsnag/bugsnag-go/v2"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

type DiagramService interface {
	CreateDiagram(request *dto.DiagramRequest) *entities.ApiError
	CreateFigDiagram(l, s string, request *dto.DiagramImage) *entities.ApiError
	UpdateFigDiagram(diagram *models.Diagram, l, s string, request *dto.DiagramImage) *entities.ApiError
	UploadFigureImage(file io.Reader, title, fileName, ext string) (dto.FigureImageResponse, *entities.ApiError)
	UpdateFigureImage(id string, file io.Reader, title, fileName, ext string) (dto.FigureImageResponse, *entities.ApiError)
	GetFigureImage(hid string) (*models.FigureImage, *entities.ApiError)
	GetFigureImages(limit, offset int, search string) ([]dto.FigureImageResponse, *entities.ApiError)
	GetDiagrams(limit, offset int, modelIds, catIds []string, f bool) ([]dto.DiagramResponse, *entities.ApiError)
	Diagrams(limit, offset int, f bool) ([]dto.DiagramResponse, *entities.ApiError)
	GetDiagramById(id string, isAdmin bool) (*dto.DiagramResponse, *entities.ApiError)
	DeleteDiagram(id string) *entities.ApiError
	DeleteFigureImage(id string) *entities.ApiError
	AddFigureImagesFromCsv() *entities.ApiError
	GetRawDiagramById(id string, isAdmin bool) (*models.Diagram, *entities.ApiError)
}

type DefaultDiagramService struct {
	repo      diagrams.DiagramRepo
	modelRepo man_models.ModelRepo
}

func (d DefaultDiagramService) AddFigureImagesFromCsv() *entities.ApiError {
	file, err := os.Open("./docs/figure_images.csv")
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to load file: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}
	defer file.Close()

	fileReader := csv.NewReader(file)
	records, err := fileReader.ReadAll()
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to read records: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	images := make([]models.FigureImage, 0, len(records))
	for _, record := range records {
		if record[0] != "id" {
			file := record[6]

			f := struct {
				FileName     string
				Url          string
				Video        string
				Description  string
				SelectedType string
			}{}

			err = json.Unmarshal([]byte(file), &f)
			if err != nil {
				bugsnag.Notify(err)
				return utils.FormatApiError(
					"failed to read records: "+err.Error(),
					configs.ServerError,
					entities.E{"server": "server error"},
				)
			}

			selectType := "image"
			if f.SelectedType != "" {
				selectType = f.SelectedType
			}

			if f.Url != "" {
				title := f.FileName
				if f.FileName != "" {
					fileSplit := strings.Split(f.FileName, ".")
					title = fileSplit[0]
				}

				img := models.FigureImage{
					Title:        title,
					SelectedType: selectType,
					File:         "https:" + f.Url,
				}

				images = append(images, img)
			}
		}
	}

	err = d.repo.SaveFromCsv(images)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to read records: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}
	return nil
}

func (d DefaultDiagramService) GetFigureImage(hid string) (*models.FigureImage, *entities.ApiError) {
	imageId, err := models.DecodeHashId(hid)
	if err != nil || imageId == 0 {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	image, err := d.repo.GetFigImage(imageId)
	if err != nil {
		return nil, utils.FormatApiError(
			"failed to get image",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	return image, nil
}

func (d DefaultDiagramService) DeleteFigureImage(hid string) *entities.ApiError {
	getImage, customErr := d.GetFigureImage(hid)
	if customErr != nil {
		return utils.FormatApiError(
			"failed to find image: "+customErr.Message,
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	err := d.repo.DeleteFigImage(getImage.ID)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to delete image",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	s3, err := awsS3.NewS3Client()
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to start s3",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	parsedLargeFileUrl, err := url.Parse(getImage.File)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to parse image url",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	key := parsedLargeFileUrl.Path[1:]
	_ = s3.DeleteResource("axisforestry", key)

	return nil
}

func (d DefaultDiagramService) UpdateFigureImage(hid string, file io.Reader, title, fileName, ext string) (dto.FigureImageResponse, *entities.ApiError) {
	getImage, customErr := d.GetFigureImage(hid)
	if customErr != nil {
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to find image: "+customErr.Message,
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}
	oldImagePath := getImage.File
	getImage.Title = title

	s3, err := awsS3.NewS3Client()
	if err != nil {
		bugsnag.Notify(err)
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	fName := fmt.Sprintf("/figure_images/figures/%s/file.%s.png", fileName, "raw")
	s3Path, err := s3.UploadImage(fName, file)
	if err != nil {
		bugsnag.Notify(err)
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	getImage.File = s3Path

	err = d.repo.UpdateFigImage(getImage)
	if err != nil {
		bugsnag.Notify(err)
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	parsedLargeFileUrl, err := url.Parse(oldImagePath)
	if err != nil {
		bugsnag.Notify(err)
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	key := parsedLargeFileUrl.Path[1:]
	_ = s3.DeleteResource("axisforestry", key)

	return getImage.ToResponse(), nil
}

func (d DefaultDiagramService) DeleteDiagram(id string) *entities.ApiError {
	rawDiagram, err := d.getDiagramById(id, true)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to get diagram, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	err = d.repo.DeleteDiagram(rawDiagram.ID)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to delete diagram",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	s3, err := awsS3.NewS3Client()
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to start s3",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	parsedLargeFileUrl, err := url.Parse(rawDiagram.BgImage)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to parse image url",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	largeKey := parsedLargeFileUrl.Path[1:]
	_ = s3.DeleteResource("axisforestry", largeKey)

	parsedSmallFileUrl, err := url.Parse(rawDiagram.Thumbnail)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to parse image url",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	smallKey := parsedSmallFileUrl.Path[1:]
	_ = s3.DeleteResource("axisforestry", smallKey)

	return nil
}

func (d DefaultDiagramService) GetRawDiagramById(id string, isAdmin bool) (*models.Diagram, *entities.ApiError) {
	getDiagram, err := d.getDiagramById(id, isAdmin)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to get diagram, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}
	return getDiagram, nil
}

func (d DefaultDiagramService) CreateFigDiagram(l, s string, request *dto.DiagramImage) *entities.ApiError {
	series := make(map[string][]string)
	var mods []models.Model
	tempFilt := make(map[uint]*models.Model)
	for _, model := range request.DiagramModels {
		modelId, err := models.DecodeHashId(model.ModelId)
		if err != nil || modelId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error"+err.Error(),
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}

		_, ok := tempFilt[modelId]
		if !ok {
			getModel, err := d.modelRepo.GetModel(modelId)
			if err != nil {
				bugsnag.Notify(err)
				return utils.FormatApiError(
					"failed to get model "+err.Error(),
					configs.ServerError,
					entities.E{"server": "failed to get model "},
				)
			}
			mods = append(mods, *getModel)

			if model.SeriesId != "" {
				serialId, err := models.DecodeHashId(model.SeriesId)
				if err != nil || serialId == 0 {
					bugsnag.Notify(err)
					return utils.FormatApiError(
						"serial id not found"+err.Error(),
						configs.ServerError,
						entities.E{"server": "serial id not found"},
					)
				}

				sl := getSeries(serialId, getModel.Serials)
				if sl != "" {
					series[getModel.Name] = []string{sl}
				}
			}

			tempFilt[modelId] = getModel
		} else {
			getModel := tempFilt[modelId]
			if model.SeriesId != "" {
				serialId, err := models.DecodeHashId(model.SeriesId)
				if err != nil || serialId == 0 {
					bugsnag.Notify(err)
					return utils.FormatApiError(
						"serial id not found"+err.Error(),
						configs.ServerError,
						entities.E{"server": "serial id not found"},
					)
				}

				sl := getSeries(serialId, getModel.Serials)
				if sl != "" {
					getActSeries := series[getModel.Name]
					getActSeries = append(getActSeries, sl)
					series[getModel.Name] = getActSeries
				}
			}
		}
	}

	ctId, err := models.DecodeHashId(request.DiagramCat.Id)
	if err != nil || ctId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error"+err.Error(),
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}
	c := models.DiagramCat{}
	c.ID = ctId

	cats := make([]models.DiagramCat, 0, 1)
	cats = append(cats, c)

	var pts []models.Part
	for _, pt := range request.Parts {
		pId, err := models.DecodeHashId(pt.Part.Id)
		if err != nil || pId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error"+err.Error(),
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}
		p := models.Part{}
		p.ID = pId
		p.PartDiagramNumber = pt.Number
		p.CatNote = pt.Catalog
		p.RecommendedQty = pt.RecommendedQty
		pts = append(pts, p)
	}

	diagram := models.Diagram{
		Name:        request.DiagramName,
		Description: request.DiagramName,
		Models:      mods,
		Parts:       pts,
		Cats:        cats,
		BgImage:     l,
		Thumbnail:   s,
		Status:      request.Status,
		Draft:       request.CanvasData,
	}

	if request.Controller.Id != "" {
		controllerId, err := models.DecodeHashId(request.Controller.Id)
		if err != nil || controllerId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error"+err.Error(),
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}

		diagram.ControllerId = &controllerId
	}

	if len(series) != 0 {
		data, err := json.Marshal(series)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to marshall json",
				configs.ServerError,
				entities.E{"error": err.Error()},
			)
		}

		diagram.Series = string(data)
	}

	err = d.repo.CreateDiagram(diagram)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create create",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}
	return nil
}

func getSeries(serialId uint, ss []*models.Serial) string {
	for _, serial := range ss {
		if serial.ID == serialId {
			return serial.SerialStart + " ~ " + serial.SerialEnd
		}
	}

	return ""
}

func (d DefaultDiagramService) UpdateFigDiagram(diagram *models.Diagram, l, s string, request *dto.DiagramImage) *entities.ApiError {
	series := make(map[string][]string)
	tempFilt := make(map[uint]*models.Model)
	var mods []models.Model
	for _, model := range request.DiagramModels {
		modelId, err := models.DecodeHashId(model.ModelId)
		if err != nil || modelId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error"+err.Error(),
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}

		_, ok := tempFilt[modelId]
		if !ok {
			getModel, err := d.modelRepo.GetModel(modelId)
			if err != nil {
				bugsnag.Notify(err)
				return utils.FormatApiError(
					"failed to get model "+err.Error(),
					configs.ServerError,
					entities.E{"server": "failed to get model "},
				)
			}
			mods = append(mods, *getModel)

			if model.SeriesId != "" {
				serialId, err := models.DecodeHashId(model.SeriesId)
				if err != nil || serialId == 0 {
					bugsnag.Notify(err)
					return utils.FormatApiError(
						"serial id not found"+err.Error(),
						configs.ServerError,
						entities.E{"server": "serial id not found"},
					)
				}

				sl := getSeries(serialId, getModel.Serials)
				if sl != "" {
					series[getModel.Name] = []string{sl}
				}
			}

			tempFilt[modelId] = getModel
		} else {
			getModel := tempFilt[modelId]
			if model.SeriesId != "" {
				serialId, err := models.DecodeHashId(model.SeriesId)
				if err != nil || serialId == 0 {
					bugsnag.Notify(err)
					return utils.FormatApiError(
						"serial id not found"+err.Error(),
						configs.ServerError,
						entities.E{"server": "serial id not found"},
					)
				}

				sl := getSeries(serialId, getModel.Serials)
				if sl != "" {
					getActSeries := series[getModel.Name]
					getActSeries = append(getActSeries, sl)
					series[getModel.Name] = getActSeries
				}
			}
		}
	}

	ctId, err := models.DecodeHashId(request.DiagramCat.Id)
	if err != nil || ctId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}
	c := models.DiagramCat{}
	c.ID = ctId

	cats := make([]models.DiagramCat, 0, 1)
	cats = append(cats, c)

	var pts []models.Part
	for _, pt := range request.Parts {
		pId, err := models.DecodeHashId(pt.Part.Id)
		if err != nil || pId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error",
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}
		p := models.Part{}
		p.ID = pId
		p.PartDiagramNumber = pt.Number
		p.CatNote = pt.Catalog
		p.RecommendedQty = pt.RecommendedQty
		pts = append(pts, p)
	}

	diagram.Cats = cats
	diagram.Description = request.DiagramName
	diagram.Name = request.DiagramName
	diagram.Models = mods
	diagram.Parts = pts
	diagram.BgImage = l
	diagram.Thumbnail = s
	diagram.Status = request.Status
	diagram.Draft = request.CanvasData

	if request.Controller.Id != "" {
		controllerId, err := models.DecodeHashId(request.Controller.Id)
		if err != nil || controllerId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error"+err.Error(),
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}

		diagram.ControllerId = &controllerId
	}

	if len(series) != 0 {
		data, err := json.Marshal(series)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to marshall json",
				configs.ServerError,
				entities.E{"error": err.Error()},
			)
		}

		diagram.Series = string(data)
	}

	err = d.repo.UpdateDiagram(diagram)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to update diagram",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}
	return nil
}

func (d DefaultDiagramService) GetFigureImages(limit, offset int, search string) ([]dto.FigureImageResponse, *entities.ApiError) {
	res, err := d.repo.GetFigureImageS(limit, offset, utils.String(search))
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	fis := make([]dto.FigureImageResponse, 0, len(res))
	for _, re := range res {
		fi := re.ToResponse()
		fis = append(fis, fi)
	}

	return fis, nil
}

func (d DefaultDiagramService) UploadFigureImage(file io.Reader, title, fileName, ext string) (dto.FigureImageResponse, *entities.ApiError) {
	s3, err := awsS3.NewS3Client()
	if err != nil {
		bugsnag.Notify(err)
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	fName := fmt.Sprintf("/figure_images/figures/%s/file.%s.png", fileName, "raw")
	s3Path, err := s3.UploadImage(fName, file)
	if err != nil {
		bugsnag.Notify(err)
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}

	res, err := d.repo.CreateFigureImage(models.FigureImage{
		Title:        title,
		SelectedType: "image",
		File:         s3Path,
	})
	if err != nil {
		bugsnag.Notify(err)
		return dto.FigureImageResponse{}, utils.FormatApiError(
			"failed to get s3 instance: "+err.Error(),
			configs.ServerError,
			entities.E{"server": "server error"},
		)
	}
	return res.ToResponse(), nil
}

func (d DefaultDiagramService) Diagrams(limit, offset int, f bool) ([]dto.DiagramResponse, *entities.ApiError) {
	getDiagrams, err := d.repo.Diagrams(limit, offset, f)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to get diagrams",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}
	var responses []dto.DiagramResponse
	responses = []dto.DiagramResponse{}
	for _, diagram := range getDiagrams {
		if f {
			diagram.IsAdmin = true
		}
		d := diagram.ToResponse()
		responses = append(responses, d)
	}

	return responses, nil
}

func (d DefaultDiagramService) CreateDiagram(request *dto.DiagramRequest) *entities.ApiError {
	controllerId, err := models.DecodeHashId(request.ControllerId)
	if err != nil || controllerId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	var images []models.Image
	for _, image := range request.ImagesId {
		imageId, err := models.DecodeHashId(image)
		if err != nil || controllerId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error",
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}
		img := models.Image{}
		img.ID = imageId
		images = append(images, img)
	}

	var mods []models.Model
	for _, model := range request.ModelsId {
		modelId, err := models.DecodeHashId(model)
		if err != nil || controllerId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error",
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}
		m := models.Model{}
		m.ID = modelId
		mods = append(mods, m)
	}

	var cats []models.DiagramCat
	for _, ct := range request.DiagramCatsId {
		ctId, err := models.DecodeHashId(ct)
		if err != nil || controllerId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error",
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}
		c := models.DiagramCat{}
		c.ID = ctId
		cats = append(cats, c)
	}

	var pts []models.Part
	for _, pt := range request.Parts {
		pId, err := models.DecodeHashId(pt.PartId)
		if err != nil || pId == 0 {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"internal server error",
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}
		p := models.Part{}
		p.ID = pId
		p.PartDiagramNumber = pt.PartDiagramNumber
		pts = append(pts, p)
	}

	diagram := models.Diagram{
		Name:        request.Name,
		Description: request.Description,
		Images:      images,
		Models:      mods,
		Parts:       pts,
		Cats:        cats,
		BgImage:     "https://picsum.photos/500/300",
		Thumbnail:   "",
	}

	err = d.repo.CreateDiagram(diagram)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create create",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	return nil
}

func (d DefaultDiagramService) GetDiagramById(hid string, isAdmin bool) (*dto.DiagramResponse, *entities.ApiError) {
	getDiagram, err := d.getDiagramById(hid, isAdmin)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to get diagram, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}
	if isAdmin {
		getDiagram.IsAdmin = true
	}
	response := getDiagram.ToResponse()

	return &response, nil
}

func (d DefaultDiagramService) getDiagramById(hid string, isAdmin bool) (*models.Diagram, error) {
	id, err := models.DecodeHashId(hid)
	if err != nil || id == 0 {
		return nil, err
	}

	getDiagram, err := d.repo.GetDiagramById(id)
	if err != nil {
		return nil, err
	}

	return getDiagram, nil
}

func (d DefaultDiagramService) GetDiagrams(limit, offset int, modelIds, catIds []string, f bool) ([]dto.DiagramResponse, *entities.ApiError) {
	var modelId []uint
	for _, id := range modelIds {
		id, _ := models.DecodeHashId(id)
		modelId = append(modelId, id)
	}

	var catId []uint
	for _, id := range catIds {
		id, _ := models.DecodeHashId(id)
		catId = append(catId, id)
	}

	getDiagrams, err := d.repo.GetDiagrams(limit, offset, modelId, catId)
	log.Println(getDiagrams, "lop")
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to create diagram category",
			configs.ServerError,
			entities.E{"error": err.Error()},
		)
	}

	var responses []dto.DiagramResponse
	responses = []dto.DiagramResponse{}
	for _, diagram := range getDiagrams {
		if f {
			diagram.IsAdmin = true
		}
		d := diagram.ToResponse()
		responses = append(responses, d)
	}

	return responses, nil
}

func NewDefaultDiagramService(repo diagrams.DiagramRepo, modelRepo man_models.ModelRepo) DiagramService {
	return &DefaultDiagramService{repo: repo, modelRepo: modelRepo}
}
