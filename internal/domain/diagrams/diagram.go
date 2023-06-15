package diagrams

import "axis/ecommerce-backend/internal/models"

type DiagramRepo interface {
	CreateDiagram(data models.Diagram) error
	UpdateDiagram(data *models.Diagram) error
	GetDiagrams(limit, offset int, modelId []uint, catId []uint) ([]models.Diagram, error)
	Diagrams(limit, offset int, f bool) ([]models.Diagram, error)
	GetDiagramById(id uint) (*models.Diagram, error)
	DeleteDiagram(id uint) error
	CreateFigureImage(data models.FigureImage) (*models.FigureImage, error)
	UpdateFigImage(data *models.FigureImage) error
	SaveFromCsv(data []models.FigureImage) error
	GetFigImage(id uint) (*models.FigureImage, error)
	DeleteFigImage(id uint) error
	GetFigureImageS(limit, offset int, search *string) ([]models.FigureImage, error)
}
