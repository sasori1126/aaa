package dto

type PartRequest struct {
	PartNumber            string   `json:"part_number"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Detail                string   `json:"detail"`
	Price                 float64  `json:"price"`
	PricePerMeter         float64  `json:"price_per_meter"`
	DealerPrice           float64  `json:"dealer_price"`
	DealerPricePercentage int      `json:"dealer_price_percentage"`
	DealerPricePerMeter   float64  `json:"dealer_price_per_meter"`
	SalePrice             float64  `json:"sale_price"`
	SalePricePerMeter     float64  `json:"sale_price_per_meter"`
	QuantityOnHand        int      `json:"quantity_on_hand"`
	QuantityOnOrder       int      `json:"quantity_on_order"`
	QuantityOnSaleOrder   int      `json:"quantity_on_sale_order"`
	QuantityRecommended   int      `json:"quantity_recommended"`
	Weight                float64  `json:"weight"`
	Length                float64  `json:"length"`
	Width                 float64  `json:"width"`
	Height                float64  `json:"height"`
	Status                string   `json:"status"`
	VideoUrl              string   `json:"video_url"`
	Seo                   string   `json:"seo"`
	MetaKeywords          string   `json:"meta_keywords"`
	MetaDescription       string   `json:"meta_description"`
	OemNumber             string   `json:"oem_number"`
	OemCompatible         string   `json:"oem_compatible"`
	Material              string   `json:"material"`
	CountryOfOrigin       string   `json:"country_of_origin"`
	Code                  string   `json:"code"`
	GuideLink             string   `json:"guide_link"`
	Featured              bool     `json:"featured"`
	CylinderType          string   `json:"cylinder_type"`
	WhereUsed             string   `json:"where_used"`
	Drivers               string   `json:"drivers"`
	PartTypes             []string `json:"types"`
}

type UpdatePartRequest struct {
	PartNumber            string   `form:"part_number" json:"part_number" binding:"required"`
	Name                  string   `form:"name" json:"name" binding:"required"`
	Description           string   `form:"description" json:"description" binding:"required"`
	Detail                string   `form:"detail" json:"detail"`
	Price                 float64  `form:"price" json:"price" binding:"required"`
	PricePerMeter         float64  `form:"price_per_meter" json:"price_per_meter"`
	DealerPrice           float64  `form:"dealer_price" json:"dealer_price"`
	DealerPricePercentage int      `form:"dealer_price_percentage" json:"dealer_price_percentage"`
	DealerPricePerMeter   float64  `form:"dealer_price_per_meter" json:"dealer_price_per_meter"`
	SalePrice             float64  `form:"sale_price" json:"sale_price"`
	SalePricePerMeter     float64  `form:"sale_price_per_meter" json:"sale_price_per_meter"`
	QuantityOnHand        int      `form:"quantity_on_hand" json:"quantity_on_hand" binding:"required"`
	QuantityOnOrder       int      `form:"quantity_on_order" json:"quantity_on_order"`
	QuantityOnSaleOrder   int      `form:"quantity_on_sale_order" json:"quantity_on_sale_order"`
	QuantityRecommended   int      `form:"quantity_recommended" json:"quantity_recommended" binding:"required"`
	Weight                float64  `form:"weight" json:"weight"`
	Length                float64  `form:"length" json:"length"`
	Width                 float64  `form:"width" json:"width"`
	Height                float64  `form:"height" json:"height"`
	Status                string   `form:"status" json:"status"`
	VideoUrl              string   `form:"video_url" json:"video_url"`
	Seo                   string   `form:"seo" json:"seo"`
	MetaKeywords          string   `form:"meta_keywords" json:"meta_keywords"`
	MetaDescription       string   `form:"meta_description" json:"meta_description"`
	OemNumber             string   `form:"oem_number" json:"oem_number" binding:"required"`
	OemCompatible         string   `form:"oem_compatible" json:"oem_compatible"`
	Material              string   `form:"material" json:"material"`
	CountryOfOrigin       string   `form:"country_of_origin" json:"country_of_origin"`
	Code                  string   `form:"code" json:"code"`
	GuideLink             string   `form:"guide_link" json:"guide_link"`
	Featured              bool     `form:"featured" json:"featured"`
	CylinderType          string   `form:"cylinder_type" json:"cylinder_type"`
	WhereUsed             string   `form:"where_used" json:"where_used"`
	Drivers               string   `form:"drivers" json:"drivers"`
	PartTypes             []string `form:"types" json:"types"`
}

type MergePartRequest struct {
	PartNumber string `json:"part_number"`
	MainItem   string `json:"main_item"`
}

type UpdatePriceRequest struct {
	ID       uint    `json:"id"`
	NewPrice float64 `json:"new_price"`
}

type PartPriceDifferenceUpdate struct {
	Id                    uint    `json:"id"`
	PartNumber            string  `json:"part_number"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	Detail                string  `json:"detail"`
	PricePerMeter         float64 `json:"price_per_meter"`
	DealerPrice           float64 `json:"dealer_price"`
	DealerPricePercentage int     `json:"dealer_price_percentage"`
	DealerPricePerMeter   float64 `json:"dealer_price_per_meter"`
	SalePrice             float64 `json:"sale_price"`
	SalePricePerMeter     float64 `json:"sale_price_per_meter"`
	OemNumber             string  `json:"oem_number"`
	OemCompatible         string  `json:"oem_compatible"`
	Code                  string  `json:"code"`
	OldPrice              float64 `json:"price"`
	NewPrice              float64 `json:"new_price"`
}

