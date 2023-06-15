package dto

type SerialRequest struct {
	SerialStart string `json:"serial_start" binding:"required"`
	SerialEnd   string `json:"serial_end" binding:"required"`
	YearStart   string `json:"year_start" binding:"required"`
	YearEnd     string `json:"year_end" binding:"required"`
}

type SerialResponse struct {
	Id          string `json:"id"`
	SerialStart string `json:"serial_start"`
	SerialEnd   string `json:"serial_end"`
	YearStart   string `json:"year_start"`
	YearEnd     string `json:"year_end"`
}
