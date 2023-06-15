package dto

type DiagramRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	ControllerId  string   `json:"controller_id"`
	DiagramCatsId []string `json:"diagram_cats_id"`
	ModelsId      []string `json:"models_id"`
	ImagesId      []string `json:"images_id"`
	Parts         []struct {
		PartDiagramNumber int    `json:"part_diagram_number"`
		PartId            string `json:"part_id"`
		CatNote           string `json:"catalogNote"`
	} `json:"parts"`
}

type FigureImageResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	File  string `json:"file"`
}

type DiagramResponse struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	BgImage     string             `json:"image"`
	Models      []EmbedModel       `json:"models"`
	Controller  ControllerResponse `json:"controller"`
	Cats        []CategoryResponse `json:"cats"`
	Draft       string             `json:"draft"`
	Parts       []EmbedPart        `json:"parts"`
}

type SearchByModelCat struct {
	ModelId string `json:"model_id"`
	CatId   string `json:"cat_id"`
}

type EmbedDiagram struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BgImage     string `json:"image"`
}
