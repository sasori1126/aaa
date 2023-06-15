package diagrams

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type DiagramRepoDb struct {
	client storage.Storage
}

func (d DiagramRepoDb) SaveFromCsv(data []models.FigureImage) error {
	return d.client.SaveImagesFromCsv(data)
}

func (d DiagramRepoDb) UpdateFigImage(data *models.FigureImage) error {
	return d.client.UpdateFigImage(data)
}

func (d DiagramRepoDb) GetFigImage(id uint) (*models.FigureImage, error) {
	return d.client.GetFigImage(id)
}

func (d DiagramRepoDb) DeleteFigImage(id uint) error {
	return d.client.DeleteFigImage(id)
}

func (d DiagramRepoDb) DeleteDiagram(id uint) error {
	return d.client.DeleteDiagram(id)
}

func (d DiagramRepoDb) UpdateDiagram(data *models.Diagram) error {
	return d.client.UpdateDiagram(data)
}

func (d DiagramRepoDb) GetFigureImageS(limit, offset int, search *string) ([]models.FigureImage, error) {
	return d.client.GetFigureImages(limit, offset, search)
}

func (d DiagramRepoDb) Diagrams(limit, offset int, f bool) ([]models.Diagram, error) {
	return d.client.Diagrams(limit, offset, f)
}

func (d DiagramRepoDb) CreateDiagram(data models.Diagram) error {
	return d.client.CreateDiagram(data)
}

func (d DiagramRepoDb) CreateFigureImage(data models.FigureImage) (*models.FigureImage, error) {
	return d.client.CreateFigureImage(data)
}

func (d DiagramRepoDb) GetDiagramById(id uint) (*models.Diagram, error) {
	return d.client.GetDiagramByField(models.FindByField{Field: "id", Value: id})
}

func (d DiagramRepoDb) GetDiagrams(limit, offset int, modelId []uint, catId []uint) ([]models.Diagram, error) {
	return d.client.GetDiagrams(limit, offset, modelId, catId)
}

func NewDiagramRepoDb(db storage.Storage) DiagramRepo {
	return &DiagramRepoDb{client: db}
}
