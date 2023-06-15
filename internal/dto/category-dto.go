package dto

type CategoryResponse struct {
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	SubCategories []SubCategoryResponse `json:"sub_categories"`
	//Diagrams      []EmbedDiagram        `json:"diagrams"`
}

type SubCategoryResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type SubCategoryRequest struct {
	CategoryId  string `json:"category_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