type DuplicatePartResponse struct {
	PartNumber string      `json:"part_number"`
	Name       string      `json:"name"`
	Parts      []EmbedPart `json:"parts"`
}

type PartResponse struct {
	Id                    string          `json:"id"`
	PartNumber            string          `json:"part_number"`
	Name                  string          `json:"name"`
	Description           string          `json:"description"`
	Detail                string          `json:"detail"`
	Price                 float64         `json:"price"`
	PricePerMeter         float64         `json:"price_per_meter"`
	DealerPrice           float64         `json:"dealer_price"`
	DealerPricePercentage int             `json:"dealer_price_percentage"`
	DealerPricePerMeter   float64         `json:"dealer_price_per_meter"`
	SalePrice             float64         `json:"sale_price"`
	SalePricePerMeter     float64         `json:"sale_price_per_meter"`
	QuantityOnHand        int             `json:"quantity_on_hand"`
	QuantityOnOrder       int             `json:"quantity_on_order"`
	QuantityOnSaleOrder   int             `json:"quantity_on_sale_order"`
	QuantityRecommended   int             `json:"quantity_recommended"`
	Weight                float64         `json:"weight"`
	Length                float64         `json:"length"`
	Width                 float64         `json:"width"`
	Height                float64         `json:"height"`
	Status                string          `json:"status"`
	VideoUrl              string          `json:"video_url"`
	Seo                   string          `json:"seo"`
	MetaKeywords          string          `json:"meta_keywords"`
	MetaDescription       string          `json:"meta_description"`
	OemNumber             string          `json:"oem_number"`
	OemCompatible         string          `json:"oem_compatible"`
	Material              string          `json:"material"`
	CountryOfOrigin       string          `json:"country_of_origin"`
	ListId                string          `json:"list_id"`
	Code                  string          `json:"code"`
	GuideLink             string          `json:"guide_link"`
	Featured              bool            `json:"featured"`
	CylinderType          string          `json:"cylinder_type"`
	WhereUsed             string          `json:"where_used"`
	Drivers               string          `json:"drivers"`
	Images                []ImageResponse `json:"images"`
	Types                 []EmbedPartType `json:"types"`
}

type EmbedPart struct {
	Id                    string          `json:"id"`
	PartNumber            string          `json:"part_number"`
	PartDiagramNumber     int             `json:"part_diagram_number"`
	Name                  string          `json:"name"`
	Description           string          `json:"description"`
	Detail                string          `json:"detail"`
	Price                 float64         `json:"price"`
	PricePerMeter         float64         `json:"price_per_meter"`
	DealerPrice           float64         `json:"dealer_price"`
	DealerPricePercentage int             `json:"dealer_price_percentage"`
	DealerPricePerMeter   float64         `json:"dealer_price_per_meter"`
	SalePrice             float64         `json:"sale_price"`
	SalePricePerMeter     float64         `json:"sale_price_per_meter"`
	QuantityOnHand        int             `json:"quantity_on_hand"`
	QuantityOnOrder       int             `json:"quantity_on_order"`
	QuantityOnSaleOrder   int             `json:"quantity_on_sale_order"`
	QuantityRecommended   int             `json:"quantity_recommended"`
	Weight                float64         `json:"weight"`
	Length                float64         `json:"length"`
	Width                 float64         `json:"width"`
	Height                float64         `json:"height"`
	Status                string          `json:"status"`
	VideoUrl              string          `json:"video_url"`
	Seo                   string          `json:"seo"`
	MetaKeywords          string          `json:"meta_keywords"`
	MetaDescription       string          `json:"meta_description"`
	OemNumber             string          `json:"oem_number"`
	OemCompatible         string          `json:"oem_compatible"`
	Material              string          `json:"material"`
	CountryOfOrigin       string          `json:"country_of_origin"`
	Code                  string          `json:"code"`
	GuideLink             string          `json:"guide_link"`
	Featured              bool            `json:"featured"`
	CylinderType          string          `json:"cylinder_type"`
	WhereUsed             string          `json:"where_used"`
	Drivers               string          `json:"drivers"`
	Images                []ImageResponse `json:"images"`
	NoOfImages            int             `json:"no_of_images"`
	Types                 []EmbedPartType `json:"types"`
}

type EmbedPartType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
