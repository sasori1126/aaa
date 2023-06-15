package dto

type ManufacturerResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EmbedManufacturer struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DiagramImage struct {
	DiagramImg    string       `json:"diagram_img" binding:"required"`
	Parts         []ActivePart `json:"parts"`
	DiagramName   string       `json:"diagram_name" binding:"required"`
	Status        string       `json:"status" binding:"required"`
	CanvasData    string       `json:"canvas_data" binding:"required"`
	DiagramModels []struct {
		Id       string `json:"id"`
		ModelId  string `json:"modelId"`
		SeriesId string `json:"seriesId"`
		Name     string `json:"name"`
	} `json:"diagram_models"`
	Figures []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		File string `json:"file"`
	} `json:"figures"`
	DiagramCat struct {
		Id          string                `json:"id"`
		Name        string                `json:"name"`
		Description string                `json:"description"`
		SubCategory []SubCategoryResponse `json:"sub_category"`
	} `json:"diagram_cat" binding:"required"`
	Controller struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"controller" binding:"required"`
}

type ActivePart struct {
	Number         int          `json:"number"`
	Catalog        string       `json:"catalogNote"`
	RecommendedQty int          `json:"recommendedQty"`
	Part           PartResponse `json:"part"`
}

type ManufacturerRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
